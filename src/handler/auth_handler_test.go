package handler_test

import (
	"errors"
	"regexp"
	"tailor-management-cli/handler"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRegister_Success(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	authH := handler.NewAuthHandler(db)

	expectedQuery := regexp.QuoteMeta(`
		INSERT INTO users (
			name, 
			email,
			password,
			role,
			phone
		)
		VALUES (?, ?, ?, 'customer', ?)
	`)

	mock.ExpectExec(expectedQuery).
		WithArgs("Ben", "ben@mail.com", "Makan123", "081234567890").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = authH.Register("Ben", "ben@mail.com", "Makan123", "081234567890")

	assert.NoError(t, err, "Registrasi aman")
	assert.NoError(t, mock.ExpectationsWereMet(), "Semua query database yang diharapkan harus terpenuhi")
}

func TestRegister_Failed_DuplicateEmail(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	authH := handler.NewAuthHandler(db)

	expectedQuery := regexp.QuoteMeta(`
		INSERT INTO users (
			name, 
			email,
			password,
			role,
			phone
		)
		VALUES (?, ?, ?, 'customer', ?)
	`)

	mock.ExpectExec(expectedQuery).
		WithArgs("Jajang", "jajang@mail.com", "Pass123", "081234567890").
		WillReturnError(errors.New("[Error duplicate]: Duplicate entry 'jajang@mail.com' for key 'users.email'"))

	err = authH.Register("Jajang", "jajang@mail.com", "Pass123", "081234567890")

	assert.Error(t, err, "return error karena duplikasi email!")
	assert.EqualError(t, err, "Email sudah terdaftar, silakan gunakan email lain!")
	assert.NoError(t, mock.ExpectationsWereMet(), "Semua query database yang diharapkan harus terpenuhi")
}
