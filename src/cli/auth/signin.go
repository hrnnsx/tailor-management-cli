package auth

import (
	"fmt"
	"regexp"
	"tailor-management-cli/config/colors"
	"tailor-management-cli/entity"
	"tailor-management-cli/handler"

	"github.com/hrnnsx/go-toolkit/stdio"
	"github.com/manifoldco/promptui"
)

func SignIn(authhandler *handler.AuthHandler) (*entity.User, error) {
	stdio.ClearScreen()
	fmt.Println("                                          ")
	fmt.Println("                  SIGN IN                 ")
	fmt.Println("                                          ")

	// ** Email PROMPT
	emailPrompt := promptui.Prompt{
		Label:       "Email",
		HideEntered: true,
	}

	var email string
	for {
		var emailErr error
		email, emailErr = emailPrompt.Run()

		if emailErr != nil {
			return nil, emailErr
		}
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
		var passwordErr error
		password, passwordErr = passwordPrompt.Run()

		if passwordErr != nil {
			return nil, passwordErr
		}
		if IsValidPassword(password) {
			break
		}

		fmt.Printf("%s[WARNING] Password minimal 6 karakter%s\n", colors.Red, colors.Reset)
	}

	user, err := authhandler.SignIn(email, password)
	if err != nil {
		return nil, fmt.Errorf("%s[FAILED]:%s %w", colors.Red, colors.Reset, err)
	}

	// Clear terminal
	stdio.ClearScreen()
	return user, nil
}

func IsValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidPassword(password string) bool {
	return len(password) >= 6
}
