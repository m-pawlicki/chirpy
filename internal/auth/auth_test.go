package auth

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestTokenCreation(t *testing.T) {
	godotenv.Load(".env")
	testUser := uuid.New()
	testSecret := os.Getenv("TEST_SECRET")
	_, err := MakeJWT(testUser, testSecret)
	if err != nil {
		t.Errorf("error making jwt: %v", err)
	}
}

func TestTokenValidation(t *testing.T) {
	godotenv.Load(".env")
	testUser := uuid.New()
	testSecret := os.Getenv("TEST_SECRET")
	jwt, _ := MakeJWT(testUser, testSecret)
	_, err := ValidateJWT(jwt, testSecret)
	if err != nil {
		t.Errorf("error validating jwt: %v", err)
	}
}
