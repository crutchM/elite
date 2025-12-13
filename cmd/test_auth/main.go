package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// AuthRequest содержит данные от Telegram Login Widget
type AuthRequest struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	PhotoURL  string `json:"photo_url,omitempty"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

type AuthResponse struct {
	Token string `json:"token"`
	Error string `json:"error,omitempty"`
}

type VoteRequest struct {
	NominantID int64 `json:"nominant_id"`
}

type VoteResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func main() {
	if len(os.Args) < 8 {
		fmt.Println("Использование:")
		fmt.Println("  go run cmd/test_auth/main.go <id> <first_name> <last_name> <username> <photo_url> <auth_date> <hash> [server_url]")
		fmt.Println("")
		fmt.Println("Пример:")
		fmt.Println("  go run cmd/test_auth/main.go 123456789 'John' 'Doe' 'johndoe' '' 1234567890 'abc123...' http://localhost:8080")
		fmt.Println("")
		fmt.Println("Для получения данных от Login Widget:")
		fmt.Println("  1. Создайте Telegram бота через @BotFather")
		fmt.Println("  2. Добавьте Login Widget на страницу (см. example_login_widget.html)")
		fmt.Println("  3. После авторизации получите данные из callback функции onTelegramAuth")
		fmt.Println("")
		fmt.Println("Или используйте HTML страницу example_login_widget.html для тестирования")
		os.Exit(1)
	}

	id, _ := strconv.ParseInt(os.Args[1], 10, 64)
	firstName := os.Args[2]
	lastName := os.Args[3]
	username := os.Args[4]
	photoURL := os.Args[5]
	authDate, _ := strconv.ParseInt(os.Args[6], 10, 64)
	hash := os.Args[7]

	serverURL := "http://localhost:8080"
	if len(os.Args) > 8 {
		serverURL = os.Args[8]
	}

	fmt.Printf("🔐 Тестирование Telegram Login Widget авторизации\n")
	fmt.Printf("Сервер: %s\n\n", serverURL)

	// Шаг 1: Авторизация
	fmt.Println("1️⃣ Отправка запроса на авторизацию...")
	token, err := authenticate(serverURL, AuthRequest{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		PhotoURL:  photoURL,
		AuthDate:  authDate,
		Hash:      hash,
	})
	if err != nil {
		fmt.Printf("❌ Ошибка авторизации: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Авторизация успешна!\n")
	fmt.Printf("Токен: %s\n\n", token)

	// Шаг 2: Тест голосования
	fmt.Println("2️⃣ Тестирование голосования...")
	fmt.Print("Введите ID номинанта (или Enter для пропуска): ")
	
	var nominantID int64
	fmt.Scanf("%d\n", &nominantID)
	
	if nominantID > 0 {
		err = testVote(serverURL, token, nominantID)
		if err != nil {
			fmt.Printf("❌ Ошибка голосования: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Голос успешно засчитан!\n")
	} else {
		fmt.Println("⏭️  Голосование пропущено")
	}

	fmt.Println("\n✅ Тестирование завершено!")
}

func authenticate(serverURL string, reqBody AuthRequest) (string, error) {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(serverURL+"/api/auth", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp AuthResponse
		json.Unmarshal(body, &errResp)
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, errResp.Error)
	}

	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return authResp.Token, nil
}

func testVote(serverURL, token string, nominantID int64) error {
	reqBody := VoteRequest{NominantID: nominantID}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", serverURL+"/api/vote", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		var errResp VoteResponse
		json.Unmarshal(body, &errResp)
		return fmt.Errorf("status %d: %s", resp.StatusCode, errResp.Error)
	}

	return nil
}

