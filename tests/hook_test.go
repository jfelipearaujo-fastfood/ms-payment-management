package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"
)

var opts = godog.Options{
	Format:      "pretty",
	Paths:       []string{"features"},
	Output:      colors.Colored(os.Stdout),
	Concurrency: 4,
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func TestFeatures(t *testing.T) {
	o := opts
	o.TestingT = t

	status := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options:             &o,
	}.Run()

	if status == 2 {
		t.SkipNow()
	}

	if status != 0 {
		t.Fatalf("zero status code expected, %d received", status)
	}
}

// Steps
const featureKey CtxKeyType = "feature"

type feature struct {
	HostApi    string
	ConnStr    string
	StateTitle string
}

var state = NewState[feature](featureKey)

func iHaveAPayment(ctx context.Context) (context.Context, error) {
	paymentId := "9dfa1386-2f52-4cca-b9aa-f9bd6887d442"

	feat := state.retrieve(ctx)

	conn, err := sql.Open("postgres", feat.ConnStr)
	if err != nil {
		return ctx, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	rows, err := conn.QueryContext(ctx, "SELECT COUNT(*) FROM payments WHERE payment_id = $1", paymentId)
	if err != nil {
		return ctx, err
	}

	var count int
	if ok := rows.Next(); !ok {
		return ctx, fmt.Errorf("payment with id %s not found", paymentId)
	}

	if err := rows.Scan(&count); err != nil {
		return ctx, err
	}

	if count == 0 {
		return ctx, fmt.Errorf("payment with id %s not found", paymentId)
	}

	return ctx, nil
}

func iSendAPaymentHookForAPayment(ctx context.Context, paymentType string) (context.Context, error) {
	paymentId := "9dfa1386-2f52-4cca-b9aa-f9bd6887d442"

	feat := state.retrieve(ctx)

	body := `{
		"approved": true
	}`

	if paymentType == "failed" {
		body = `{
			"approved": false
		}`
	}

	route := fmt.Sprintf("%s/payments/webhook/%s", feat.HostApi, paymentId)
	req, err := http.NewRequest("PATCH", route, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return ctx, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ctx, err
	}

	if res.StatusCode != http.StatusCreated {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return ctx, err
		}
		defer res.Body.Close()

		fmt.Printf("Body: %s", string(body))

		return ctx, fmt.Errorf("payment hook failed with status: %d", res.StatusCode)
	}

	defer res.Body.Close()

	var payment map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&payment); err != nil {
		return ctx, err
	}

	stateTitle, ok := payment["state_title"].(string)
	if !ok {
		return ctx, fmt.Errorf("State title not found")
	}

	feat.StateTitle = stateTitle

	return state.enrich(ctx, feat), nil
}

func thePaymentStateShouldBe(ctx context.Context, expectedState string) (context.Context, error) {
	feat := state.retrieve(ctx)

	if feat.StateTitle != expectedState {
		return ctx, fmt.Errorf("Expected state %s, got %s", expectedState, feat.StateTitle)
	}

	return ctx, nil
}

type testContext struct {
	network    *testcontainers.DockerNetwork
	containers []testcontainers.Container
}

var (
	containers = make(map[string]testContext)
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		network, err := network.New(ctx, network.WithCheckDuplicate(), network.WithDriver("bridge"))
		if err != nil {
			return ctx, err
		}

		postgresContainer, ctx, err := createPostgresContainer(ctx, network)
		if err != nil {
			return ctx, err
		}

		localstack, ctx, err := createLocalstackContainer(ctx, network)
		if err != nil {
			return ctx, err
		}

		apiContainer, ctx, err := createApiContainer(ctx, network)
		if err != nil {
			return ctx, err
		}

		containers[sc.Id] = testContext{
			network: network,
			containers: []testcontainers.Container{
				postgresContainer,
				localstack,
				apiContainer,
			},
		}

		return ctx, nil
	})

	ctx.Step(`^I have a payment$`, iHaveAPayment)
	ctx.Step(`^I send a payment hook for a "([^"]*)" payment$`, iSendAPaymentHookForAPayment)
	ctx.Step(`^the payment state should be "([^"]*)"$`, thePaymentStateShouldBe)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if err != nil {
			return ctx, err
		}

		tc := containers[sc.Id]

		for _, c := range tc.containers {
			err := c.Terminate(ctx)
			if err != nil {
				return ctx, err
			}
		}

		err = tc.network.Remove(ctx)

		return ctx, err
	})
}

func createPostgresContainer(ctx context.Context, network *testcontainers.DockerNetwork) (testcontainers.Container, context.Context, error) {
	dbScript, err := filepath.Abs(filepath.Join(".", "testdata", "init-db.sql"))
	if err != nil {
		return nil, ctx, err
	}

	dbScriptReader, err := os.Open(dbScript)
	if err != nil {
		return nil, ctx, err
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "postgres:16.0",
			ExposedPorts: []string{
				"5432",
			},
			Env: map[string]string{
				"POSTGRES_DB":       "payment_db",
				"POSTGRES_USER":     "payment",
				"POSTGRES_PASSWORD": "payment",
			},
			Networks: []string{
				network.Name,
			},
			NetworkAliases: map[string][]string{
				network.Name: {
					"test",
				},
			},
			Files: []testcontainers.ContainerFile{
				{
					Reader:            dbScriptReader,
					ContainerFilePath: "/docker-entrypoint-initdb.d/init.sql",
					FileMode:          0644,
				},
			},
			WaitingFor: wait.ForLog("PostgreSQL init process complete; ready for start up").WithStartupTimeout(120 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to start postgres container: %w", err)
	}

	postgresIp, err := container.Host(ctx)
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to get postgres ip: %w", err)
	}

	postgresPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to get postgres port: %w", err)
	}

	feat := state.retrieve(ctx)
	feat.ConnStr = fmt.Sprintf("postgres://payment:payment@%s:%s/payment_db?sslmode=disable", postgresIp, postgresPort.Port())

	return container, state.enrich(ctx, feat), nil
}

func createLocalstackContainer(ctx context.Context, network *testcontainers.DockerNetwork) (testcontainers.Container, context.Context, error) {
	snsScript, err := filepath.Abs(filepath.Join(".", "testdata", "init-sns.sh"))
	if err != nil {
		return nil, ctx, err
	}

	sqsScript, err := filepath.Abs(filepath.Join(".", "testdata", "init-sqs.sh"))
	if err != nil {
		return nil, ctx, err
	}

	smScript, err := filepath.Abs(filepath.Join(".", "testdata", "init-sm.sh"))
	if err != nil {
		return nil, ctx, err
	}

	snsScriptReader, err := os.Open(snsScript)
	if err != nil {
		return nil, ctx, err
	}

	sqsScriptReader, err := os.Open(sqsScript)
	if err != nil {
		return nil, ctx, err
	}

	smScriptReader, err := os.Open(smScript)
	if err != nil {
		return nil, ctx, err
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "localstack/localstack:latest",
			ExposedPorts: []string{
				"4566",
			},
			Env: map[string]string{
				"SERVICES":       "secretsmanager,sqs,sns",
				"DEFAULT_REGION": "us-east-1",
				"DOCKER_HOST":    "unix:///var/run/docker.sock",
			},
			Networks: []string{
				network.Name,
			},
			NetworkAliases: map[string][]string{
				network.Name: {
					"test",
				},
			},
			Files: []testcontainers.ContainerFile{
				{
					Reader:            snsScriptReader,
					ContainerFilePath: "/etc/localstack/init/ready.d/init-sns.sh",
					FileMode:          0777,
				},
				{
					Reader:            sqsScriptReader,
					ContainerFilePath: "/etc/localstack/init/ready.d/init-sqs.sh",
					FileMode:          0777,
				},
				{
					Reader:            smScriptReader,
					ContainerFilePath: "/etc/localstack/init/ready.d/init-sm.sh",
					FileMode:          0777,
				},
			},
			WaitingFor: wait.ForListeningPort("4566/tcp").WithStartupTimeout(120 * time.Second),
		},
		Started: true,
	})

	if err != nil {
		return nil, ctx, err
	}

	return container, ctx, nil
}

func createApiContainer(ctx context.Context, network *testcontainers.DockerNetwork) (testcontainers.Container, context.Context, error) {
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    "../",
				Dockerfile: "Dockerfile",
				KeepImage:  true,
			},
			ExposedPorts: []string{
				"8080",
			},
			Env: map[string]string{
				"API_PORT":                        "8080",
				"API_ENV_NAME":                    "development",
				"API_VERSION":                     "v1",
				"DB_URL":                          "todo",
				"DB_URL_SECRET_NAME":              "db-secret-url",
				"AWS_ACCESS_KEY_ID":               "test",
				"AWS_SECRET_ACCESS_KEY":           "test",
				"AWS_REGION":                      "us-east-1",
				"AWS_BASE_ENDPOINT":               "http://test:4566",
				"AWS_ORDER_PRODUCTION_TOPIC_NAME": "OrderProductionTopic",
				"AWS_UPDATE_ORDER_TOPIC_NAME":     "UpdateOrderTopic",
				"AWS_ORDER_PAYMENT_QUEUE_NAME":    "OrderPaymentQueue",
			},
			Networks: []string{
				network.Name,
			},
			NetworkAliases: map[string][]string{
				network.Name: {
					"test",
				},
			},
			WaitingFor: wait.ForLog("Server started").WithStartupTimeout(120 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		return nil, ctx, err
	}

	ports, err := container.Ports(ctx)
	if err != nil {
		return nil, ctx, err
	}

	if len(ports["8080/tcp"]) == 0 {
		return nil, ctx, fmt.Errorf("Port 8080/tcp not found")
	}

	port := ports["8080/tcp"][0].HostPort

	res, err := http.Get(fmt.Sprintf("http://localhost:%s/health", port))
	if err != nil {
		return nil, ctx, err
	}

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, ctx, err
		}
		defer res.Body.Close()

		fmt.Printf("Body: %s", string(body))

		return nil, ctx, fmt.Errorf("API health check failed with status: %d", res.StatusCode)
	}

	feat := state.retrieve(ctx)
	feat.HostApi = fmt.Sprintf("http://localhost:%s/api/v1", port)

	return container, state.enrich(ctx, feat), nil
}
