package cli

import (
	"bufio"
	"fmt"
	"os"

	"tailor-management-cli/config/colors"
	"tailor-management-cli/entity"
	"tailor-management-cli/handler"

	"github.com/hrnnsx/go-toolkit/stdio"
	"github.com/manifoldco/promptui"
)

var adminSelectTemplates = &promptui.SelectTemplates{
	Active:   "▸ {{ . | cyan }}",
	Inactive: "  {{ . }}",
	Selected: "✔ {{ . | green }}",
}

func adminPause(message string) {
	fmt.Print(message)
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func adminReturn() {
	adminPause("\nTekan ENTER untuk kembali...")
	stdio.ClearScreen()
}

// ============================================================
// CHECK ALL ORDER
// ============================================================

func CheckAllOrder(orderH *handler.OrderHandler) {
	stdio.ClearScreen()

	fmt.Println("==============================================================")
	fmt.Println("                    CEK & ASSIGN ORDER")
	fmt.Println("==============================================================")
	fmt.Println()

	orders, err := orderH.CheckAllOrder()
	if err != nil {
		fmt.Printf("%s[ERROR]%s %v\n", colors.Red, colors.Reset, err)
		adminReturn()
		return
	}

	if len(orders) == 0 {
		fmt.Println("[WARNING] Tidak ada order yang sedang menunggu assignment.")
		adminReturn()
		return
	}

	fmt.Println("Order yang menunggu assignment:")
	fmt.Println()

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Printf(
		"%-4s %-25s %-16s %-8s %-10s %-20s\n",
		"ID",
		"ORDER CODE",
		"CUSTOMER",
		"SIZE",
		"BAHAN",
		"CREATED AT",
	)
	fmt.Println("--------------------------------------------------------------------------")

	for _, order := range orders {
		fmt.Printf(
			"%-4d %-25s %-16s %-8s %-10d %-20s\n",
			order.ID,
			order.OrderCode,
			order.CustomerName,
			order.DeterminedSize,
			order.CMUsed,
			order.CreatedAt.Format("02-01-2006 15:04"),
		)
	}

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println()

	orderItems := make([]string, len(orders))

	for i, order := range orders {
		orderItems[i] = fmt.Sprintf(
			"%s | %s | Size %s | %d cm",
			order.OrderCode,
			order.CustomerName,
			order.DeterminedSize,
			order.CMUsed,
		)
	}

	orderPrompt := promptui.Select{
		Label:     "Pilih order",
		Items:     orderItems,
		Templates: adminSelectTemplates,
	}

	orderIndex, _, err := orderPrompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	selectedOrder := orders[orderIndex]

	fmt.Println()
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("                      ORDER TERPILIH")
	fmt.Println("--------------------------------------------------------------")
	fmt.Printf("Order Code : %s\n", selectedOrder.OrderCode)
	fmt.Printf("Customer   : %s\n", selectedOrder.CustomerName)
	fmt.Printf("Size       : %s\n", selectedOrder.DeterminedSize)
	fmt.Printf("Bahan      : %d cm\n", selectedOrder.CMUsed)
	fmt.Println("--------------------------------------------------------------")
	fmt.Println()

	workers, err := orderH.GetAvailableWorkers()
	if err != nil {
		fmt.Printf("%s[ERROR]%s %v\n", colors.Red, colors.Reset, err)
		adminReturn()
		return
	}

	if len(workers) == 0 {
		fmt.Println("[WARNING] Tidak ada worker yang sedang tersedia.")
		adminReturn()
		return
	}

	workerItems := make([]string, len(workers))

	for i, worker := range workers {
		workerItems[i] = fmt.Sprintf(
			"%s (ID: %d)",
			worker.Name,
			worker.ID,
		)
	}

	workerPrompt := promptui.Select{
		Label:     "Pilih worker",
		Items:     workerItems,
		Templates: adminSelectTemplates,
	}

	workerIndex, _, err := workerPrompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	selectedWorker := workers[workerIndex]

	fmt.Println()
	fmt.Println("==============================================================")
	fmt.Println("                  KONFIRMASI ASSIGNMENT")
	fmt.Println("==============================================================")
	fmt.Println()
	fmt.Printf("Order  : %s\n", selectedOrder.OrderCode)
	fmt.Printf("Worker : %s\n", selectedWorker.Name)
	fmt.Println()

	confirmPrompt := promptui.Select{
		Label: "Assign order ke worker ini?",
		Items: []string{
			"Ya, assign order",
			"Batal",
		},
		Templates: adminSelectTemplates,
	}

	confirmIndex, _, err := confirmPrompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	if confirmIndex == 1 {
		fmt.Println()
		fmt.Println("[WARNING] Assignment dibatalkan.")
		adminReturn()
		return
	}

	err = orderH.AssignOrder(
		selectedOrder.ID,
		selectedWorker.ID,
	)

	if err != nil {
		fmt.Printf("\n%s[ERROR]%s %v\n", colors.Red, colors.Reset, err)
		adminReturn()
		return
	}

	fmt.Println()
	fmt.Println("==============================================================")
	fmt.Println("                  ASSIGNMENT BERHASIL")
	fmt.Println("==============================================================")
	fmt.Println()
	fmt.Printf("Order  : %s\n", selectedOrder.OrderCode)
	fmt.Printf("Worker : %s\n", selectedWorker.Name)
	fmt.Println()
	fmt.Println("==============================================================")

	adminReturn()
}

// ============================================================
// CHECK PAYMENT
// ============================================================

func CheckPayment(paymentH *handler.PaymentHandler, user entity.User) {
	stdio.ClearScreen()

	fmt.Println("==============================================================")
	fmt.Println("                CEK & VERIFIKASI PEMBAYARAN")
	fmt.Println("==============================================================")
	fmt.Println()

	payments, err := paymentH.GetPendingPayments()
	if err != nil {
		fmt.Printf("%s[ERROR]%s %v\n", colors.Red, colors.Reset, err)
		adminReturn()
		return
	}

	if len(payments) == 0 {
		fmt.Println("[WARNING] Tidak ada pembayaran yang perlu diverifikasi.")
		adminReturn()
		return
	}

	fmt.Println("Pembayaran yang menunggu verifikasi:")
	fmt.Println()

	items := make([]string, len(payments))

	for i, payment := range payments {
		items[i] = fmt.Sprintf(
			"%s | %s | Rp %.2f | %s",
			payment.OrderCode,
			payment.CustomerName,
			payment.Amount,
			payment.Status,
		)
	}

	prompt := promptui.Select{
		Label:     "Pilih pembayaran",
		Items:     items,
		Templates: adminSelectTemplates,
	}

	index, _, err := prompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	selectedPayment := payments[index]

	fmt.Println()
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("                    DETAIL PEMBAYARAN")
	fmt.Println("--------------------------------------------------------------")
	fmt.Printf("Order Code : %s\n", selectedPayment.OrderCode)
	fmt.Printf("Customer   : %s\n", selectedPayment.CustomerName)
	fmt.Printf("Amount     : Rp %.2f\n", selectedPayment.Amount)
	fmt.Printf("Status     : %s\n", selectedPayment.Status)
	fmt.Println("--------------------------------------------------------------")
	fmt.Println()

	confirmPrompt := promptui.Select{
		Label: "Verifikasi pembayaran ini?",
		Items: []string{
			"Ya, verifikasi",
			"Batal",
		},
		Templates: adminSelectTemplates,
	}

	confirmIndex, _, err := confirmPrompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	if confirmIndex == 1 {
		fmt.Println()
		fmt.Println("[WARNING] Verifikasi pembayaran dibatalkan.")
		adminReturn()
		return
	}

	err = paymentH.VerifyPayment(
		selectedPayment.ID,
		user.ID,
	)

	if err != nil {
		fmt.Printf("\n%s[ERROR]%s %v\n", colors.Red, colors.Reset, err)
		adminReturn()
		return
	}

	fmt.Println()
	fmt.Println("==============================================================")
	fmt.Println("                 PEMBAYARAN BERHASIL")
	fmt.Println("==============================================================")
	fmt.Println()
	fmt.Printf("Order Code : %s\n", selectedPayment.OrderCode)
	fmt.Printf("Customer   : %s\n", selectedPayment.CustomerName)
	fmt.Println("Status     : paid")
	fmt.Println()
	fmt.Println("==============================================================")

	adminReturn()
}

// ============================================================
// CHECK REPORT
// ============================================================

func CheckReport(orderH *handler.OrderHandler) {
	stdio.ClearScreen()

	fmt.Println("==============================================================")
	fmt.Println("                    LAPORAN PENJUALAN")
	fmt.Println("==============================================================")
	fmt.Println()

	report, err := orderH.GetSalesReport()
	if err != nil {
		fmt.Printf("%s[ERROR]%s %v\n", colors.Red, colors.Reset, err)
		adminReturn()
		return
	}

	fmt.Println("Ringkasan penjualan:")
	fmt.Println()

	fmt.Println("--------------------------------------------------------------")
	fmt.Printf("%-25s : %d\n", "Total Orders", report.TotalOrders)
	fmt.Printf("%-25s : %d\n", "Finished Orders", report.FinishedOrders)
	fmt.Printf("%-25s : %d\n", "Paid Orders", report.PaidOrders)
	fmt.Printf("%-25s : %d\n", "Unpaid Orders", report.UnpaidOrders)
	fmt.Printf("%-25s : Rp %.2f\n", "Total Revenue", report.TotalRevenue)
	fmt.Println("--------------------------------------------------------------")

	adminReturn()
}

// ============================================================
// CREATE FABRIC
// ============================================================

func CreateFabricCLI(fabricH *handler.FabricHandler) {
	stdio.ClearScreen()

	fmt.Println("==============================================================")
	fmt.Println("                    TAMBAH JENIS KAIN")
	fmt.Println("==============================================================")
	fmt.Println()

	namePrompt := promptui.Prompt{
		Label: "Nama kain",
	}

	name, err := namePrompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	fmt.Println()

	pricePrompt := promptui.Prompt{
		Label: "Harga per cm",
	}

	priceInput, err := pricePrompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	var price float64

	_, err = fmt.Sscanf(priceInput, "%f", &price)
	if err != nil {
		fmt.Printf(
			"\n%s[ERROR]%s Harga harus berupa angka.\n",
			colors.Red,
			colors.Reset,
		)
		adminReturn()
		return
	}

	if price <= 0 {
		fmt.Printf(
			"\n%s[ERROR]%s Harga harus lebih dari 0.\n",
			colors.Red,
			colors.Reset,
		)
		adminReturn()
		return
	}

	fmt.Println()
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("                     RINGKASAN KAIN")
	fmt.Println("--------------------------------------------------------------")
	fmt.Printf("Nama kain  : %s\n", name)
	fmt.Printf("Harga / cm : Rp %.2f\n", price)
	fmt.Println("--------------------------------------------------------------")
	fmt.Println()

	confirmPrompt := promptui.Select{
		Label: "Tambahkan jenis kain ini?",
		Items: []string{
			"Ya, tambahkan",
			"Batal",
		},
		Templates: adminSelectTemplates,
	}

	confirmIndex, _, err := confirmPrompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	if confirmIndex == 1 {
		fmt.Println()
		fmt.Println("[WARNING] Penambahan kain dibatalkan.")
		adminReturn()
		return
	}

	err = fabricH.CreateFabric(name, price)
	if err != nil {
		fmt.Printf(
			"\n%s[ERROR]%s %v\n",
			colors.Red,
			colors.Reset,
			err,
		)
		adminReturn()
		return
	}

	fmt.Println()
	fmt.Println("==============================================================")
	fmt.Println("                 KAIN BERHASIL DITAMBAHKAN")
	fmt.Println("==============================================================")
	fmt.Println()
	fmt.Printf("Nama kain  : %s\n", name)
	fmt.Printf("Harga / cm : Rp %.2f\n", price)
	fmt.Println()
	fmt.Println("==============================================================")

	adminReturn()
}

// ============================================================
// RESTOCK FABRIC
// ============================================================

func RestockFabricCLI(fabricH *handler.FabricHandler) {
	stdio.ClearScreen()

	fmt.Println("==============================================================")
	fmt.Println("                       RESTOCK KAIN")
	fmt.Println("==============================================================")
	fmt.Println()

	fabrics, err := fabricH.GetFabricPatterns()
	if err != nil {
		fmt.Printf("%s[ERROR]%s %v\n", colors.Red, colors.Reset, err)
		adminReturn()
		return
	}

	if len(fabrics) == 0 {
		fmt.Println("[WARNING] Belum ada kombinasi kain dan pattern.")
		adminReturn()
		return
	}

	fmt.Println("Daftar kombinasi kain dan pattern:")
	fmt.Println()

	fmt.Println("--------------------------------------------------------------")
	fmt.Printf(
		"%-4s %-20s %-20s %-10s\n",
		"ID",
		"FABRIC",
		"PATTERN",
		"STOCK",
	)
	fmt.Println("--------------------------------------------------------------")

	for _, fabric := range fabrics {
		fmt.Printf(
			"%-4d %-20s %-20s %d cm\n",
			fabric.ID,
			fabric.FabricName,
			fabric.PatternName,
			fabric.StockCM,
		)
	}

	fmt.Println("--------------------------------------------------------------")
	fmt.Println()

	prompt := promptui.Prompt{
		Label: "Masukkan ID kombinasi kain",
	}

	input, err := prompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	var fabricPatternID int

	_, err = fmt.Sscanf(input, "%d", &fabricPatternID)
	if err != nil {
		fmt.Printf(
			"\n%s[ERROR]%s ID harus berupa angka.\n",
			colors.Red,
			colors.Reset,
		)
		adminReturn()
		return
	}

	var selectedFabric *entity.FabricPattern

	for i := range fabrics {
		if fabrics[i].ID == fabricPatternID {
			selectedFabric = &fabrics[i]
			break
		}
	}

	if selectedFabric == nil {
		fmt.Printf(
			"\n%s[ERROR]%s Kombinasi kain tidak ditemukan.\n",
			colors.Red,
			colors.Reset,
		)
		adminReturn()
		return
	}

	fmt.Println()
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("                    KAIN TERPILIH")
	fmt.Println("--------------------------------------------------------------")
	fmt.Printf("Kain    : %s\n", selectedFabric.FabricName)
	fmt.Printf("Pattern : %s\n", selectedFabric.PatternName)
	fmt.Printf("Stok    : %d cm\n", selectedFabric.StockCM)
	fmt.Println("--------------------------------------------------------------")
	fmt.Println()

	quantityPrompt := promptui.Prompt{
		Label: "Jumlah restock (cm)",
	}

	quantityInput, err := quantityPrompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	var quantityCM int

	_, err = fmt.Sscanf(quantityInput, "%d", &quantityCM)
	if err != nil {
		fmt.Printf(
			"\n%s[ERROR]%s Jumlah harus berupa angka.\n",
			colors.Red,
			colors.Reset,
		)
		adminReturn()
		return
	}

	if quantityCM <= 0 {
		fmt.Printf(
			"\n%s[ERROR]%s Jumlah restock harus lebih dari 0 cm.\n",
			colors.Red,
			colors.Reset,
		)
		adminReturn()
		return
	}

	fmt.Println()
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("                    RINGKASAN RESTOCK")
	fmt.Println("--------------------------------------------------------------")
	fmt.Printf("Kain            : %s\n", selectedFabric.FabricName)
	fmt.Printf("Pattern         : %s\n", selectedFabric.PatternName)
	fmt.Printf("Stok sebelumnya : %d cm\n", selectedFabric.StockCM)
	fmt.Printf("Jumlah restock  : %d cm\n", quantityCM)
	fmt.Printf("Stok sekarang   : %d cm\n", selectedFabric.StockCM+quantityCM)
	fmt.Println("--------------------------------------------------------------")
	fmt.Println()

	confirmPrompt := promptui.Select{
		Label: "Tambahkan stok ini?",
		Items: []string{
			"Ya, tambahkan",
			"Batal",
		},
		Templates: adminSelectTemplates,
	}

	confirmIndex, _, err := confirmPrompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	if confirmIndex == 1 {
		fmt.Println()
		fmt.Println("[WARNING] Restock dibatalkan.")
		adminReturn()
		return
	}

	err = fabricH.RestockFabric(
		fabricPatternID,
		quantityCM,
	)

	if err != nil {
		fmt.Printf(
			"\n%s[ERROR]%s %v\n",
			colors.Red,
			colors.Reset,
			err,
		)
		adminReturn()
		return
	}

	fmt.Println()
	fmt.Println("==============================================================")
	fmt.Println("                    RESTOCK BERHASIL")
	fmt.Println("==============================================================")
	fmt.Println()
	fmt.Printf("Kain            : %s\n", selectedFabric.FabricName)
	fmt.Printf("Pattern         : %s\n", selectedFabric.PatternName)
	fmt.Printf("Stok sebelumnya : %d cm\n", selectedFabric.StockCM)
	fmt.Printf("Restock         : %d cm\n", quantityCM)
	fmt.Printf("Stok sekarang   : %d cm\n", selectedFabric.StockCM+quantityCM)
	fmt.Println()
	fmt.Println("==============================================================")

	adminReturn()
}

// ============================================================
// CREATE FABRIC PATTERN
// ============================================================

func CreateFabricPatternCLI(fabricH *handler.FabricHandler) {
	stdio.ClearScreen()

	fmt.Println("==============================================================")
	fmt.Println("                 TAMBAH KOMBINASI KAIN")
	fmt.Println("==============================================================")
	fmt.Println()

	fabrics, err := fabricH.GetFabrics()
	if err != nil {
		fmt.Printf("%s[ERROR]%s %v\n", colors.Red, colors.Reset, err)
		adminReturn()
		return
	}

	if len(fabrics) == 0 {
		fmt.Println("[WARNING] Belum ada jenis kain.")
		adminReturn()
		return
	}

	fabricItems := make([]string, len(fabrics))

	for i, fabric := range fabrics {
		fabricItems[i] = fmt.Sprintf(
			"%s - Rp %.2f/cm",
			fabric.Name,
			fabric.PricePerCM,
		)
	}

	fabricPrompt := promptui.Select{
		Label:     "Pilih kain",
		Items:     fabricItems,
		Templates: adminSelectTemplates,
	}

	fabricIndex, _, err := fabricPrompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	selectedFabric := fabrics[fabricIndex]

	fmt.Println()
	fmt.Printf("Kain terpilih : %s\n", selectedFabric.Name)
	fmt.Println()

	patterns, err := fabricH.GetPatterns()
	if err != nil {
		fmt.Printf("%s[ERROR]%s %v\n", colors.Red, colors.Reset, err)
		adminReturn()
		return
	}

	if len(patterns) == 0 {
		fmt.Println("[WARNING] Belum ada pattern.")
		adminReturn()
		return
	}

	patternItems := make([]string, len(patterns))

	for i, pattern := range patterns {
		patternItems[i] = pattern.Name
	}

	patternPrompt := promptui.Select{
		Label:     "Pilih pattern",
		Items:     patternItems,
		Templates: adminSelectTemplates,
	}

	patternIndex, _, err := patternPrompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	selectedPattern := patterns[patternIndex]

	fmt.Println()
	fmt.Printf("Pattern terpilih : %s\n", selectedPattern.Name)
	fmt.Println()

	stockPrompt := promptui.Prompt{
		Label: "Stok awal (cm)",
	}

	stockInput, err := stockPrompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	var stockCM int

	_, err = fmt.Sscanf(stockInput, "%d", &stockCM)
	if err != nil {
		fmt.Printf(
			"\n%s[ERROR]%s Stok harus berupa angka.\n",
			colors.Red,
			colors.Reset,
		)
		adminReturn()
		return
	}

	if stockCM < 0 {
		fmt.Printf(
			"\n%s[ERROR]%s Stok tidak boleh negatif.\n",
			colors.Red,
			colors.Reset,
		)
		adminReturn()
		return
	}

	fmt.Println()
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("                  RINGKASAN KOMBINASI")
	fmt.Println("--------------------------------------------------------------")
	fmt.Printf("Kain      : %s\n", selectedFabric.Name)
	fmt.Printf("Pattern   : %s\n", selectedPattern.Name)
	fmt.Printf("Stok awal : %d cm\n", stockCM)
	fmt.Println("--------------------------------------------------------------")
	fmt.Println()

	confirmPrompt := promptui.Select{
		Label: "Tambahkan kombinasi ini?",
		Items: []string{
			"Ya, tambahkan",
			"Batal",
		},
		Templates: adminSelectTemplates,
	}

	confirmIndex, _, err := confirmPrompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	if confirmIndex == 1 {
		fmt.Println()
		fmt.Println("[WARNING] Penambahan kombinasi dibatalkan.")
		adminReturn()
		return
	}

	err = fabricH.CreateFabricPattern(
		selectedFabric.ID,
		selectedPattern.ID,
		stockCM,
	)

	if err != nil {
		fmt.Printf(
			"\n%s[ERROR]%s %v\n",
			colors.Red,
			colors.Reset,
			err,
		)
		adminReturn()
		return
	}

	fmt.Println()
	fmt.Println("==============================================================")
	fmt.Println("             KOMBINASI BERHASIL DITAMBAHKAN")
	fmt.Println("==============================================================")
	fmt.Println()
	fmt.Printf("Kain      : %s\n", selectedFabric.Name)
	fmt.Printf("Pattern   : %s\n", selectedPattern.Name)
	fmt.Printf("Stok awal : %d cm\n", stockCM)
	fmt.Println()
	fmt.Println("==============================================================")

	adminReturn()
}
