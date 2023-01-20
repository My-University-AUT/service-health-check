package token

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var checkTokenIsValid = `SELECT COUNT(*) FROM token WHERE token.id=?`
var revokeToken = `DELETE FROM token WHERE token.id=?`

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issue_at"`
	ExpiredAt time.Time `json:"expired_at"`
	Role      string    `json:"role"`
}

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (*Payload, string, error)
	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}

const minSecretKeySize = 32

var ErrExpiredToken = errors.New("token has expired")
var ErrInvalidToken = errors.New("token is not validdd")

type JWTMaker struct {
	secretKey string
	DB        *sql.DB
}

func NewPayload(username, role string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (maker *JWTMaker) CreateToken(username, role string, duration time.Duration) (*Payload, string, error) {

	payload, err := NewPayload(username, role, duration)
	if err != nil {
		return nil, "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, error := jwtToken.SignedString([]byte(maker.secretKey))
	return payload, token, error
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	// fmt.Println("token is:", token)
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	// check token in database with sql function
	exists := false

	rows := maker.DB.QueryRow(checkTokenIsValid, payload.ID)
	err = rows.Scan(&exists)

	fmt.Println("shittttttt ", err, exists, " jj")
	if err != nil || !exists {
		return nil, errors.New("token not found!!")
	}
	return payload, nil
}

func (maker *JWTMaker) RevokeToken(token string) (map[string]interface{}, error) {
	// fmt.Println("token is:", token)
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	// check token in database with sql function
	res, err := maker.DB.Exec(revokeToken, payload.ID)
	if err == nil {
		count := int64(0)
		count, err = res.RowsAffected()
		if err == nil {
			fmt.Println("here is count", count)
			if count == 0 {
				return nil, ErrExpiredToken
			}
		}
	}

	fmt.Println("shittttttt ", err, " jj")
	if err != nil {
		return nil, errors.New("token not found!!")
	}

	return map[string]interface{}{"result": "token revoked successfully!"}, nil
}

func NewJWTMaker(secretKey string, db *sql.DB) (*JWTMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey, db}, nil
}
