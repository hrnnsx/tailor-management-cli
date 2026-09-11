package handler

import (
	"database/sql"
	"fmt"
	"tailor-management-cli/entity"
)

type WorkerHandler struct {
	db *sql.DB
}

func NewWorkerHandler(db *sql.DB) *WorkerHandler {
	return &WorkerHandler{
		db: db,
	}
}
func (h *WorkerHandler) GetMyOrders(workerID int) ([]entity.AdminOrder, error) {
	query := `
		SELECT
			o.id,
			o.order_code,
			u.name,
			o.determined_size,
			o.cm_used,
			o.assigned_worker_id,
			o.status,
			o.progress,
			o.created_at
		FROM orders o
		JOIN users u ON o.customer_id = u.id
		WHERE o.assigned_worker_id = ?
		  AND o.status = 'in progress'
		ORDER BY o.created_at ASC
	`

	rows, err := h.db.Query(query, workerID)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal mengambil order worker: %w",
			err,
		)
	}
	defer rows.Close()

	var orders []entity.AdminOrder

	for rows.Next() {
		var order entity.AdminOrder

		err := rows.Scan(
			&order.ID,
			&order.OrderCode,
			&order.CustomerName,
			&order.DeterminedSize,
			&order.CMUsed,
			&order.AssignedWorkerID,
			&order.Status,
			&order.Progress,
			&order.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"gagal membaca order worker: %w",
				err,
			)
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"gagal membaca hasil query: %w",
			err,
		)
	}

	return orders, nil
}

// ============================================================
// UPDATE PROGRESS
// ============================================================

func (h *WorkerHandler) UpdateProgress(
	workerID int,
	orderID int,
	progress string,
) error {

	// Validasi progress
	switch progress {
	case "Order diterima",
		"Persiapan bahan",
		"Pemotongan kain",
		"Proses jahit",
		"Finishing",
		"Selesai":

	default:
		return fmt.Errorf("progress tidak valid: %s", progress)
	}

	// Pastikan order memang milik worker tersebut
	var status string
	var paymentStatus string

	err := h.db.QueryRow(`
		SELECT
			status,
			payment_status
		FROM orders
		WHERE id = ?
		  AND assigned_worker_id = ?
	`, orderID, workerID).Scan(
		&status,
		&paymentStatus,
	)

	if err == sql.ErrNoRows {
		return fmt.Errorf(
			"order tidak ditemukan atau bukan milik worker ini",
		)
	}

	if err != nil {
		return fmt.Errorf(
			"gagal mengecek order: %w",
			err,
		)
	}

	if status != "in progress" {
		return fmt.Errorf(
			"order tidak bisa di-update karena statusnya %s",
			status,
		)
	}

	// ========================================================
	// FINISHING
	// ========================================================

	if progress == "Finishing" {

		// Kalau sudah bayar → langsung selesai
		if paymentStatus == "paid" {
			return h.finishOrder(orderID, workerID)
		}

		// Kalau belum bayar → waiting payment
		_, err := h.db.Exec(`
			UPDATE orders
			SET
				progress = 'Finishing',
				status = 'waiting payment'
			WHERE id = ?
			  AND assigned_worker_id = ?
			  AND status = 'in progress'
		`, orderID, workerID)

		if err != nil {
			return fmt.Errorf(
				"gagal mengubah order ke waiting payment: %w",
				err,
			)
		}

		// Worker kembali tersedia
		_, err = h.db.Exec(`
			UPDATE workers
			SET availability = TRUE
			WHERE user_id = ?
		`, workerID)

		if err != nil {
			return fmt.Errorf(
				"gagal mengembalikan availability worker: %w",
				err,
			)
		}

		return nil
	}

	// ========================================================
	// SELESAI
	// ========================================================

	if progress == "Selesai" {
		return h.finishOrder(orderID, workerID)
	}

	// ========================================================
	// PROGRESS BIASA
	// ========================================================

	_, err = h.db.Exec(`
		UPDATE orders
		SET progress = ?
		WHERE id = ?
		  AND assigned_worker_id = ?
		  AND status = 'in progress'
	`, progress, orderID, workerID)

	if err != nil {
		return fmt.Errorf(
			"gagal update progress: %w",
			err,
		)
	}

	return nil
}

// ============================================================
// FINISH ORDER
// ============================================================

func (h *WorkerHandler) finishOrder(
	orderID int,
	workerID int,
) error {

	// Pastikan pembayaran sudah lunas
	var paymentStatus string

	err := h.db.QueryRow(`
		SELECT payment_status
		FROM orders
		WHERE id = ?
		  AND assigned_worker_id = ?
	`, orderID, workerID).Scan(&paymentStatus)

	if err == sql.ErrNoRows {
		return fmt.Errorf(
			"order tidak ditemukan atau bukan milik worker ini",
		)
	}

	if err != nil {
		return fmt.Errorf(
			"gagal mengecek pembayaran: %w",
			err,
		)
	}

	if paymentStatus != "paid" {
		return fmt.Errorf(
			"order belum dibayar, tidak bisa diselesaikan",
		)
	}

	// Selesaikan order
	_, err = h.db.Exec(`
		UPDATE orders
		SET
			progress = 'Selesai',
			status = 'finished'
		WHERE id = ?
		  AND assigned_worker_id = ?
	`, orderID, workerID)

	if err != nil {
		return fmt.Errorf(
			"gagal menyelesaikan order: %w",
			err,
		)
	}

	// Worker kembali tersedia
	_, err = h.db.Exec(`
		UPDATE workers
		SET availability = TRUE
		WHERE user_id = ?
	`, workerID)

	if err != nil {
		return fmt.Errorf(
			"gagal mengembalikan availability worker: %w",
			err,
		)
	}

	return nil
}
