package cli

import (
	"fmt"
	"strconv"
	"strings"

	"tailor-management-cli/entity"
	"tailor-management-cli/handler"

	"github.com/manifoldco/promptui"
)

func BikinBajuCLI(orderHandler *handler.OrderHandler, currentUser entity.User) {
	fmt.Printf("\n=== Menu Pemesanan Custom - Pelanggan: %s ===\n", currentUser.Name)

	measurement, err := orderHandler.GetLatestMeasurement(currentUser.ID)
	if err != nil {
		fmt.Printf("[Error]: %v\n", err)
		return
	}

	if measurement != nil {
		fmt.Printf("\nData ukuran tubuh Anda tersimpan:\n- TB: %.1f cm\n- Lingkar Dada: %.1f cm\n- Lingkar Pinggang: %.1f cm\n",
			measurement.HeightCM, measurement.ChestCircumference, measurement.WaistCircumference)

		measMenu := promptui.Select{
			Label: "Pilih Opsi Ukuran",
			Items: []string{
				"Lanjut menggunakan data yang telah disimpan",
				"Ukur ulang",
				"Kembali",
			},
		}
		idx, _, err := measMenu.Run()
		if err != nil || idx == 2 {
			return
		}

		if idx == 1 {
			measurement = inputNewMeasurementCLI(orderHandler, currentUser.ID)
			if measurement == nil {
				return
			}
		}
	} else {
		fmt.Println("\nAnda belum memiliki catatan ukuran tubuh. Silakan isi terlebih dahulu.")
		measurement = inputNewMeasurementCLI(orderHandler, currentUser.ID)
		if measurement == nil {
			return
		}
	}

	size, reqCM, err := orderHandler.CalculateSizeAndRequirement(
		measurement.HeightCM,
		measurement.ChestCircumference,
		measurement.WaistCircumference,
	)
	if err != nil {
		fmt.Printf("[Error]: %v\n", err)
		return
	}
	fmt.Printf("\n[Hasil Ukuran]: Ukuran baju: %s | Kebutuhan kain: %d cm\n", size, reqCM)

	for {
		fabrics, err := orderHandler.GetFabrics()
		if err != nil || len(fabrics) == 0 {
			fmt.Println("[Error]: Data katalog bahan tidak tersedia.")
			return
		}

		var fabricItems []string
		for _, f := range fabrics {
			fabricItems = append(fabricItems, fmt.Sprintf("%s - Rp %.2f/cm", f.Name, f.PricePerCM))
		}
		fabricItems = append(fabricItems, "Kembali ke menu utama")

		fabricPrompt := promptui.Select{
			Label: "Pilih Jenis Bahan Kain",
			Items: fabricItems,
		}
		fIdx, _, err := fabricPrompt.Run()
		if err != nil || fIdx == len(fabricItems)-1 {
			return
		}
		selectedFabric := fabrics[fIdx]

		selectedPattern := selectPatternWithSoldOutGuard(orderHandler, selectedFabric.ID, reqCM)
		if selectedPattern == nil {
			continue
		}

		summary := orderHandler.CalculateSummary(selectedFabric, *selectedPattern, size, reqCM)
		fmt.Println("\n================ RINGKASAN PESANAN ================")
		fmt.Printf("Bahan Kain     : %s\n", summary.FabricName)
		fmt.Printf("Corak Motif    : %s\n", summary.PatternName)
		fmt.Printf("Estimasi Size  : %s\n", summary.Size)
		fmt.Printf("Bahan Terpakai : %d cm\n", summary.RequiredCM)
		fmt.Printf("Harga/cm       : Rp %.2f\n", summary.PricePerCM)
		fmt.Printf("Total Biaya    : Rp %.2f\n", summary.TotalPrice)
		fmt.Println("===================================================")

		confirmPrompt := promptui.Select{
			Label: "Konfirmasi pembuatan pesanan?",
			Items: []string{
				"Ya, proses order",
				"Batal dan kembali",
			},
		}
		cIdx, _, err := confirmPrompt.Run()
		if err != nil || cIdx == 1 {
			fmt.Println("Pemesanan dibatalkan.")
			return
		}

		err = orderHandler.SubmitOrder(
			currentUser.ID,
			measurement.ID,
			selectedPattern.FabricPatternID,
			summary.Size,
			summary.RequiredCM,
			summary.PricePerCM,
			summary.TotalPrice,
		)
		if err != nil {
			fmt.Printf("[Error]: %v\n", err)
			return
		}

		fmt.Println("\nOrder berhasil dibuat! Status pesanan dapat dipantau di menu utama.")
		return
	}
}

func selectPatternWithSoldOutGuard(orderHandler *handler.OrderHandler, fabricID int, reqCM int) *entity.FabricPatternOption {
	patterns, err := orderHandler.GetPatternsByFabric(fabricID)
	if err != nil || len(patterns) == 0 {
		fmt.Println("[Pemberitahuan]: Corak untuk bahan ini belum tersedia.")
		return nil
	}

	var patternItems []string
	for _, p := range patterns {
		if p.StockCM < reqCM {
			patternItems = append(patternItems, fmt.Sprintf("%s [SOLD OUT - Sisa %d cm]", p.PatternName, p.StockCM))
		} else {
			patternItems = append(patternItems, fmt.Sprintf("%s (Tersedia %d cm)", p.PatternName, p.StockCM))
		}
	}
	patternItems = append(patternItems, "Kembali pilih bahan lain")

	for {
		patternPrompt := promptui.Select{
			Label: "Pilih Corak Motif",
			Items: patternItems,
		}
		pIdx, _, err := patternPrompt.Run()
		if err != nil || pIdx == len(patternItems)-1 {
			return nil
		}

		chosen := patterns[pIdx]
		if chosen.StockCM < reqCM {
			fmt.Printf("\n[Peringatan]: Corak '%s' sedang SOLD OUT / stok tidak cukup (butuh: %d cm, sisa: %d cm). Silakan pilih corak lain.\n\n",
				chosen.PatternName, reqCM, chosen.StockCM)
			continue
		}

		return &chosen
	}
}

func inputNewMeasurementCLI(orderHandler *handler.OrderHandler, userID int) *entity.UserMeasurement {
	inputTemplate := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . }} ",
		Invalid: "{{ . }} ",
		Success: " ",
	}

	promptFloat := func(label string) (float64, error) {
		fmt.Printf("\nMasukkan %s:\n", label)
		p := promptui.Prompt{
			Label:     ">",
			Templates: inputTemplate,
			Validate: func(input string) error {
				v, err := strconv.ParseFloat(strings.TrimSpace(input), 64)
				if err != nil || v <= 0 {
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

	tb, err := promptFloat("Tinggi Badan (cm)")
	if err != nil {
		return nil
	}
	ld, err := promptFloat("Lingkar Dada (cm)")
	if err != nil {
		return nil
	}
	lp, err := promptFloat("Lingkar Pinggang (cm)")
	if err != nil {
		return nil
	}

	m, err := orderHandler.SaveMeasurement(userID, tb, ld, lp)
	if err != nil {
		fmt.Printf("[Error]: %v\n", err)
		return nil
	}
	return m
}
