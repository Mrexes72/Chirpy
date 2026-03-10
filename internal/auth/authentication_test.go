package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword returned an empty string")
	}

	if hash == password {
		t.Fatal("HashPassword did not hash the password")
	}
}

func TestCheckPasswordHash_Correct(t *testing.T) {
	password := "correctpassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash failed: %v", err)
	}

	if !match {
		t.Error("Expected password to match hash")
	}

}
func TestCheckPasswordHash_Incorrect(t *testing.T) {
	password := "correctpassword"
	wrongPassword := "wrongpassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	match, err := CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash failed: %v", err)
	}

	if match {
		t.Error("Expected wrong password to not match hash")
	}
}

func TestHashPassword_Uniqueness(t *testing.T) {
	password := "samepassword"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("First HashPassword failed: %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Second HashPassword failed: %v", err)
	}

	if hash1 == hash2 {
		t.Error("Two hashes of same password should be different (salted)")
	}

	match1, _ := CheckPasswordHash(password, hash1)
	match2, _ := CheckPasswordHash(password, hash2)

	if !match1 || !match2 {
		t.Error("Both hashes should verify the same password")
	}
}

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "my-secret-key"
	expiresIn := time.Hour

	tokenString, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	if tokenString == "" {
		t.Error("Token should not be empty")
	}
}

func TestValidateJWT_Success(t *testing.T) {
	userID := uuid.New()
	secret := "my-secret-key"
	expiresIn := time.Hour

	tokenString, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	parsedUserID, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if parsedUserID != userID {
		t.Errorf("Expected userID %v, got %v", userID, parsedUserID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "my-secret-key"
	expiresIn := -time.Hour // Token som allerede er utløpt

	tokenString, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(tokenString, secret)
	if err == nil {
		t.Error("Expected error for expired token")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	secret := "correct-secret"
	wrongSecret := "wrong-secret"
	expiresIn := time.Hour

	tokenString, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(tokenString, wrongSecret)
	if err == nil {
		t.Error("Expected error for wrong secret")
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	secret := "my-secret-key"

	_, err := ValidateJWT("not.a.valid.token", secret)
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestValidateJWT_MalformedToken(t *testing.T) {
	secret := "my-secret-key"

	_, err := ValidateJWT("garbage", secret)
	if err == nil {
		t.Error("Expected error for malformed token")
	}
}

func TestGetBearerToken_Success(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer my-token-string")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token != "my-token-string" {
		t.Errorf("Expected 'my-token-string', got %q", token)
	}
}

func TestGetBearerToken_NoHeader(t *testing.T) {
	headers := http.Header{}

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Error("Expected error for missing header")
	}
}

func TestGetBearerToken_WrongPrefix(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Basic my-token")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Error("Expected error for wrong prefix")
	}
}

func TestGetBearerToken_NoToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Error("Expected error for missing token")
	}
}
