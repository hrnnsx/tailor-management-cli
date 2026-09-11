package auth

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"tailor-management-cli/config/colors"
	"tailor-management-cli/entity"
	"tailor-management-cli/handler"

	"github.com/hrnnsx/go-toolkit/stdio"
	"github.com/manifoldco/promptui"
)

func SignIn(authhandler *handler.AuthHandler) (*entity.User, error) {
	stdio.ClearScreen()

	fmt.Println("                                          ")
	fmt.Println(">>>               SIGN IN             <<<")
	fmt.Println("                                          ")

	// ** Email PROMPT
	emailPrompt := promptui.Prompt{
		Label:       "Email",
		HideEntered: true,
	}

	var email string
	for {
		var err error
		email, err = emailPrompt.Run()

		if err != nil {
			return nil, err
		}

		email = strings.TrimSpace(email)

		if IsValidEmail(email) {
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
			return nil, err
		}

		if IsValidPassword(password) {
			break
		}

		fmt.Printf(
			"%s[WARNING] Password minimal 6 karakter%s\n",
			colors.Red,
			colors.Reset,
		)
	}
	fmt.Print("Password: ********\n")

	// ** Sign In
	user, err := authhandler.SignIn(email, password)
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
		return nil, err
	}

	fmt.Printf(
		"%s[SUCCESS]:%s Login berhasil!\n",
		colors.Green,
		colors.Reset,
	)

	fmt.Printf("Selamat datang, %s.\n", user.Name)
	fmt.Println()

	fmt.Print("Tekan Enter untuk melanjutkan...")
	bufio.NewReader(os.Stdin).ReadString('\n')

	// Clear terminal setelah user menekan Enter
	stdio.ClearScreen()

	return user, nil
}

func IsValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile(
		`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
	)

	return emailRegex.MatchString(email)
}

func IsValidPassword(password string) bool {
	return len(password) >= 6
}
