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
				"Login",
				"Signup",
				"Exit",
			},
		}
		idx, _, err := prompt.Run()
		if err != nil {
			return
		}

		switch idx {
		case 0:
			user := auth.LoginCLI(authH)
			if user != nil {
				// Arahkan ke menu customer jika role customer
				CustomerMenu(orderH, *user)
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
			fmt.Println("\n(Fitur cek status order dalam pengerjaan)")
		case 2:
			fmt.Println("\n(Fitur cek status pembayaran dalam pengerjaan)")
		case 3:
			fmt.Println("\nBerhasil logout.")
			return
		}
	}
}
