package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"tailor-management-cli/config/colors"
	"tailor-management-cli/entity"
	"tailor-management-cli/handler"

	"github.com/manifoldco/promptui"
)

var selectTemplates = &promptui.SelectTemplates{
	Active:   "▸ {{ . | cyan }}",
	Inactive: "  {{ . }}",
	Selected: "✔ {{ . | green }}",
}

func BikinBajuCLI(
	orderHandler *handler.OrderHandler,
	currentUser entity.User,
) {
	fmt.Println()
	fmt.Println("==============================================================")
	fmt.Printf("                    PEMESANAN CUSTOM\n")
	fmt.Printf("                    Pelanggan: %s\n", currentUser.Name)
	fmt.Println("==============================================================")
	fmt.Println()

	measurement, err := orderHandler.GetLatestMeasurement(currentUser.ID)
	if err != nil {
		fmt.Printf(
			"%s[FAILED]:%s %v\n",
			colors.Red,
			colors.Reset,
			err,
		)
		pause("Tekan Enter untuk kembali...")
		return
	}

	// Data ukuran
	if measurement != nil {
		fmt.Println("Data ukuran tubuh Anda tersimpan:")
		fmt.Printf("  Tinggi Badan      : %.1f cm\n", measurement.HeightCM)
		fmt.Printf("  Lingkar Dada      : %.1f cm\n", measurement.ChestCircumference)
		fmt.Printf("  Lingkar Pinggang  : %.1f cm\n", measurement.WaistCircumference)
		fmt.Println()

		measMenu := promptui.Select{
			Label: "Pilih Opsi Ukuran",
			Items: []string{
				"Lanjut menggunakan data yang telah disimpan",
				"Ukur ulang",
				"Kembali",
			},
			Templates: selectTemplates,
		}

		idx, _, err := measMenu.Run()
		if err != nil || idx == 2 {
			return
		}

		if idx == 1 {
			fmt.Println()
			fmt.Println("Silakan masukkan ukuran tubuh baru.")
			fmt.Println()

			measurement = inputNewMeasurementCLI(
				orderHandler,
				currentUser.ID,
			)

			if measurement == nil {
				return
			}
		}
	} else {
		fmt.Println(
			"Anda belum memiliki catatan ukuran tubuh.",
		)
		fmt.Println(
			"Silakan isi ukuran tubuh terlebih dahulu.",
		)
		fmt.Println()

		measurement = inputNewMeasurementCLI(
			orderHandler,
			currentUser.ID,
		)

		if measurement == nil {
			return
		}
	}

	// Hasil perhitungan ukuran
	size, reqCM, err := orderHandler.CalculateSizeAndRequirement(
		measurement.HeightCM,
		measurement.ChestCircumference,
		measurement.WaistCircumference,
	)

	if err != nil {
		fmt.Printf(
			"\n%s[FAILED]:%s %v\n",
			colors.Red,
			colors.Reset,
			err,
		)
		pause("Tekan Enter untuk kembali...")
		return
	}

	fmt.Println()
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("                    HASIL PERHITUNGAN")
	fmt.Println("--------------------------------------------------------------")
	fmt.Printf("Ukuran baju       : %s\n", size)
	fmt.Printf("Kebutuhan kain    : %d cm\n", reqCM)
	fmt.Println("--------------------------------------------------------------")
	fmt.Println()

	// Pilih bahan
	for {
		fabrics, err := orderHandler.GetFabrics()
		if err != nil || len(fabrics) == 0 {
			fmt.Printf(
				"%s[FAILED]:%s Data katalog bahan tidak tersedia.\n",
				colors.Red,
				colors.Reset,
			)
			pause("Tekan Enter untuk kembali...")
			return
		}

		var fabricItems []string

		for _, f := range fabrics {
			fabricItems = append(
				fabricItems,
				fmt.Sprintf(
					"%s - Rp %.2f/cm",
					f.Name,
					f.PricePerCM,
				),
			)
		}

		fabricItems = append(
			fabricItems,
			"Kembali ke menu customer",
		)

		fmt.Println("Pilih bahan kain yang ingin digunakan.")
		fmt.Println()

		fabricPrompt := promptui.Select{
			Label:     "Pilih Jenis Bahan Kain",
			Items:     fabricItems,
			Templates: selectTemplates,
		}

		fIdx, _, err := fabricPrompt.Run()
		if err != nil || fIdx == len(fabricItems)-1 {
			return
		}

		selectedFabric := fabrics[fIdx]

		fmt.Println()

		// Pilih corak
		selectedPattern := selectPatternWithSoldOutGuard(
			orderHandler,
			selectedFabric.ID,
			reqCM,
		)

		if selectedPattern == nil {
			continue
		}

		// Ringkasan pesanan
		summary := orderHandler.CalculateSummary(
			selectedFabric,
			*selectedPattern,
			size,
			reqCM,
		)

		fmt.Println()
		fmt.Println("==============================================================")
		fmt.Println("                    RINGKASAN PESANAN")
		fmt.Println("==============================================================")
		fmt.Printf("  Bahan Kain        : %s\n", summary.FabricName)
		fmt.Printf("  Corak Motif       : %s\n", summary.PatternName)
		fmt.Printf("  Estimasi Size     : %s\n", summary.Size)
		fmt.Printf("  Bahan Terpakai    : %d cm\n", summary.RequiredCM)
		fmt.Printf("  Harga/cm          : Rp %.2f\n", summary.PricePerCM)
		fmt.Printf("  Total Biaya       : Rp %.2f\n", summary.TotalPrice)
		fmt.Println("==============================================================")
		fmt.Println()

		confirmPrompt := promptui.Select{
			Label: "Konfirmasi pembuatan pesanan?",
			Items: []string{
				"Ya, proses order",
				"Batal dan kembali",
			},
			Templates: selectTemplates,
		}

		cIdx, _, err := confirmPrompt.Run()
		if err != nil || cIdx == 1 {
			fmt.Println()
			fmt.Println("Pemesanan dibatalkan.")
			return
		}

		// SUvmit order
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
			fmt.Printf(
				"\n%s[FAILED]:%s %v\n",
				colors.Red,
				colors.Reset,
				err,
			)
			pause("Tekan Enter untuk kembali...")
			return
		}

		fmt.Printf(
			"\n%s[SUCCESS]:%s Order berhasil dibuat!\n",
			colors.Green,
			colors.Reset,
		)
		fmt.Println(
			"Status pesanan dapat dipantau melalui menu utama.",
		)

		pause("\nTekan Enter untuk kembali...")
		return
	}
}

func selectPatternWithSoldOutGuard(
	orderHandler *handler.OrderHandler,
	fabricID int,
	reqCM int,
) *entity.FabricPatternOption {
	patterns, err := orderHandler.GetPatternsByFabric(fabricID)

	if err != nil || len(patterns) == 0 {
		fmt.Printf(
			"%s[WARNING]:%s Corak untuk bahan ini belum tersedia.\n",
			colors.Red,
			colors.Reset,
		)
		fmt.Println()
		return nil
	}

	var patternItems []string

	for _, p := range patterns {
		if p.StockCM < reqCM {
			patternItems = append(
				patternItems,
				fmt.Sprintf(
					"%s [SOLD OUT - Sisa %d cm]",
					p.PatternName,
					p.StockCM,
				),
			)
		} else {
			patternItems = append(
				patternItems,
				fmt.Sprintf(
					"%s (Tersedia %d cm)",
					p.PatternName,
					p.StockCM,
				),
			)
		}
	}

	patternItems = append(
		patternItems,
		"Kembali pilih bahan lain",
	)

	for {
		patternPrompt := promptui.Select{
			Label:     "Pilih Corak Motif",
			Items:     patternItems,
			Templates: selectTemplates,
		}

		pIdx, _, err := patternPrompt.Run()
		if err != nil || pIdx == len(patternItems)-1 {
			return nil
		}

		chosen := patterns[pIdx]

		if chosen.StockCM < reqCM {
			fmt.Printf(
				"\n%s[WARNING]:%s Corak '%s' sedang SOLD OUT / stok tidak cukup.\n",
				colors.Red,
				colors.Reset,
				chosen.PatternName,
			)
			fmt.Printf(
				"  Kebutuhan : %d cm\n",
				reqCM,
			)
			fmt.Printf(
				"  Stok      : %d cm\n",
				chosen.StockCM,
			)
			fmt.Println()
			continue
		}

		return &chosen
	}
}

func inputNewMeasurementCLI(
	orderHandler *handler.OrderHandler,
	userID int,
) *entity.UserMeasurement {
	inputTemplate := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . }} ",
		Invalid: "{{ . }} ",
		Success: " ",
	}

	promptFloat := func(label string) (float64, error) {
		fmt.Printf("Masukkan %s:\n", label)

		p := promptui.Prompt{
			Label:     ">",
			Templates: inputTemplate,
			Validate: func(input string) error {
				v, err := strconv.ParseFloat(
					strings.TrimSpace(input),
					64,
				)

				if err != nil || v <= 0 {
					return fmt.Errorf(
						"masukkan angka positif yang valid",
					)
				}

				return nil
			},
		}

		res, err := p.Run()
		if err != nil {
			return 0, err
		}

		return strconv.ParseFloat(
			strings.TrimSpace(res),
			64,
		)
	}

	tb, err := promptFloat("Tinggi Badan (cm)")
	if err != nil {
		return nil
	}

	fmt.Println()

	ld, err := promptFloat("Lingkar Dada (cm)")
	if err != nil {
		return nil
	}

	fmt.Println()

	lp, err := promptFloat("Lingkar Pinggang (cm)")
	if err != nil {
		return nil
	}

	m, err := orderHandler.SaveMeasurement(
		userID,
		tb,
		ld,
		lp,
	)

	if err != nil {
		fmt.Printf(
			"\n%s[FAILED]:%s %v\n",
			colors.Red,
			colors.Reset,
			err,
		)
		return nil
	}

	fmt.Printf(
		"\n%s[SUCCESS]:%s Data ukuran berhasil disimpan.\n",
		colors.Green,
		colors.Reset,
	)
	fmt.Println()

	return m
}

func CheckAllOrderStatus(
	orderH *handler.OrderHandler,
	customerID int,
) {
	fmt.Println()
	fmt.Println("==============================================================")
	fmt.Println("                    STATUS PESANAN")
	fmt.Println("==============================================================")
	fmt.Println()

	orders, err := orderH.CheckOrder(customerID)

	if err != nil {
		fmt.Printf(
			"%s[FAILED]:%s %v\n",
			colors.Red,
			colors.Reset,
			err,
		)
		pause("\nTekan Enter untuk kembali...")
		return
	}

	if len(orders) == 0 {
		fmt.Println("Belum ada order.")
		pause("\nTekan Enter untuk kembali...")
		return
	}

	fmt.Println(
		"---------------------------------------------------------------------",
	)
	fmt.Printf(
		"%-18s %-8s %-20s %-22s\n",
		"ORDER CODE",
		"SIZE",
		"STATUS",
		"PROGRESS",
	)
	fmt.Println(
		"---------------------------------------------------------------------",
	)

	for _, order := range orders {
		fmt.Printf(
			"%-18s %-8s %-20s %-22s\n",
			order.OrderCode,
			order.DeterminedSize,
			order.Status,
			order.Progress,
		)
	}

	fmt.Println(
		"---------------------------------------------------------------------",
	)

	pause("\nTekan Enter untuk kembali...")
}

func PayOrderCLI(
	paymentHandler *handler.PaymentHandler,
	orders []entity.CustomerOrder,
) {
	items := make([]string, len(orders))

	for i, order := range orders {
		items[i] = fmt.Sprintf(
			"%s - Rp%.2f",
			order.OrderCode,
			order.TotalPrice,
		)
	}

	fmt.Println()
	fmt.Println("Pilih order yang ingin dibayar.")
	fmt.Println()

	prompt := promptui.Select{
		Label:     "Pilih Order",
		Items:     items,
		Templates: selectTemplates,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return
	}

	selectedOrder := orders[index]

	fmt.Println()
	fmt.Println("--------------------------------------------------------------")
	fmt.Printf("  Order Code : %s\n", selectedOrder.OrderCode)
	fmt.Printf("  Total      : Rp%.2f\n", selectedOrder.TotalPrice)
	fmt.Printf("  Status     : %s\n", selectedOrder.PaymentStatus)
	fmt.Println("--------------------------------------------------------------")
	fmt.Println()

	confirmPrompt := promptui.Select{
		Label: "Bayar order ini?",
		Items: []string{
			"Ya, bayar",
			"Kembali",
		},
		Templates: selectTemplates,
	}

	confirmIndex, _, err := confirmPrompt.Run()
	if err != nil {
		return
	}

	switch confirmIndex {
	case 0:
		err := paymentHandler.CreatePayment(
			selectedOrder.ID,
			selectedOrder.TotalPrice,
		)

		if err != nil {
			fmt.Printf(
				"\n%s[FAILED]:%s %v\n",
				colors.Red,
				colors.Reset,
				err,
			)
			pause("\nTekan Enter untuk kembali...")
			return
		}

		fmt.Printf(
			"\n%s[SUCCESS]:%s Pembayaran berhasil dibuat.\n",
			colors.Green,
			colors.Reset,
		)
		fmt.Println("Status pembayaran: pending.")
		fmt.Println("Silakan tunggu verifikasi admin.")
		fmt.Println()

		pause("Tekan Enter untuk kembali...")

	case 1:
		return
	}
}

func CheckBill(
	orderHandler *handler.OrderHandler,
	paymentHandler *handler.PaymentHandler,
	userID int,
) {
	fmt.Println()
	fmt.Println("==============================================================")
	fmt.Println("                     TAGIHAN PESANAN")
	fmt.Println("==============================================================")
	fmt.Println()

	orders, err := orderHandler.CheckOrder(userID)

	if err != nil {
		fmt.Printf(
			"%s[FAILED]:%s %v\n",
			colors.Red,
			colors.Reset,
			err,
		)
		pause("\nTekan Enter untuk kembali...")
		return
	}

	if len(orders) == 0 {
		fmt.Println("Belum ada order.")
		pause("\nTekan Enter untuk kembali...")
		return
	}

	fmt.Println(
		"-----------------------------------------------------------",
	)
	fmt.Printf(
		"%-28s %-15s %-17s\n",
		"ORDER CODE",
		"PRICE",
		"PAYMENT STATUS",
	)
	fmt.Println(
		"-----------------------------------------------------------",
	)

	var unpaidOrders []entity.CustomerOrder
	var totalPrice float64

	for _, order := range orders {
		fmt.Printf(
			"%-28s %-15.2f %-17s\n",
			order.OrderCode,
			order.TotalPrice,
			order.PaymentStatus,
		)

		if order.PaymentStatus == "unpaid" {
			unpaidOrders = append(
				unpaidOrders,
				order,
			)
			totalPrice += order.TotalPrice
		}
	}

	fmt.Println(
		"-----------------------------------------------------------",
	)
	fmt.Printf(
		"Total Unpaid : Rp%.2f\n",
		totalPrice,
	)
	fmt.Println()

	if len(unpaidOrders) == 0 {
		fmt.Printf(
			"%s[SUCCESS]:%s Tidak ada tagihan yang harus dibayar.\n",
			colors.Green,
			colors.Reset,
		)
		pause("\nTekan Enter untuk kembali...")
		return
	}

	PayOrderCLI(
		paymentHandler,
		unpaidOrders,
	)
}

func pause(message string) {
	fmt.Print(message)
	bufio.NewReader(os.Stdin).ReadString('\n')
}
