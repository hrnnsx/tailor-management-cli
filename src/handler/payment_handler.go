package handler

import (
	"database/sql"
	"fmt"

	"tailor-management-cli/entity"
)

type PaymentHandler struct {
	db *sql.DB
}

func NewPaymentHandler(db *sql.DB) *PaymentHandler {
	return &PaymentHandler{db: db}
}

func (h *PaymentHandler) CreatePayment(orderID int, amount float64) error {
	tx, err := h.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Pastikan order ada dan masih unpaid
	var paymentStatus string

	err = tx.QueryRow(`
		SELECT payment_status
		FROM orders
		WHERE id = ?
		FOR UPDATE
	`, orderID).Scan(&paymentStatus)

	if err == sql.ErrNoRows {
		return fmt.Errorf("order tidak ditemukan")
	}

	if err != nil {
		return err
	}

	if paymentStatus == "paid" {
		return fmt.Errorf("order ini sudah dibayar")
	}

	if paymentStatus != "unpaid" {
		return fmt.Errorf("status pembayaran order tidak valid: %s", paymentStatus)
	}

	// Pastikan belum ada payment yang masih pending
	var paymentID int

	err = tx.QueryRow(`
		SELECT id
		FROM payments
		WHERE order_id = ?
		  AND status = 'pending'
		LIMIT 1
	`, orderID).Scan(&paymentID)

	if err == nil {
		return fmt.Errorf(
			"order ini sudah memiliki pembayaran yang sedang menunggu verifikasi",
		)
	}

	if err != sql.ErrNoRows {
		return err
	}

	// Buat payment baru
	_, err = tx.Exec(`
		INSERT INTO payments (
			order_id,
			amount,
			status
		)
		VALUES (?, ?, 'pending')
	`, orderID, amount)

	if err != nil {
		return fmt.Errorf("gagal membuat pembayaran: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("gagal menyimpan pembayaran: %w", err)
	}

	return nil
}

func (h *PaymentHandler) GetPendingPayments() ([]entity.Payment, error) {
	var payments []entity.Payment

	query := `
		SELECT
			p.id,
			p.order_id,
			o.order_code,
			u.name,
			p.amount,
			p.status,
			p.created_at
		FROM payments p
		JOIN orders o ON p.order_id = o.id
		JOIN users u ON o.customer_id = u.id
		WHERE p.status = 'pending'
		ORDER BY p.created_at ASC
	`

	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var payment entity.Payment

		err := rows.Scan(
			&payment.ID,
			&payment.OrderID,
			&payment.OrderCode,
			&payment.CustomerName,
			&payment.Amount,
			&payment.Status,
			&payment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (h *PaymentHandler) VerifyPayment(paymentID int, adminID int) error {
	tx, err := h.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Ambil payment yang masih pending
	var orderID int

	err = tx.QueryRow(`
		SELECT order_id
		FROM payments
		WHERE id = ?
		  AND status = 'pending'
		FOR UPDATE
	`, paymentID).Scan(&orderID)

	if err == sql.ErrNoRows {
		return fmt.Errorf(
			"payment tidak ditemukan atau sudah diverifikasi",
		)
	}

	if err != nil {
		return fmt.Errorf(
			"gagal mengambil payment: %w",
			err,
		)
	}

	// Verifikasi payment
	result, err := tx.Exec(`
		UPDATE payments
		SET
			status = 'verified',
			verified_by = ?
		WHERE id = ?
		  AND status = 'pending'
	`, adminID, paymentID)

	if err != nil {
		return fmt.Errorf(
			"gagal memverifikasi payment: %w",
			err,
		)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"gagal mengecek hasil verifikasi: %w",
			err,
		)
	}

	if rowsAffected == 0 {
		return fmt.Errorf(
			"payment tidak ditemukan atau sudah diverifikasi",
		)
	}

	// Ubah payment status order menjadi paid
	_, err = tx.Exec(`
		UPDATE orders
		SET payment_status = 'paid'
		WHERE id = ?
	`, orderID)

	if err != nil {
		return fmt.Errorf(
			"gagal mengubah status pembayaran order: %w",
			err,
		)
	}

	// Kalau order sudah menunggu pembayaran,
	// otomatis selesai.
	_, err = tx.Exec(`
		UPDATE orders
		SET
			status = 'finished',
			progress = 'Selesai'
		WHERE id = ?
		  AND status = 'waiting payment'
	`, orderID)

	if err != nil {
		return fmt.Errorf(
			"gagal menyelesaikan order: %w",
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(
			"gagal menyimpan verifikasi pembayaran: %w",
			err,
		)
	}

	return nil
}
