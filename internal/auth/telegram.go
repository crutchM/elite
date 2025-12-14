package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TelegramAuth struct {
	botToken  string
	jwtSecret string
}

func NewTelegramAuth(botToken, jwtSecret string) *TelegramAuth {
	return &TelegramAuth{
		botToken:  botToken,
		jwtSecret: jwtSecret,
	}
}

type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	PhotoURL  string `json:"photo_url,omitempty"`
}

// LoginWidgetData представляет данные от Telegram Login Widget
type LoginWidgetData struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	PhotoURL  string `json:"photo_url,omitempty"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

func (ta *TelegramAuth) ValidateLoginWidgetData(data *LoginWidgetData) (*TelegramUser, error) {
	if data.ID == 0 {
		return nil, errors.New("user ID is required")
	}

	// Проверяем свежесть данных
	if time.Now().Unix()-data.AuthDate > 86400 {
		return nil, errors.New("auth data is too old")
	}

	// Собираем все поля кроме hash
	fields := make(map[string]string)
	fields["id"] = strconv.FormatInt(data.ID, 10)
	fields["first_name"] = data.FirstName
	if data.LastName != "" {
		fields["last_name"] = data.LastName
	}
	if data.Username != "" {
		fields["username"] = data.Username
	}
	if data.PhotoURL != "" {
		fields["photo_url"] = data.PhotoURL
	}
	fields["auth_date"] = strconv.FormatInt(data.AuthDate, 10)

	// Сортируем ключи и формируем data_check_string
	var keys []string
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var dataCheckParts []string
	for _, k := range keys {
		dataCheckParts = append(dataCheckParts, fmt.Sprintf("%s=%s", k, fields[k]))
	}
	dataCheckStr := strings.Join(dataCheckParts, "\n")

	// Telegram требует: secret_key = SHA256(bot_token)
	secretKey := sha256.Sum256([]byte(ta.botToken))

	// HMAC-SHA256
	mac := hmac.New(sha256.New, secretKey[:])
	mac.Write([]byte(dataCheckStr))
	calculatedHash := hex.EncodeToString(mac.Sum(nil))

	if calculatedHash != data.Hash {
		return nil, fmt.Errorf("invalid hash: got %s, expected %s", data.Hash, calculatedHash)
	}

	return &TelegramUser{
		ID:        data.ID,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Username:  data.Username,
		PhotoURL:  data.PhotoURL,
	}, nil
}

func (ta *TelegramAuth) GenerateToken(tgUserID int64) (string, error) {
	claims := jwt.MapClaims{
		"tg_user_id": tgUserID,
		"exp":        time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(ta.jwtSecret))
}

func (ta *TelegramAuth) ValidateToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ta.jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tgUserID, ok := claims["tg_user_id"].(float64)
		if !ok {
			return 0, errors.New("invalid token claims")
		}
		return int64(tgUserID), nil
	}

	return 0, errors.New("invalid token")
}
