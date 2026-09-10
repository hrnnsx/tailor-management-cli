package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"tailor-management-cli/entity"
)

type OrderHandler struct {
	db *sql.DB
}

func NewOrderHandler(db *sql.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

// 1. Ambil ukuran tubuh terakhir milik user
func (h *OrderHandler) GetLatestMeasurement(userID int) (*entity.UserMeasurement, error) {
	var m entity.UserMeasurement
	query := `
		SELECT
			id, user_id, title, height_cm, chest_circumference, waist_circumference 
		FROM
			user_measurements
		WHERE
			user_id = ?
		ORDER BY
			id DESC LIMIT 1
	`
	err := h.db.QueryRow(query, userID).Scan(&m.ID, &m.UserID, &m.Title, &m.HeightCM, &m.ChestCircumference, &m.WaistCircumference)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data ukuran: %w", err)
	}
	return &m, nil
}

// 2. Simpan ukuran tubuh baru
func (h *OrderHandler) SaveMeasurement(userID int, height, chest, waist float64) (*entity.UserMeasurement, error) {
	title := "Default Profile"
	query := `
		INSERT INTO
			user_measurements (user_id, title, height_cm, chest_circumference, waist_circumference) 
		VALUES (?, ?, ?, ?, ?)
	`
	res, err := h.db.Exec(query, userID, title, height, chest, waist)
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan ukuran badan: %w", err)
	}

	lastID, _ := res.LastInsertId()
	return &entity.UserMeasurement{
		ID:                 int(lastID),
		UserID:             userID,
		Title:              title,
		HeightCM:           height,
		ChestCircumference: chest,
		WaistCircumference: waist,
	}, nil
}

// 3. Kebutuhan bahan kain dari size_requirements
func (h *OrderHandler) CalculateSizeAndRequirement(height, chest, waist float64) (string, int, error) {

	var size string
	switch {
	case chest <= 88 && waist <= 74:
		size = "XS"
	case chest <= 94 && waist <= 80:
		size = "S"
	case chest <= 102 && waist <= 86:
		size = "M"
	case chest <= 110 && waist <= 94:
		size = "L"
	case chest <= 118 && waist <= 102:
		size = "XL"
	default:
		size = "XXL"
	}

	// tinggi banget tp ramping
	if height > 180 && (size == "XS" || size == "S") {
		size = "M" // Supaya ga ngatung jadi tanktop
	}

	// Ambil kebutuhan kain dari size_requirements di database
	var reqCM int
	err := h.db.QueryRow("SELECT required_cm FROM size_requirements WHERE size = ?", size).Scan(&reqCM)
	if err != nil {
		reqCM = 240 // Default fallback 2.4 meter
	}

	return size, reqCM, nil
}

// 4. Ambil daftar kain dengan harga per cm
func (h *OrderHandler) GetFabrics() ([]entity.Fabric, error) {
	rows, err := h.db.Query("SELECT id, name, price_per_cm FROM fabrics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fabrics []entity.Fabric
	for rows.Next() {
		var f entity.Fabric
		if err := rows.Scan(&f.ID, &f.Name, &f.PricePerCM); err != nil {
			return nil, err
		}
		fabrics = append(fabrics, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return fabrics, nil
}

// 5. Ambil varian corak dengan harga per cm
func (h *OrderHandler) GetPatternsByFabric(fabricID int) ([]entity.FabricPatternOption, error) {
	query := `
		SELECT
			fp.id, fp.fabric_id, fp.pattern_id, p.name, fp.stock_cm, f.price_per_cm
		FROM
			fabric_patterns fp JOIN patterns p
				ON fp.pattern_id = p.id JOIN fabrics f
					ON fp.fabric_id = f.id
		WHERE
			fp.fabric_id = ?
		ORDER BY
			p.id ASC
	`
	rows, err := h.db.Query(query, fabricID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []entity.FabricPatternOption
	for rows.Next() {
		var opt entity.FabricPatternOption
		if err := rows.Scan(&opt.FabricPatternID, &opt.FabricID, &opt.PatternID, &opt.PatternName, &opt.StockCM, &opt.PricePerCM); err != nil {
			return nil, err
		}
		options = append(options, opt)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return options, nil
}

// 6. Hitung total harga langsung dalam satuan cm (tanpa konversi /100)
func (h *OrderHandler) CalculateSummary(fabric entity.Fabric, pattern entity.FabricPatternOption, size string, reqCM int) entity.OrderSummary {
	totalPrice := float64(reqCM) * fabric.PricePerCM
	return entity.OrderSummary{
		FabricName:  fabric.Name,
		PatternName: pattern.PatternName,
		Size:        size,
		RequiredCM:  reqCM,
		PricePerCM:  fabric.PricePerCM,
		TotalPrice:  totalPrice,
	}
}

// 7. Simpan ke orders dengan price_per_cm_snapshot
func (h *OrderHandler) SubmitOrder(customerID, measurementID, fabricPatternID int, size string, cmUsed int, pricePerCM, totalPrice float64) error {
	tx, err := h.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentStock int
	err = tx.QueryRow("SELECT stock_cm FROM fabric_patterns WHERE id = ? FOR UPDATE", fabricPatternID).Scan(&currentStock)
	if err != nil {
		return err
	}
	if currentStock < cmUsed {
		return errors.New("bahan tidak cukup / sold out")
	}

	_, err = tx.Exec("UPDATE fabric_patterns SET stock_cm = stock_cm - ? WHERE id = ?", cmUsed, fabricPatternID)
	if err != nil {
		return fmt.Errorf("gagal update stok kain: %w", err)
	}

	orderCode := fmt.Sprintf("ORD-%d", time.Now().UnixNano())
	orderQuery := `
		INSERT INTO orders (
			order_code, customer_id,
			user_measurement_id,
			fabric_pattern_id,
			determined_size,
			cm_used,
			price_per_cm_snapshot,
			total_price,
			payment_status,
			status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'unpaid', 'pending')
	`
	_, err = tx.Exec(
		orderQuery,
		orderCode,
		customerID,
		measurementID,
		fabricPatternID,
		size,
		cmUsed,
		pricePerCM,
		totalPrice,
	)
	if err != nil {
		return fmt.Errorf("gagal insert pesanan: %w", err)
	}

	return tx.Commit()
}

func (h *OrderHandler) CheckOrder(customerID int) ([]entity.CustomerOrder, error) {

	var orders []entity.CustomerOrder

	query := `
		SELECT
			id,
			order_code,
			determined_size,
			total_price,
			payment_status,
			status,
			created_at
		FROM orders
		WHERE customer_id = ?
		ORDER BY created_at DESC;
	`
	rows, err := h.db.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.CustomerOrder

		err := rows.Scan(
			&order.ID,
			&order.OrderCode,
			&order.DeterminedSize,
			&order.TotalPrice,
			&order.PaymentStatus,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
