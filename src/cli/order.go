package cli

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"tailor-management-cli/entity"

	"github.com/manifoldco/promptui"
)

func CreateOrderFlow(db *sql.DB, user entity.User) error {
	fmt.Printf("\n=== Halo %s, Selamat Datang di Pemesanan Custom Tailor ===\n", user.Name)

	// 1. Cek User Measurements
	measurement, err := getOrInputMeasurement(db, user.ID)
	if err != nil {
		return err
	}
	if measurement == nil {
		return nil // User memilih kembali
	}

	// 2. Estimasi Size & Kebutuhan Bahan (cm)
	determinedSize := determineSize(measurement.ChestCircumference)
	var requiredCM int
	err = db.QueryRow("SELECT required_cm FROM size_requirements WHERE size = ?", determinedSize).Scan(&requiredCM)
	if err != nil {
		// Default fallback jika size_requirements belum terisi
		requiredCM = 250
	}

	fmt.Printf("\n[Hasil Ukuran] Estimasi Size Anda: %s (Kebutuhan Kain: %d cm / %.2f m)\n",
		determinedSize, requiredCM, float64(requiredCM)/100.0)

	// 3. Alur Pilih Bahan & Corak (Looping jika stok kurang atau ingin ganti bahan)
	for {
		fabric, err := selectFabric(db)
		if err != nil {
			return err
		}
		if fabric == nil {
			fmt.Println("Kembali ke menu utama.")
			return nil
		}

		selectedPattern, err := selectPattern(db, fabric.ID, requiredCM)
		if err != nil {
			return err
		}
		if selectedPattern == nil {
			// Kembali ke pemilihan bahan
			continue
		}

		// 4. Hitung Biaya & Konfirmasi Order
		meterUsed := float64(requiredCM) / 100.0
		totalPrice := meterUsed * fabric.PricePerCM

		fmt.Println("\n--- Ringkasan Pesanan ---")
		fmt.Printf("Bahan       : %s\n", fabric.Name)
		fmt.Printf("Corak       : %s\n", selectedPattern.PatternName)
		fmt.Printf("Estimasi Size: %s\n", determinedSize)
		fmt.Printf("Bahan Dipakai: %d cm (%.2f m)\n", requiredCM, meterUsed)
		fmt.Printf("Harga/meter : Rp %.2f\n", fabric.PricePerCM)
		fmt.Printf("Total Biaya : Rp %.2f\n", totalPrice)

		confirmPrompt := promptui.Select{
			Label: "Konfirmasi pembuatan pesanan?",
			Items: []string{"Ya, buat pesanan", "Batal / Kembali"},
		}
		cIdx, _, err := confirmPrompt.Run()
		if err != nil || cIdx == 1 {
			fmt.Println("Pemesanan dibatalkan.")
			return nil
		}

		// 5. Simpan Order & Potong Stok (Database Transaction)
		err = submitOrderTransaction(db, user.ID, measurement.ID, selectedPattern.FabricPatternID, determinedSize, requiredCM, fabric.PricePerCM, totalPrice)
		if err != nil {
			fmt.Printf("Gagal memproses order: %v\n", err)
			return err
		}

		fmt.Println("\nOrder berhasil dibuat! Status pesanan dapat dipantau di main menu.")
		return nil
	}
}

// Helpers untuk Step 1: Measurement
func getOrInputMeasurement(db *sql.DB, userID int) (*entity.UserMeasurement, error) {
	var m entity.UserMeasurement
	query := `SELECT id, user_id, title, height_cm, chest_circumference, waist_circumference 
	          FROM user_measurements WHERE user_id = ? ORDER BY id DESC LIMIT 1`
	err := db.QueryRow(query, userID).Scan(&m.ID, &m.UserID, &m.Title, &m.HeightCM, &m.ChestCircumference, &m.WaistCircumference)

	if err == sql.ErrNoRows {
		fmt.Println("\nBelum ada data ukuran tubuh. Mari isi ukuran Anda terlebih dahulu.")
		return promptNewMeasurement(db, userID)
	} else if err != nil {
		return nil, err
	}

	// Jika data sudah ada
	fmt.Printf("\nData Ukuran Tubuh Tersimpan:\n- TB: %.1f cm\n- Lingkar Dada: %.1f cm\n- Lingkar Pinggang: %.1f cm\n",
		m.HeightCM, m.ChestCircumference, m.WaistCircumference)

	prompt := promptui.Select{
		Label: "Pilihan Ukuran",
		Items: []string{"Lanjut menggunakan data yang telah disimpan", "Ukur ulang", "Kembali"},
	}
	idx, _, err := prompt.Run()
	if err != nil || idx == 2 {
		return nil, nil
	}

	if idx == 1 {
		return promptNewMeasurement(db, userID)
	}

	return &m, nil
}

func promptNewMeasurement(db *sql.DB, userID int) (*entity.UserMeasurement, error) {
	parseFloat := func(label string) (float64, error) {
		p := promptui.Prompt{
			Label: label,
			Validate: func(s string) error {
				val, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
				if err != nil || val <= 0 {
					return fmt.Errorf("masukkan angka positif yang valid")
				}
				return nil
			},
		}
		res, err := p.Run()
		if err != nil {
			return 0, err
		}
		return strconv.ParseFloat(strings.TrimSpace(res), 64)
	}

	tb, err := parseFloat("Tinggi Badan (cm)")
	if err != nil {
		return nil, err
	}
	ld, err := parseFloat("Lingkar Dada (cm)")
	if err != nil {
		return nil, err
	}
	lp, err := parseFloat("Lingkar Pinggang (cm)")
	if err != nil {
		return nil, err
	}

	title := "Default Profile"
	res, err := db.Exec(`INSERT INTO user_measurements (user_id, title, height_cm, chest_circumference, waist_circumference) 
	                     VALUES (?, ?, ?, ?, ?)`, userID, title, tb, ld, lp)
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan ukuran badan: %w", err)
	}

	lastID, _ := res.LastInsertId()
	return &entity.UserMeasurement{
		ID:                 int(lastID),
		UserID:             userID,
		Title:              title,
		HeightCM:           tb,
		ChestCircumference: ld,
		WaistCircumference: lp,
	}, nil
}

func determineSize(chest float64) string {
	switch {
	case chest < 88:
		return "XS"
	case chest <= 96:
		return "S"
	case chest <= 104:
		return "M"
	case chest <= 112:
		return "L"
	case chest <= 120:
		return "XL"
	default:
		return "XXL"
	}
}

// Helpers untuk Step 2 & 3: Bahan & Corak
func selectFabric(db *sql.DB) (*entity.Fabric, error) {
	rows, err := db.Query("SELECT id, name, price_per_meter FROM fabrics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fabrics []entity.Fabric
	var items []string
	for rows.Next() {
		var f entity.Fabric
		if err := rows.Scan(&f.ID, &f.Name, &f.PricePerCM); err != nil {
			return nil, err
		}
		fabrics = append(fabrics, f)
		items = append(items, fmt.Sprintf("%s (Rp %.2f/m)", f.Name, f.PricePerCM))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	items = append(items, "<- Kembali ke menu utama")

	prompt := promptui.Select{
		Label: "Pilih Jenis Bahan Kain",
		Items: items,
	}
	idx, _, err := prompt.Run()
	if err != nil || idx == len(items)-1 {
		return nil, nil
	}

	return &fabrics[idx], nil
}

func selectPattern(db *sql.DB, fabricID int, requiredCM int) (*entity.FabricPatternOption, error) {
	query := `
		SELECT fp.id, fp.fabric_id, fp.pattern_id, p.name, fp.stock_cm, f.price_per_meter
		FROM fabric_patterns fp
		JOIN patterns p ON fp.pattern_id = p.id
		JOIN fabrics f ON fp.fabric_id = f.id
		WHERE fp.fabric_id = ?`
	rows, err := db.Query(query, fabricID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []entity.FabricPatternOption
	var items []string
	for rows.Next() {
		var opt entity.FabricPatternOption
		if err := rows.Scan(&opt.FabricPatternID, &opt.FabricID, &opt.PatternID, &opt.PatternName, &opt.StockCM, &opt.PricePerCM); err != nil {
			return nil, err
		}
		options = append(options, opt)
		statusStock := fmt.Sprintf("%d cm", opt.StockCM)
		if opt.StockCM < requiredCM {
			statusStock += " [STOK TIDAK CUKUP]"
		}
		items = append(items, fmt.Sprintf("Corak: %s | Sisa: %s", opt.PatternName, statusStock))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	items = append(items, "<- Kembali pilih bahan lain")

	for {
		prompt := promptui.Select{
			Label: "Pilih Corak Motif",
			Items: items,
		}
		idx, _, err := prompt.Run()
		if err != nil || idx == len(items)-1 {
			return nil, nil // Switch ke pilih bahan lagi
		}

		selected := options[idx]
		if selected.StockCM < requiredCM {
			fmt.Printf("\n[Error]: Bahan tidak cukup / sold out (Dibutuhkan: %d cm, Tersedia: %d cm). Silakan pilih corak/bahan lain.\n\n",
				requiredCM, selected.StockCM)
			continue
		}

		return &selected, nil
	}
}

// Helpers untuk Step 4: Transaksi Pembuatan Order
func submitOrderTransaction(db *sql.DB, customerID, measurementID, fabricPatternID int, size string, cmUsed int, pricePerMeter, totalPrice float64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Kurangi stok bahan di fabric_patterns
	_, err = tx.Exec("UPDATE fabric_patterns SET stock_cm = stock_cm - ? WHERE id = ?", cmUsed, fabricPatternID)
	if err != nil {
		return fmt.Errorf("gagal mengupdate stok kain: %w", err)
	}

	// 2. Insert ke orders
	orderCode := fmt.Sprintf("ORD-%d", time.Now().UnixNano())
	orderQuery := `
		INSERT INTO orders (
			order_code, customer_id, user_measurement_id, fabric_pattern_id,
			determined_size, cm_used, price_per_meter_snapshot, total_price,
			payment_status, status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'unpaid', 'pending')`
	_, err = tx.Exec(orderQuery, orderCode, customerID, measurementID, fabricPatternID, size, cmUsed, pricePerMeter, totalPrice)
	if err != nil {
		return fmt.Errorf("gagal menyimpan data order: %w", err)
	}

	return tx.Commit()
}
