package handler

import (
	"database/sql"
	"fmt"
	"strings"

	"tailor-management-cli/entity"
)

type FabricHandler struct {
	db *sql.DB
}

func NewFabricHandler(db *sql.DB) *FabricHandler {
	return &FabricHandler{
		db: db,
	}
}

func (h *FabricHandler) GetFabricPatterns() ([]entity.FabricPattern, error) {
	query := `
		SELECT
			fp.id,
			fp.fabric_id,
			f.name,
			fp.pattern_id,
			p.name,
			fp.stock_cm,
			f.price_per_cm
		FROM fabric_patterns fp
		JOIN fabrics f ON f.id = fp.fabric_id
		JOIN patterns p ON p.id = fp.pattern_id
		ORDER BY fp.id ASC
	`

	rows, err := h.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal mengambil data stok kain: %w",
			err,
		)
	}
	defer rows.Close()

	var fabrics []entity.FabricPattern

	for rows.Next() {
		var fabric entity.FabricPattern

		err := rows.Scan(
			&fabric.ID,
			&fabric.FabricID,
			&fabric.FabricName,
			&fabric.PatternID,
			&fabric.PatternName,
			&fabric.StockCM,
			&fabric.PricePerCM,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"gagal membaca data stok kain: %w",
				err,
			)
		}

		fabrics = append(fabrics, fabric)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"gagal membaca hasil query: %w",
			err,
		)
	}

	return fabrics, nil
}

func (h *FabricHandler) RestockFabric(fabricPatternID int, quantityCM int) error {

	if quantityCM <= 0 {
		return fmt.Errorf(
			"jumlah restock harus lebih dari 0 cm",
		)
	}

	query := `
		UPDATE fabric_patterns
		SET stock_cm = stock_cm + ?
		WHERE id = ?
	`

	result, err := h.db.Exec(
		query,
		quantityCM,
		fabricPatternID,
	)

	if err != nil {
		return fmt.Errorf(
			"gagal melakukan restock kain: %w",
			err,
		)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"gagal memeriksa hasil restock: %w",
			err,
		)
	}

	if rowsAffected == 0 {
		return fmt.Errorf(
			"kombinasi kain dan pattern tidak ditemukan",
		)
	}

	return nil
}

func (h *FabricHandler) CreateFabric(name string, pricePerCM float64) error {
	// Bersihkan spasi di awal dan akhir nama
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("nama kain tidak boleh kosong")
	}

	if pricePerCM <= 0 {
		return fmt.Errorf("harga kain harus lebih dari 0")
	}

	// Cek apakah nama kain sudah ada
	var exists bool

	err := h.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM fabrics
			WHERE LOWER(name) = LOWER(?)
		)
	`, name).Scan(&exists)

	if err != nil {
		return fmt.Errorf(
			"gagal mengecek nama kain: %w",
			err,
		)
	}

	if exists {
		return fmt.Errorf(
			"kain %q sudah ada",
			name,
		)
	}

	// Insert kain baru
	query := `
		INSERT INTO fabrics (
			name,
			price_per_cm
		)
		VALUES (?, ?)
	`

	_, err = h.db.Exec(
		query,
		name,
		pricePerCM,
	)

	if err != nil {
		return fmt.Errorf(
			"gagal menambahkan jenis kain: %w",
			err,
		)
	}

	return nil
}

func (h *FabricHandler) CreatePattern(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("nama pattern tidak boleh kosong")
	}

	// Cek apakah pattern sudah ada
	var exists bool

	err := h.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM patterns
			WHERE LOWER(name) = LOWER(?)
		)
	`, name).Scan(&exists)

	if err != nil {
		return fmt.Errorf(
			"gagal mengecek nama pattern: %w",
			err,
		)
	}

	if exists {
		return fmt.Errorf(
			"pattern %q sudah ada",
			name,
		)
	}

	// Insert pattern baru
	query := `
		INSERT INTO patterns (
			name
		)
		VALUES (?)
	`

	_, err = h.db.Exec(
		query,
		name,
	)

	if err != nil {
		return fmt.Errorf(
			"gagal menambahkan pattern: %w",
			err,
		)
	}

	return nil
}

func (h *FabricHandler) CreateFabricPattern(fabricID int, patternID int, stockCM int) error {

	if fabricID <= 0 {
		return fmt.Errorf("ID kain tidak valid")
	}

	if patternID <= 0 {
		return fmt.Errorf("ID pattern tidak valid")
	}

	if stockCM < 0 {
		return fmt.Errorf("stok tidak boleh negatif")
	}

	var fabricExists bool

	err := h.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM fabrics
			WHERE id = ?
		)
	`, fabricID).Scan(&fabricExists)

	if err != nil {
		return fmt.Errorf(
			"gagal mengecek kain: %w",
			err,
		)
	}

	if !fabricExists {
		return fmt.Errorf(
			"kain dengan ID %d tidak ditemukan",
			fabricID,
		)
	}

	var patternExists bool

	err = h.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM patterns
			WHERE id = ?
		)
	`, patternID).Scan(&patternExists)

	if err != nil {
		return fmt.Errorf(
			"gagal mengecek pattern: %w",
			err,
		)
	}

	if !patternExists {
		return fmt.Errorf(
			"pattern dengan ID %d tidak ditemukan",
			patternID,
		)
	}

	var combinationExists bool

	err = h.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM fabric_patterns
			WHERE fabric_id = ?
			  AND pattern_id = ?
		)
	`, fabricID, patternID).Scan(&combinationExists)

	if err != nil {
		return fmt.Errorf(
			"gagal mengecek kombinasi kain dan pattern: %w",
			err,
		)
	}

	if combinationExists {
		return fmt.Errorf(
			"kombinasi kain dan pattern tersebut sudah ada",
		)
	}

	query := `
		INSERT INTO fabric_patterns (
			fabric_id,
			pattern_id,
			stock_cm
		)
		VALUES (?, ?, ?)
	`

	_, err = h.db.Exec(
		query,
		fabricID,
		patternID,
		stockCM,
	)

	if err != nil {
		return fmt.Errorf(
			"gagal menambahkan kombinasi kain dan pattern: %w",
			err,
		)
	}

	return nil
}

func (h *FabricHandler) GetFabrics() ([]entity.Fabric, error) {
	query := `
		SELECT
			id,
			name,
			price_per_cm
		FROM fabrics
		ORDER BY id ASC
	`

	rows, err := h.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal mengambil data kain: %w",
			err,
		)
	}
	defer rows.Close()

	var fabrics []entity.Fabric

	for rows.Next() {
		var fabric entity.Fabric

		err := rows.Scan(
			&fabric.ID,
			&fabric.Name,
			&fabric.PricePerCM,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"gagal membaca data kain: %w",
				err,
			)
		}

		fabrics = append(fabrics, fabric)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"gagal membaca hasil query: %w",
			err,
		)
	}

	return fabrics, nil
}

func (h *FabricHandler) GetPatterns() ([]entity.Pattern, error) {
	query := `
		SELECT
			id,
			name
		FROM patterns
		ORDER BY id ASC
	`

	rows, err := h.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal mengambil data pattern: %w",
			err,
		)
	}
	defer rows.Close()

	var patterns []entity.Pattern

	for rows.Next() {
		var pattern entity.Pattern

		err := rows.Scan(
			&pattern.ID,
			&pattern.Name,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"gagal membaca data pattern: %w",
				err,
			)
		}

		patterns = append(patterns, pattern)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"gagal membaca hasil query: %w",
			err,
		)
	}

	return patterns, nil
}
