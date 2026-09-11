package cli

import (
	"database/sql"
	"fmt"

	"tailor-management-cli/cli/auth"
	"tailor-management-cli/entity"
	"tailor-management-cli/handler"

	"github.com/hrnnsx/go-toolkit/stdio"
	"github.com/manifoldco/promptui"
)

func MainMenu(db *sql.DB) {
	authH := handler.NewAuthHandler(db)
	orderH := handler.NewOrderHandler(db)
	paymentH := handler.NewPaymentHandler(db)
	fabricH := handler.NewFabricHandler(db)
	workerH := handler.NewWorkerHandler(db)

	for {
		stdio.ClearScreen()

		fmt.Println("                                          ")
		fmt.Println(">>>             MAIN MENU             <<<")
		fmt.Println("                                          ")

		prompt := promptui.Select{
			Label: "Pilih Menu Utama",
			Items: []string{
				"Sign In",
				"Sign Up",
				"Exit",
			},
			Templates: &promptui.SelectTemplates{
				Active:   "▸ {{ . | cyan }}",
				Inactive: "  {{ . }}",
				Selected: "✔ {{ . | green }}",
			},
		}

		idx, _, err := prompt.Run()
		if err != nil {
			return
		}

		switch idx {
		case 0:
			user, err := auth.SignIn(authH)
			if err != nil {
				continue
			}

			switch user.Role {
			case "customer":
				CustomerMenu(orderH, paymentH, *user)
			case "admin":
				AdminMenu(paymentH, orderH, fabricH, *user)
			case "worker":
				WorkerMenu(workerH, *user)
			}

		case 1:
			auth.SignUpCLI(authH)

		case 2:
			stdio.ClearScreen()

			fmt.Println()
			fmt.Println("Terima kasih, program ditutup.")
			return
		}
	}
}

func CustomerMenu(
	orderH *handler.OrderHandler,
	paymentH *handler.PaymentHandler,
	user entity.User,
) {
	for {
		stdio.ClearScreen()

		fmt.Println("                                          ")
		fmt.Println(">>>           CUSTOMER MENU           <<<")
		fmt.Println("                                          ")

		prompt := promptui.Select{
			Label: fmt.Sprintf(
				"Hi, %s, ada yang bisa kami bantu?",
				user.Name,
			),
			Items: []string{
				"Bikin baju",
				"Cek status order yang sedang berjalan",
				"Cek & Bayar Tagihan",
				"Log out",
			},
			Templates: &promptui.SelectTemplates{
				Active:   "▸ {{ . | cyan }}",
				Inactive: "  {{ . }}",
				Selected: "✔ {{ . | green }}",
			},
		}

		idx, _, err := prompt.Run()
		if err != nil {
			return
		}

		switch idx {
		case 0:
			BikinBajuCLI(orderH, user)

		case 1:
			CheckAllOrderStatus(orderH, user.ID)

		case 2:
			CheckBill(orderH, paymentH, user.ID)

		case 3:
			stdio.ClearScreen()

			fmt.Println()
			fmt.Println("Berhasil logout.")
			return
		}
	}
}

func FabricMenu(fabricH *handler.FabricHandler) {
	for {
		stdio.ClearScreen()

		fmt.Println("                                          ")
		fmt.Println(">>>            FABRIC MENU            <<<")
		fmt.Println("                                          ")

		prompt := promptui.Select{
			Label: "Kelola Kain",
			Items: []string{
				"Restock Kain",
				"Tambah Jenis Kain",
				"Tambah Kombinasi Kain",
				"Log Out",
			},
			Templates: &promptui.SelectTemplates{
				Active:   "▸ {{ . | cyan }}",
				Inactive: "  {{ . }}",
				Selected: "✔ {{ . | green }}",
			},
		}

		idx, _, err := prompt.Run()
		if err != nil {
			return
		}

		switch idx {
		case 0:
			RestockFabricCLI(fabricH)

		case 1:
			CreateFabricCLI(fabricH)

		case 2:
			CreateFabricPatternCLI(fabricH)

		case 3:
			stdio.ClearScreen()

			fmt.Println()
			fmt.Println("Kembali ke menu admin.")
			return
		}
	}
}

func AdminMenu(
	paymentH *handler.PaymentHandler,
	orderH *handler.OrderHandler,
	fabricH *handler.FabricHandler,
	user entity.User,
) {
	fmt.Println("                                          ")
	fmt.Println(">>>             ADMIN MENU             <<<")
	fmt.Println("                                          ")
	for {

		prompt := promptui.Select{
			Label: fmt.Sprintf(
				"Hi, %s, ada apa yang perlu dilakukan saat ini?",
				user.Name,
			),
			Items: []string{
				"Cek & Assign Order",
				"Cek Report Penjualan",
				"Cek & Verifikasi Pembayaran",
				"Check & Restock Kain",
				"Log Out",
			},
			Templates: &promptui.SelectTemplates{
				Active:   "▸ {{ . | cyan }}",
				Inactive: "  {{ . }}",
				Selected: "✔ {{ . | green }}",
			},
		}

		idx, _, err := prompt.Run()
		if err != nil {
			return
		}

		switch idx {
		case 0:
			CheckAllOrder(orderH)

		case 1:
			CheckReport(orderH)

		case 2:
			CheckPayment(paymentH, user)

		case 3:
			FabricMenu(fabricH)

		case 4:
			stdio.ClearScreen()

			fmt.Println()
			fmt.Println("Berhasil logout.")
			return
		}
	}
}

func WorkerMenu(workerH *handler.WorkerHandler, user entity.User) {
	for {
		fmt.Println()
		fmt.Println(">>>             WORKER MENU             <<<")
		fmt.Println()

		prompt := promptui.Select{
			Label: fmt.Sprintf(
				"Hi, %s, ada pekerjaan apa?",
				user.Name,
			),
			Items: []string{
				"Lihat Order Saya",
				"Update Progress Order",
				"Log out",
			},
			Templates: workerSelectTemplates,
		}

		index, _, err := prompt.Run()
		if err != nil {
			stdio.ClearScreen()
			return
		}

		switch index {
		case 0:
			ShowMyOrders(workerH, user.ID)

		case 1:
			UpdateOrderProgressCLI(workerH, user.ID)

		case 2:
			stdio.ClearScreen()
			fmt.Println()
			fmt.Println("[SUCCESS] Berhasil logout.")
			return
		}
	}
}
