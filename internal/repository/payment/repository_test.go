package payment

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/entity/payment_entity"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	t.Run("Should create payment", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectBegin()

		mock.ExpectExec("INSERT INTO (.+)?payments(.+)?").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("INSERT INTO (.+)?payment_items(.+)?").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		repo := NewPaymentRepository(db)

		// Act
		err = repo.Create(ctx, &payment_entity.Payment{
			Items: []payment_entity.PaymentItem{
				payment_entity.NewPaymentItem(uuid.NewString(), "item", 1),
			},
		})

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return error if an error occurs", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectBegin()

		mock.ExpectExec("INSERT INTO (.+)?payments(.+)?").
			WillReturnError(assert.AnError)

		mock.ExpectRollback()

		repo := NewPaymentRepository(db)

		// Act
		err = repo.Create(ctx, &payment_entity.Payment{})

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should rollback if an error occurs when inserting payment items", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectBegin()

		mock.ExpectExec("INSERT INTO (.+)?payments(.+)?").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("INSERT INTO (.+)?payment_items(.+)?").
			WillReturnError(assert.AnError)

		mock.ExpectRollback()

		repo := NewPaymentRepository(db)

		// Act
		err = repo.Create(ctx, &payment_entity.Payment{
			Items: []payment_entity.PaymentItem{
				payment_entity.NewPaymentItem(uuid.NewString(), "item", 1),
			},
		})

		// Assert
		assert.NoError(t, err)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("Should get payment by id", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		expectedPaymentItems := []payment_entity.PaymentItem{
			payment_entity.NewPaymentItem(uuid.NewString(), "item", 1),
		}

		expectedPayment := payment_entity.Payment{
			OrderId:    "order_id",
			PaymentId:  "payment_id",
			Items:      expectedPaymentItems,
			TotalItems: 1,
			Amount:     1.0,
			State:      payment_entity.WaitingForApproval,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		mock.ExpectQuery("SELECT (.+)?payments(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"order_id", "payment_id", "total_items", "amount", "state", "created_at", "updated_at"}).
				AddRow(expectedPayment.OrderId, expectedPayment.PaymentId, expectedPayment.TotalItems, expectedPayment.Amount, expectedPayment.State, expectedPayment.CreatedAt, expectedPayment.UpdatedAt))

		mock.ExpectQuery("SELECT (.+)?payment_items(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "quantity"}).
				AddRow(expectedPaymentItems[0].Id, expectedPaymentItems[0].Name, expectedPaymentItems[0].Quantity))

		repo := NewPaymentRepository(db)

		// Act
		payment, err := repo.GetByID(ctx, "payment_id")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedPayment, payment)
	})

	t.Run("Should return error if an error occurs", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectQuery("SELECT (.+)?payments(.+)?").
			WillReturnError(assert.AnError)

		repo := NewPaymentRepository(db)

		// Act
		payment, err := repo.GetByID(ctx, "payment_id")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, payment)
	})

	t.Run("Should return error payment if no payment found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectQuery("SELECT (.+)?payments(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"order_id", "payment_id", "total_items", "amount", "state", "created_at", "updated_at"}))

		repo := NewPaymentRepository(db)

		// Act
		payment, err := repo.GetByID(ctx, "payment_id")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, payment)
	})

	t.Run("Should return error if scan fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectQuery("SELECT (.+)?payments(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"order_id", "payment_id", "total_items", "amount", "state", "created_at", "updated_at"}).
				AddRow("order_id", "payment_id", 1, "abc", payment_entity.WaitingForApproval, time.Now(), time.Now()))

		repo := NewPaymentRepository(db)

		// Act
		payment, err := repo.GetByID(ctx, "payment_id")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, payment)
	})
}

func TestGetByOrderID(t *testing.T) {
	t.Run("Should get payments by order id", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		expectedPayments := []payment_entity.Payment{
			{
				OrderId:    "order_id",
				PaymentId:  "payment_id",
				TotalItems: 1,
				Amount:     1.0,
				State:      payment_entity.WaitingForApproval,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			{
				OrderId:    "order_id",
				PaymentId:  "payment_id",
				TotalItems: 1,
				Amount:     1.0,
				State:      payment_entity.WaitingForApproval,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		}

		mock.ExpectQuery("SELECT (.+)?payments(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"order_id", "payment_id", "total_items", "amount", "state", "created_at", "updated_at"}).
				AddRow(expectedPayments[0].OrderId, expectedPayments[0].PaymentId, expectedPayments[0].TotalItems, expectedPayments[0].Amount, expectedPayments[0].State, expectedPayments[0].CreatedAt, expectedPayments[0].UpdatedAt).
				AddRow(expectedPayments[1].OrderId, expectedPayments[1].PaymentId, expectedPayments[1].TotalItems, expectedPayments[1].Amount, expectedPayments[1].State, expectedPayments[1].CreatedAt, expectedPayments[1].UpdatedAt))

		repo := NewPaymentRepository(db)

		// Act
		payments, err := repo.GetByOrderID(ctx, "order_id")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedPayments, payments)
	})

	t.Run("Should return error if an error occurs", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectQuery("SELECT (.+)?payments(.+)?").
			WillReturnError(assert.AnError)

		repo := NewPaymentRepository(db)

		// Act
		payments, err := repo.GetByOrderID(ctx, "order_id")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, payments)
	})

	t.Run("Should return empty payments if no payment found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectQuery("SELECT (.+)?payments(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"order_id", "payment_id", "total_items", "amount", "state", "created_at", "updated_at"}))

		repo := NewPaymentRepository(db)

		// Act
		payments, err := repo.GetByOrderID(ctx, "order_id")

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, payments)
	})

	t.Run("Should return error if scan fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectQuery("SELECT (.+)?payments(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"order_id", "payment_id", "total_items", "amount", "state", "created_at", "updated_at"}).
				AddRow("order_id", "payment_id", 1, "abc", payment_entity.WaitingForApproval, time.Now(), time.Now()))

		repo := NewPaymentRepository(db)

		// Act
		payments, err := repo.GetByOrderID(ctx, "order_id")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, payments)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Should update payment", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		expectedPayment := payment_entity.Payment{
			OrderId:    "order_id",
			PaymentId:  "payment_id",
			TotalItems: 1,
			Amount:     1.0,
			State:      payment_entity.WaitingForApproval,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		mock.ExpectExec("UPDATE (.+)?payments(.+)?").
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewPaymentRepository(db)

		// Act
		err = repo.Update(ctx, &expectedPayment)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return error if an error occurs", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectExec("UPDATE (.+)?payments(.+)?").
			WillReturnError(assert.AnError)

		repo := NewPaymentRepository(db)

		// Act
		err = repo.Update(ctx, &payment_entity.Payment{})

		// Assert
		assert.Error(t, err)
	})
}
