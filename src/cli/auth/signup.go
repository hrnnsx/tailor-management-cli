package auth

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"tailor-management-cli/handler"

	"github.com/manifoldco/promptui"
)

func SignUpCLI(authHandler *handler.AuthHandler) {
	fmt.Println("\n==============================")
	fmt.Println("    FORM REGISTRASI AKUN      ")
	fmt.Println("==============================")

	// Template tanpa mengulang kembali inputan teks user ke terminal
	inputTemplate := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . }} ",
		Invalid: "{{ . }} ",
		Success: " ",
	}

	// 1. Nama Lengkap
	fmt.Println("\nMasukkan Nama Lengkap:")
	namePrompt := promptui.Prompt{
		Label:     ">",
		Templates: inputTemplate,
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return errors.New("nama tidak boleh kosong")
			}
			return nil
		},
	}
	name, err := namePrompt.Run()
	if err != nil {
		return
	}

	// 2. Email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	fmt.Println("\nMasukkan Email (contoh: user@mail.com):")
	emailPrompt := promptui.Prompt{
		Label:     ">",
		Templates: inputTemplate,
		Validate: func(input string) error {
			if !emailRegex.MatchString(strings.TrimSpace(input)) {
				return errors.New("format email tidak valid (contoh: x@y.z)")
			}
			return nil
		},
	}
	email, err := emailPrompt.Run()
	if err != nil {
		return
	}

	// 3. Password
	fmt.Println("\nMasukkan Password (minimal 6 karakter):")
	passPrompt := promptui.Prompt{
		Label:     ">",
		Mask:      '*',
		Templates: inputTemplate,
		Validate: func(input string) error {
			if len(input) < 6 {
				return errors.New("password minimal 6 karakter")
			}
			return nil
		},
	}
	password, err := passPrompt.Run()
	if err != nil {
		return
	}

	// 4. Konfirmasi Password
	fmt.Println("\nKonfirmasi Ulang Password:")
	confirmPassPrompt := promptui.Prompt{
		Label:     ">",
		Mask:      '*',
		Templates: inputTemplate,
		Validate: func(input string) error {
			if input != password {
				return errors.New("konfirmasi password tidak cocok")
			}
			return nil
		},
	}
	_, err = confirmPassPrompt.Run()
	if err != nil {
		return
	}

	// 5. Nomor Handphone
	fmt.Println("\nMasukkan Nomor HP (contoh: 08xxxxxxxxxx):")
	phonePrompt := promptui.Prompt{
		Label:     ">",
		Templates: inputTemplate,
		Validate: func(input string) error {
			input = strings.TrimSpace(input)
			if !strings.HasPrefix(input, "08") || len(input) < 10 {
				return errors.New("nomor HP wajib diawali '08' dan minimal 10 digit")
			}
			return nil
		},
	}
	phone, err := phonePrompt.Run()
	if err != nil {
		return
	}

	// Eksekusi logic query ke Handler
	err = authHandler.Register(name, email, password, phone)
	if err != nil {
		fmt.Printf("\n[Error]: %v\n", err)
		return
	}

	fmt.Println("\nRegistrasi akun berhasil! Silakan login untuk masuk ke sistem.")
}
