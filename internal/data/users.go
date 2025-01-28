package data

import (
	"runtime"
	"time"

	"github.com/Robert-litts/fantasy-football-archive/internal/validator"
	"github.com/alexedwards/argon2id"
)

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type password struct {
	plaintext *string
	hash      *string
}

var params = &argon2id.Params{
	Memory:      128 * 1024, //128MB
	Iterations:  4,
	Parallelism: uint8(runtime.NumCPU()),
	SaltLength:  16,
	KeyLength:   32,
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := argon2id.CreateHash(plaintextPassword, params)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = &hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(plaintextPassword, *p.hash)
	if err != nil || !match {
		return false, err
	}

	return match, nil
}

func (p *password) Hash() *string {
	return p.hash
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	// Call the standalone ValidateEmail() helper.
	ValidateEmail(v, user.Email)

	// If the plaintext password is not nil, call the standalone
	// ValidatePasswordPlaintext() helper.
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	// If the password hash is ever nil, this will be due to a logic error in our
	// codebase,raise a panic.
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}
