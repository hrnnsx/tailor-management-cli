package cli

import (
	"bufio"
	"fmt"
	"os"

	"tailor-management-cli/config/colors"
	"tailor-management-cli/handler"

	"github.com/hrnnsx/go-toolkit/stdio"
	"github.com/manifoldco/promptui"
)

var workerSelectTemplates = &promptui.SelectTemplates{
	Active:   "▸ {{ . | cyan }}",
	Inactive: "  {{ . }}",
	Selected: "✔ {{ . | green }}",
}

func workerPause(message string) {
	fmt.Print(message)
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func workerReturn() {
	workerPause("\nTekan ENTER untuk kembali...")
	stdio.ClearScreen()
}

func ShowMyOrders(
	workerH *handler.WorkerHandler,
	workerID int,
) {
	stdio.ClearScreen()

	fmt.Println()
	fmt.Println("==========================================================================")
	fmt.Println("                         ORDER SAYA")
	fmt.Println("==========================================================================")
	fmt.Println()

	orders, err := workerH.GetMyOrders(workerID)
	if err != nil {
		fmt.Printf("%s[ERROR]%s %v\n", colors.Red, colors.Reset, err)
		workerReturn()
		return
	}

	if len(orders) == 0 {
		fmt.Printf("%s[WARNING]%s Tidak ada order yang sedang dikerjakan.\n",
			colors.Red,
			colors.Reset,
		)
		workerReturn()
		return
	}

	fmt.Println("Daftar order yang sedang dikerjakan:")
	fmt.Println()

	fmt.Println("==========================================================================")
	fmt.Printf(
		"%-4s %-18s %-18s %-8s %-8s %-22s\n",
		"ID",
		"ORDER CODE",
		"CUSTOMER",
		"SIZE",
		"CM",
		"PROGRESS",
	)
	fmt.Println("==========================================================================")

	for _, order := range orders {
		fmt.Printf(
			"%-4d %-18s %-18s %-8s %-8d %-22s\n",
			order.ID,
			order.OrderCode,
			order.CustomerName,
			order.DeterminedSize,
			order.CMUsed,
			order.Progress,
		)
	}

	fmt.Println("==========================================================================")

	workerReturn()
}

func UpdateOrderProgressCLI(
	workerH *handler.WorkerHandler,
	workerID int,
) {
	stdio.ClearScreen()

	fmt.Println()
	fmt.Println("==========================================================================")
	fmt.Println("                       UPDATE PROGRESS")
	fmt.Println("==========================================================================")
	fmt.Println()

	orders, err := workerH.GetMyOrders(workerID)
	if err != nil {
		fmt.Printf("%s[ERROR]%s %v\n", colors.Red, colors.Reset, err)
		workerReturn()
		return
	}

	if len(orders) == 0 {
		fmt.Printf("%s[WARNING]%s Tidak ada order yang sedang dikerjakan.\n",
			colors.Red,
			colors.Reset,
		)
		workerReturn()
		return
	}

	items := make([]string, len(orders))

	for i, order := range orders {
		items[i] = fmt.Sprintf(
			"%s | %s | %s",
			order.OrderCode,
			order.CustomerName,
			order.Progress,
		)
	}

	orderPrompt := promptui.Select{
		Label:     "Pilih order yang ingin di-update",
		Items:     items,
		Templates: workerSelectTemplates,
	}

	orderIndex, _, err := orderPrompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	selectedOrder := orders[orderIndex]

	fmt.Println()
	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("                         DETAIL ORDER")
	fmt.Println("--------------------------------------------------------------------------")
	fmt.Printf("Order Code : %s\n", selectedOrder.OrderCode)
	fmt.Printf("Customer   : %s\n", selectedOrder.CustomerName)
	fmt.Printf("Size       : %s\n", selectedOrder.DeterminedSize)
	fmt.Printf("Progress   : %s\n", selectedOrder.Progress)
	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println()

	progressItems := []string{
		"Order diterima",
		"Persiapan bahan",
		"Pemotongan kain",
		"Proses jahit",
		"Finishing",
	}

	progressPrompt := promptui.Select{
		Label:     "Update progress ke",
		Items:     progressItems,
		Templates: workerSelectTemplates,
	}

	progressIndex, _, err := progressPrompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	selectedProgress := progressItems[progressIndex]

	fmt.Println()
	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("                       KONFIRMASI UPDATE")
	fmt.Println("--------------------------------------------------------------------------")
	fmt.Printf("Order    : %s\n", selectedOrder.OrderCode)
	fmt.Printf("Sebelum  : %s\n", selectedOrder.Progress)
	fmt.Printf("Menjadi  : %s\n", selectedProgress)
	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println()

	confirmPrompt := promptui.Select{
		Label: "Update progress?",
		Items: []string{
			"Ya, update",
			"Batal",
		},
		Templates: workerSelectTemplates,
	}

	confirmIndex, _, err := confirmPrompt.Run()
	if err != nil {
		stdio.ClearScreen()
		return
	}

	if confirmIndex == 1 {
		fmt.Println()
		fmt.Printf("%s[WARNING]%s Update progress dibatalkan.\n",
			colors.Red,
			colors.Reset,
		)
		workerReturn()
		return
	}

	err = workerH.UpdateProgress(
		workerID,
		selectedOrder.ID,
		selectedProgress,
	)

	if err != nil {
		fmt.Printf("\n%s[ERROR]%s %v\n",
			colors.Red,
			colors.Reset,
			err,
		)
		workerReturn()
		return
	}

	fmt.Println()
	fmt.Println("==========================================================================")
	fmt.Println("                    PROGRESS BERHASIL DIPERBARUI")
	fmt.Println("==========================================================================")
	fmt.Printf("Order    : %s\n", selectedOrder.OrderCode)
	fmt.Printf("Progress : %s\n", selectedProgress)
	fmt.Println("==========================================================================")

	if selectedProgress == "Finishing" {
		fmt.Println()
		fmt.Println("[INFO] Order sudah masuk tahap finishing.")
		fmt.Println("[INFO] Status order akan menyesuaikan pembayaran.")
	}

	workerReturn()
}
