package handler_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"tailor-management-cli/entity"
	"tailor-management-cli/handler"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// signIn() test
func TestAuthHandler_SignIn_Success(t *testing.T) {
	// Buat database mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handler.NewAuthHandler(db)

	// Query yang diharapkan
	query := `
		SELECT
			id,
			name,
			email,
			password,
			role,
			phone,
			created_at
		FROM users WHERE email = ? LIMIT 1
	`

	// Data hasil query
	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"email",
		"password",
		"role",
		"phone",
		"created_at",
	}).AddRow(
		1,
		"John Doe",
		"john@example.com",
		"123456",
		"customer",
		"08123456789",
		time.Now(),
	)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("john@example.com").
		WillReturnRows(rows)

	user, err := handler.SignIn(
		"john@example.com",
		"123456",
	)

	assert.NoError(t, err)
	assert.NotNil(t, user)

	assert.Equal(t, entity.User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "123456",
		Role:     "customer",
		Phone:    "08123456789",
	}, entity.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
		Phone:    user.Phone,
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthHandler_SignIn_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handler.NewAuthHandler(db)

	query := `
		SELECT
			id,
			name,
			email,
			password,
			role,
			phone,
			created_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("unknown@example.com").
		WillReturnError(sql.ErrNoRows)

	user, err := handler.SignIn(
		"unknown@example.com",
		"123456",
	)

	assert.Nil(t, user)
	assert.EqualError(t, err, "email atau password salah")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthHandler_SignIn_WrongPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handler.NewAuthHandler(db)

	query := `
		SELECT
			id,
			name,
			email,
			password,
			role,
			phone,
			created_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`

	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"email",
		"password",
		"role",
		"phone",
		"created_at",
	}).AddRow(
		1,
		"John Doe",
		"john@example.com",
		"correct-password",
		"customer",
		"08123456789",
		time.Now(),
	)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("john@example.com").
		WillReturnRows(rows)

	user, err := handler.SignIn(
		"john@example.com",
		"wrong-password",
	)

	assert.Nil(t, user)
	assert.EqualError(t, err, "email atau password salah")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthHandler_SignIn_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := handler.NewAuthHandler(db)

	query := `
		SELECT
			id,
			name,
			email,
			password,
			role,
			phone,
			created_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("john@example.com").
		WillReturnError(errors.New("database connection failed"))

	user, err := handler.SignIn(
		"john@example.com",
		"123456",
	)

	assert.Nil(t, user)
	assert.EqualError(
		t,
		err,
		"gagal mengambil data user: database connection failed",
	)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// signUp() test
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
