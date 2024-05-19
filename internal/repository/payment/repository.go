package payment

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/doug-martin/goqu/v9"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/entity/payment_entity"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/shared/custom_error"
)

type PaymentRepository struct {
	conn *sql.DB
}

func NewPaymentRepository(conn *sql.DB) *PaymentRepository {
	return &PaymentRepository{
		conn: conn,
	}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *payment_entity.Payment) error {
	queryInsertPayment := `
		INSERT INTO payments (
			order_id,
			payment_id,
			total_items,
			amount,
			state,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7);
	`
	queryInsertPaymentItems := `
		INSERT INTO payment_items (
			id,
			order_id,
			payment_id,
			name,
			quantity
		)
		VALUES ($1,$2,$3,$4,$5);
	`

	tx, err := r.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		queryInsertPayment,
		payment.OrderId,
		payment.PaymentId,
		payment.TotalItems,
		payment.Amount,
		payment.State,
		payment.CreatedAt,
		payment.UpdatedAt)
	if err != nil {
		slog.ErrorContext(ctx, "error creating payment", "error", err)
		errTx := tx.Rollback()
		if errTx != nil {
			slog.ErrorContext(ctx, "error rolling back transaction", "error", errTx)
			return errTx
		}
		return err
	}

	for _, item := range payment.Items {
		_, err = tx.ExecContext(ctx,
			queryInsertPaymentItems,
			item.Id,
			payment.OrderId,
			payment.PaymentId,
			item.Name,
			item.Quantity)
		if err != nil {
			slog.ErrorContext(ctx, "error creating payment item", "item_id", item.Id, "error", err)
			errTx := tx.Rollback()
			if errTx != nil {
				slog.ErrorContext(ctx, "error rolling back transaction", "error", errTx)
				return errTx
			}
			return err
		}
	}

	return tx.Commit()
}

func (r *PaymentRepository) GetByID(ctx context.Context, paymentId string) (payment_entity.Payment, error) {
	sql, params, err := goqu.
		From("payments").
		Select("order_id", "payment_id", "total_items", "amount", "state", "created_at", "updated_at").
		Where(goqu.C("payment_id").Eq(paymentId)).
		ToSQL()
	if err != nil {
		return payment_entity.Payment{}, err
	}

	statement, err := r.conn.QueryContext(ctx, sql, params...)
	if err != nil {
		return payment_entity.Payment{}, err
	}
	defer statement.Close()

	var payment payment_entity.Payment

	for statement.Next() {
		err = statement.Scan(
			&payment.OrderId,
			&payment.PaymentId,
			&payment.TotalItems,
			&payment.Amount,
			&payment.State,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return payment_entity.Payment{}, err
		}
	}

	if payment.OrderId == "" {
		return payment_entity.Payment{}, custom_error.ErrPaymentNotFound
	}

	sql, params, err = goqu.
		From("payment_items").
		Select("id", "name", "quantity").
		Where(goqu.C("payment_id").Eq(paymentId)).
		ToSQL()
	if err != nil {
		return payment_entity.Payment{}, err
	}

	statement, err = r.conn.QueryContext(ctx, sql, params...)
	if err != nil {
		return payment_entity.Payment{}, err
	}
	defer statement.Close()

	payment.Items = make([]payment_entity.PaymentItem, 0)

	for statement.Next() {
		var item payment_entity.PaymentItem

		err = statement.Scan(
			&item.Id,
			&item.Name,
			&item.Quantity,
		)
		if err != nil {
			return payment_entity.Payment{}, err
		}

		payment.Items = append(payment.Items, item)
	}

	return payment, nil
}

func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderId string) ([]payment_entity.Payment, error) {
	var payments []payment_entity.Payment

	sql, params, err := goqu.
		From("payments").
		Where(goqu.C("order_id").Eq(orderId)).
		ToSQL()
	if err != nil {
		return payments, err
	}

	statement, err := r.conn.QueryContext(ctx, sql, params...)
	if err != nil {
		return payments, err
	}
	defer statement.Close()

	for statement.Next() {
		var payment payment_entity.Payment

		err = statement.Scan(
			&payment.OrderId,
			&payment.PaymentId,
			&payment.TotalItems,
			&payment.Amount,
			&payment.State,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return payments, err
		}

		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *PaymentRepository) Update(ctx context.Context, payment *payment_entity.Payment) error {
	sql, params, err := goqu.
		Update("payments").
		Set(goqu.Record{
			"state":      payment.State,
			"updated_at": payment.UpdatedAt,
		}).
		Where(goqu.C("payment_id").Eq(payment.PaymentId)).
		ToSQL()
	if err != nil {
		return err
	}

	_, err = r.conn.ExecContext(ctx, sql, params...)
	if err != nil {
		return err
	}

	return nil
}
