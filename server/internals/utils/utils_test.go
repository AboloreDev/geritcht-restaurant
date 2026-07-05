package utils

import (
	"testing"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/config"
)

func Test_GenerateVerificationToken(t *testing.T) {
	result, err := GenerateVerificationToken()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == "" {
		t.Error("Expected token, got empty string. ", err)
	}

	if len(result) != 6 {
		t.Error("Expected token with length 6, got ", len(result))
	}
}

func Test_HashToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid token", "123456", true},
		{"empty token", "", false},
		{"token with special chars", "123@456", true},
		{"token with spaces", "123 456", false},
		{"token with numbers", "123456", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HashToken(tt.input)
			if tt.expected && result == "" {
				t.Errorf("Expected hashed token, got empty string")
			}
			if tt.expected && result == "123 456" {
				t.Errorf("Expected hashed token, got token with spaces")
			}
		})
	}
}

func Test_GenerateOTP(t *testing.T) {
	result, err := GenerateOTP()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == "" {
		t.Error("Expected otp, got empty string. ", err)
	}

	if len(result) != 6 {
		t.Error("Expected otp with length 6, got ", len(result))
	}
}

func Test_IsValidExtension(t *testing.T) {
	Tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid jpg", ".jpg", true},
		{"valid jpeg", ".jpeg", true},
		{"valid png", ".png", true},
		{"valid gif", ".gif", true},
		{"valid webp", ".webp", true},

		{"invalid exe", ".exe", false},
		{"invalid txt", ".txt", false},
		{"empty string", "", false},
	}

	for _, tt := range Tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidExtensions(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}

}

func Test_GeneratePaymentReference(t *testing.T) {
	result := GeneratePaymentReference()
	if result == "" {
		t.Error("Expected reference, got empty string")
	}
	if result == "payment_ref" {
		t.Errorf("Expected PAY-%X, got default value", result)
	}

}

func Test_GenerateIdempotencyKey(t *testing.T) {
	result := GenerateIdempotencyKey()
	if result == "" {
		t.Error("Expected idempotency key, got empty string")
	}

}

func Test_ValidateToken(t *testing.T) {
	cfg := &config.JWTConfig{
		JWTSecret:                 "1234567890",
		JWTTokenExpiration:        time.Minute,
		JWTRefreshTokenExpiration: time.Hour,
	}
	// Generate a token for testing
	token, _, err := GenerateTokenPair(
		cfg, 1, "test@gmail.com", "test")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tests := []struct {
		name      string
		token     string
		secret    string
		expectErr bool
	}{
		{"valid token", token, cfg.JWTSecret, false},
		{"invalid token", "12345", cfg.JWTSecret, true},
		{"invalid token long", "invalid.token.", cfg.JWTSecret, true},
		{"empty token", "", cfg.JWTSecret, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			claims, err := ValidateToken(tt.token, tt.secret)

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if claims == nil {
				t.Fatal("expected claims, got nil")
			}

			if claims.UserID != 1 {
				t.Errorf("expected userID 1, got %d", claims.UserID)
			}
		})
	}
}

func Test_HashPassword(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid password", "password", true},
		{"empty password", "", false},
		{"password with special chars", "pass@word123", true},
		{"password with spaces", "pass word", true},
		{"password with numbers", "password123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HashPassword(tt.input)
			if err != nil && tt.expected {
				t.Fatalf("Expected no error, got %v", err)
			}

			if tt.expected && result == "" {
				t.Errorf("Expected hashed password, got empty string")
			}
		})
	}
}

func Test_CheckPassword(t *testing.T) {
	hashPassword, _ := HashPassword("password")
	password := "password"

	result := CheckPassword(hashPassword, password)
	if !result {
		t.Error("Expected true, got false")
	}
	if result {
		t.Log("Password matched")
	} else {
		t.Error("Password did not match")
	}
}

func Test_Pagination(t *testing.T) {
	tests := []struct {
		name           string
		page           int
		pageSize       int
		expectedOffset int
	}{
		{"valid page and limit", 1, 10, 0},
		{"valid page and limit 2", 2, 10, 10},
		{"valid page and limit 3", 3, 10, 20},
		{"page 0", 0, 10, 0},
		{"limit 0", 1, 0, 0},
		{"page and limit 0", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset := Pagination(tt.page, tt.pageSize)
			if offset != tt.expectedOffset {
				t.Errorf("Expected offset %d, got %d", tt.expectedOffset, offset)
			}
			if offset < 0 {
				t.Errorf("Expected offset to be positive, got %d", offset)
			}
			if offset > 100 {
				t.Errorf("Expected offset to be less than 100, got %d", offset)
			}

		})
	}
}

func Test_IsValidTimeSlot(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid time slot", "22:00:00", true},
		{"invalid time slot", "25:00:00", false},
		{"invalid time slot 2", "09:70:00", false},
		{"invalid time slot 3", "ab:cd", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidTimeSlots(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}

			time, err := ParseToDataTypesTime(tt.input)
			if err != nil && result {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expected && time.String() == "0001-01-01 00:00:00 +0000 UTC" {
				t.Errorf("Expected time, got zero value")
			}

			timeStr := FormatDataTypesTime(time)
			if timeStr == "0001-01-01 00:00:00 +0000 UTC" && result {
				t.Errorf("Expected formatted time, got zero value")
			}
		})
	}
}
