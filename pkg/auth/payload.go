package auth

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Payload struct {
	UserId    string    `json:"userId"`
	AccountId uuid.UUID `json:"accountId"`
	Email     string    `json:"email"`
	Role      int64     `json:"role"`
}

type TokenData struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

func (payload *Payload) GetDataFromClaims(claims string) error {
	err := json.Unmarshal([]byte(claims), &payload)
	if err != nil {
		return err
	}
	return nil
}

func (payload *Payload) GetTokenDataJWT() (*TokenData, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	token, err := GetJWTToken(string(data))
	if err != nil {
		return nil, err
	}

	tokenData := &TokenData{
		Token: token,
		Type:  "Bearer",
	}

	return tokenData, nil
}
