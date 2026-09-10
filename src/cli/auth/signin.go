package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"tailor-management-cli/entity"

	"github.com/hrnnsx/go-toolkit/stdio"
	"github.com/manifoldco/promptui"
)

func SignIn(db *sql.DB) (entity.User, error) {
	fmt.Println("MASUK SIGN IN")

	// ** Email PROMPT
	emailPrompt := promptui.Prompt{
		Label:       "Email",
		HideEntered: true,
	}

	var email string
	var emailErr error

	for {
		email, emailErr = emailPrompt.Run()
		if emailErr != nil {
			return entity.User{}, fmt.Errorf("Error email input: %v", emailErr)
		}
		if IsValidEmail(email) {
			break
		} else {
			fmt.Println("Invalid Email input...")
		}
	}

	fmt.Printf("Email: %s\n", email)

	// ** Password PROMPT
	passwordPrompt := promptui.Prompt{
		Label:       "Password",
		Mask:        '*',
		HideEntered: true,
	}
	password, passwordErr := passwordPrompt.Run()

	if passwordErr != nil {
		return entity.User{}, fmt.Errorf("Error input password: %v", passwordErr)
	}

	// Clear terminal
	stdio.ClearScreen()
	fmt.Printf("\nSIGNING IN...\n\n")

	// Backend PROC
	// TODO: Hashing password

	var user entity.User

	query := `SELECT id, name, email, password, role, phone, created_at FROM users WHERE email = ?`
	dbErr := db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Phone, &user.CreatedAt)

	if dbErr == sql.ErrNoRows {
		return entity.User{}, errors.New("Invalid Email or Password")
	}
	if dbErr != nil {
		return entity.User{}, fmt.Errorf("Database error: %w", dbErr)
	}

	// TODO: Hashed password matching
	if password == user.Password {
		fmt.Println("Cuayo Sign In")
		return user, nil
	}

	return entity.User{}, errors.New("Invalid Email or Password")
}

func IsValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
