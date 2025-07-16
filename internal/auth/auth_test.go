package auth

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestTokenCreation(t *testing.T) {
	godotenv.Load(".env")
	testUser := uuid.New()
	testSecret := os.Getenv("TEST_SECRET")
	_, err := MakeJWT(testUser, testSecret, time.Second*60)
	if err != nil {
		t.Errorf("error making jwt: %v", err)
	}
}

func TestTokenValidation(t *testing.T) {
	godotenv.Load(".env")
	testUser := uuid.New()
	testSecret := os.Getenv("TEST_SECRET")
	jwt, _ := MakeJWT(testUser, testSecret, time.Second*60)
	_, err := ValidateJWT(jwt, testSecret)
	if err != nil {
		t.Errorf("error validating jwt: %v", err)
	}
}

func TestTokenExpired(t *testing.T) {
	godotenv.Load(".env")
	testUser := uuid.New()
	testSecret := os.Getenv("TEST_SECRET")
	jwt, _ := MakeJWT(testUser, testSecret, time.Second*2)
	time.Sleep(time.Second * 5)
	_, err := ValidateJWT(jwt, testSecret)
	if err == nil {
		t.Errorf("expected expired token error but got nil instead")
	}
}

func TestBearerToken(t *testing.T) {
}

func TestEmptyAuthHead(t *testing.T) {
}
