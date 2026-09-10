package cli

import (
	"database/sql"
	"fmt"

	"tailor-management-cli/cli/auth"
	"tailor-management-cli/entity"
	"tailor-management-cli/handler"

	"github.com/manifoldco/promptui"
)

func MainMenu(db *sql.DB) {
	authH := handler.NewAuthHandler(db)
	orderH := handler.NewOrderHandler(db)

	for {
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
				fmt.Println(err)
				continue
			}
			switch user.Role {
			case "customer":
				CustomerMenu(orderH, *user)
			case "admin":
			case "worker":
			}
		case 1:
			auth.SignUpCLI(authH)
		case 2:
			fmt.Println("Terima kasih, program ditutup.")
			return
		}
	}
}

func CustomerMenu(orderH *handler.OrderHandler, user entity.User) {
	for {
		prompt := promptui.Select{
			Label: fmt.Sprintf("Hi, %s, ada yang bisa kami bantu?", user.Name),
			Items: []string{
				"Bikin baju",
				"Cek status order yang sedang berjalan",
				"Cek status pembayaran",
				"Log out",
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
			CheckAllOrder(orderH, user.ID)
		case 2:
			fmt.Println("\n(Fitur cek status pembayaran dalam pengerjaan)")
		case 3:
			fmt.Println("\nBerhasil logout.")
			return
		}
	}
}
