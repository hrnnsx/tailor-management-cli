package auth

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"tailor-management-cli/config/colors"
	"tailor-management-cli/handler"

	"github.com/hrnnsx/go-toolkit/stdio"
	"github.com/manifoldco/promptui"
)

func SignUpCLI(authHandler *handler.AuthHandler) {
	stdio.ClearScreen()

	fmt.Println("                                          ")
	fmt.Println(">>>               SIGN UP             <<<")
	fmt.Println("                                          ")

	// ** Name PROMPT
	namePrompt := promptui.Prompt{
		Label:       "Name",
		HideEntered: true,
	}

	var name string
	for {
		var err error
		name, err = namePrompt.Run()

		if err != nil {
			return
		}

		name = strings.TrimSpace(name)

		if name != "" {
			break
		}

		fmt.Printf(
			"%s[WARNING] Nama tidak boleh kosong%s\n",
			colors.Red,
			colors.Reset,
		)
	}
	fmt.Printf("Name: %s\n", name)

	// ** Email PROMPT
	emailRegex := regexp.MustCompile(
		`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`,
	)

	emailPrompt := promptui.Prompt{
		Label:       "Email",
		HideEntered: true,
	}

	var email string
	for {
		var err error
		email, err = emailPrompt.Run()

		if err != nil {
			return
		}

		email = strings.TrimSpace(email)

		if emailRegex.MatchString(email) {
			break
		}

		fmt.Printf(
			"%s[WARNING] Invalid Email (ex: johndoe@mail.com)%s\n",
			colors.Red,
			colors.Reset,
		)
	}
	fmt.Printf("Email: %s\n", email)

	// ** Password PROMPT
	passwordPrompt := promptui.Prompt{
		Label:       "Password",
		Mask:        '*',
		HideEntered: true,
	}

	var password string
	for {
		var err error
		password, err = passwordPrompt.Run()

		if err != nil {
			return
		}

		if len(password) >= 6 {
			break
		}

		fmt.Printf(
			"%s[WARNING] Password minimal 6 karakter%s\n",
			colors.Red,
			colors.Reset,
		)
	}
	fmt.Print("Password: ********\n")

	// ** Confirm Password PROMPT
	confirmPasswordPrompt := promptui.Prompt{
		Label:       "Confirm Password",
		Mask:        '*',
		HideEntered: true,
	}

	var confirmPassword string
	for {
		var err error
		confirmPassword, err = confirmPasswordPrompt.Run()

		if err != nil {
			return
		}

		if confirmPassword == password {
			break
		}

		fmt.Printf(
			"%s[WARNING] Konfirmasi password tidak cocok%s\n",
			colors.Red,
			colors.Reset,
		)
	}
	fmt.Print("Confirm Password: ********\n")

	// ** Phone PROMPT
	phonePrompt := promptui.Prompt{
		Label:       "Phone",
		HideEntered: true,
	}

	var phone string
	for {
		var err error
		phone, err = phonePrompt.Run()

		if err != nil {
			return
		}

		phone = strings.TrimSpace(phone)

		if strings.HasPrefix(phone, "08") && len(phone) >= 10 {
			break
		}

		fmt.Printf(
			"%s[WARNING] Nomor HP wajib diawali '08' dan minimal 10 digit%s\n",
			colors.Red,
			colors.Reset,
		)
	}
	fmt.Printf("Phone: %s\n", phone)

	// ** Register
	err := authHandler.Register(name, email, password, phone)
	if err != nil {
		fmt.Printf(
			"%s[FAILED]:%s %v\n",
			colors.Red,
			colors.Reset,
			err,
		)

		fmt.Print("\nTekan Enter untuk kembali...")
		bufio.NewReader(os.Stdin).ReadString('\n')

		stdio.ClearScreen()
		return
	}

	fmt.Printf(
		"%s[SUCCESS]:%s Registrasi akun berhasil!\n",
		colors.Green,
		colors.Reset,
	)
	fmt.Println("Silakan login untuk masuk ke sistem.")
	fmt.Println()

	fmt.Print("Tekan Enter untuk melanjutkan...")
	bufio.NewReader(os.Stdin).ReadString('\n')

	// Clear terminal setelah user menekan Enter
	stdio.ClearScreen()
}
