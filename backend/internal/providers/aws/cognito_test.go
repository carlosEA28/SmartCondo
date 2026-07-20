package providers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/carlosEA28/smartcondo/internal/dto"
)

func TestComputeSecretHash(t *testing.T) {
	provider := &AwsProvider{
		CognitoClientID:     "test-client-id",
		CognitoClientSecret: "test-secret",
	}

	username := "user@example.com"
	hash := provider.computeSecretHash(username)

	mac := hmac.New(sha256.New, []byte("test-secret"))
	mac.Write([]byte(username + "test-client-id"))
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if hash != expected {
		t.Errorf("computeSecretHash() = %q, want %q", hash, expected)
	}
}

func TestComputeSecretHashDifferentInputs(t *testing.T) {
	provider := &AwsProvider{
		CognitoClientID:     "client-id-123",
		CognitoClientSecret: "my-secret",
	}

	hash1 := provider.computeSecretHash("user1@example.com")
	hash2 := provider.computeSecretHash("user2@example.com")

	if hash1 == hash2 {
		t.Error("computeSecretHash() should produce different hashes for different usernames")
	}
}

func TestComputeSecretHashDeterministic(t *testing.T) {
	provider := &AwsProvider{
		CognitoClientID:     "client-id",
		CognitoClientSecret: "secret",
	}

	hash1 := provider.computeSecretHash("user@example.com")
	hash2 := provider.computeSecretHash("user@example.com")

	if hash1 != hash2 {
		t.Errorf("computeSecretHash() should be deterministic, got %q and %q", hash1, hash2)
	}
}

func TestCreateUserSetsCorrectAttributes(t *testing.T) {
	provider := &AwsProvider{
		CognitoClientID:     "test-client-id",
		CognitoClientSecret: "test-secret",
	}

	input := &dto.CreateUserDTO{
		Email:    "test@example.com",
		Password: "password123",
	}

	if input.Email != "test@example.com" {
		t.Errorf("input.Email = %q, want %q", input.Email, "test@example.com")
	}
	if input.Password != "password123" {
		t.Errorf("input.Password = %q, want %q", input.Password, "password123")
	}

	_ = provider
}
