package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/layer8s/home-dashboard-app/internal/db"
	"github.com/layer8s/home-dashboard-app/internal/validator"
)

const (
	ScopeActivation = "activation"
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

type TokenModel struct {
	queries *db.Queries
}

func NewTokenService(queries *db.Queries) *TokenModel {
	return &TokenModel{
		queries: queries,
	}
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	// Create a new Token instance with the provided values.
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	// Initialize a zero-valued byte slice with a length of 16 bytes.
	randomBytes := make([]byte, 16)

	// Use the Read() function from the crypto/rand package
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// Encode the byte slice to a base-32-encoded string and assign it to the token
	// Plaintext field. This will be the token string that we send to the user in their
	// welcome email. They will look similar to this:
	//
	// Y3QMGX3PJ3WLRL2YRTQGQ6KRHU
	//
	// Note that by default base-32 strings may be padded at the end with the =
	// character. We don't need this padding character for the purpose of our tokens, so
	// we use the WithPadding(base32.NoPadding) method in the line below to omit them.
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// Generate a SHA-256 hash of the plaintext token string.
	hash := sha256.Sum256([]byte(token.Plaintext))
	// convert it to a slice using the [:] operator before storing it.
	token.Hash = hash[:]

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

// The New() method is a shortcut which creates a new Token struct and then inserts the
// data in the tokens table.
func (m *TokenModel) New(ctx context.Context, userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(ctx, token)
	if err != nil {
		return nil, err
	}
	return token, err
}

// Insert() adds the data for a specific token to the tokens table.
func (m TokenModel) Insert(ctx context.Context, token *Token) error {
	params := db.InsertTokenParams{
		Hash:   token.Hash,
		UserID: token.UserID,
		Expiry: token.Expiry,
		Scope:  token.Scope,
	}

	err := m.queries.InsertToken(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAllForUser() deletes all tokens for a specific user and scope.
func (m TokenModel) DeleteAllForUser(ctx context.Context, scope string, userID int64) error {

	params := db.DeleteTokenParams{
		Scope:  scope,
		UserID: userID,
	}

	return m.queries.DeleteToken(ctx, params)
}
