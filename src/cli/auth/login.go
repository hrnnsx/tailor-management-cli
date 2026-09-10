package auth

import (
	"errors"
	"fmt"
	"strings"

	"tailor-management-cli/entity"
	"tailor-management-cli/handler"

	"github.com/manifoldco/promptui"
)

func LoginCLI(authHandler *handler.AuthHandler) *entity.User {
	fmt.Println("\n==============================")
	fmt.Println("         FORM LOGIN           ")
	fmt.Println("==============================")

	inputTemplate := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . }} ",
		Invalid: "{{ . }} ",
		Success: " ",
	}

	// 1. Input Email
	fmt.Println("\nMasukkan Email:")
	emailPrompt := promptui.Prompt{
		Label:     ">",
		Templates: inputTemplate,
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return errors.New("email tidak boleh kosong")
			}
			return nil
		},
	}
	email, err := emailPrompt.Run()
	if err != nil {
		return nil
	}

	// 2. Input Password
	fmt.Println("\nMasukkan Password:")
	passPrompt := promptui.Prompt{
		Label:     ">",
		Mask:      '*',
		Templates: inputTemplate,
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return errors.New("password tidak boleh kosong")
			}
			return nil
		},
	}
	password, err := passPrompt.Run()
	if err != nil {
		return nil
	}

	// 3. Panggil Handler
	user, err := authHandler.Login(email, password)
	if err != nil {
		fmt.Printf("\n[Login Gagal]: %v\n\n", err)
		return nil
	}

	fmt.Printf("\nLogin berhasil! Selamat datang, %s.\n", user.Name)
	return user
}
