package providers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/carlosEA28/smartcondo/internal/dto"
)

type fakeCognitoClient struct {
	signUpResult     *cognitoidentityprovider.SignUpOutput
	signUpErr        error
	deleteUserErr    error
	signUpCalled     bool
	lastSignUpInput  *cognitoidentityprovider.SignUpInput
}

func (f *fakeCognitoClient) SignUp(_ context.Context, input *cognitoidentityprovider.SignUpInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.SignUpOutput, error) {
	f.signUpCalled = true
	f.lastSignUpInput = input
	return f.signUpResult, f.signUpErr
}

func (f *fakeCognitoClient) AdminDeleteUser(_ context.Context, _ *cognitoidentityprovider.AdminDeleteUserInput, _ ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDeleteUserOutput, error) {
	return nil, f.deleteUserErr
}

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

func TestCreateUserSuccess(t *testing.T) {
	client := &fakeCognitoClient{
		signUpResult: &cognitoidentityprovider.SignUpOutput{
			UserConfirmed: true,
		},
	}
	provider := &AwsProvider{
		CognitoClientID:     "test-client-id",
		CognitoClientSecret: "test-secret",
		cognitoClient:       client,
	}

	confirmed, err := provider.CreateUser(context.Background(), &dto.CreateUserDTO{
		Email:    "user@example.com",
		Password: "Password123!",
	})

	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}
	if !confirmed {
		t.Fatal("CreateUser() confirmed = false, want true")
	}
	if !client.signUpCalled {
		t.Fatal("CreateUser() did not call SignUp")
	}
	if client.lastSignUpInput == nil {
		t.Fatal("CreateUser() did not capture SignUpInput")
	}
	if *client.lastSignUpInput.Username != "user@example.com" {
		t.Fatalf("CreateUser() Username = %q, want %q", *client.lastSignUpInput.Username, "user@example.com")
	}
}

func TestCreateUserNotConfirmed(t *testing.T) {
	client := &fakeCognitoClient{
		signUpResult: &cognitoidentityprovider.SignUpOutput{
			UserConfirmed: false,
		},
	}
	provider := &AwsProvider{
		CognitoClientID:     "test-client-id",
		CognitoClientSecret: "test-secret",
		cognitoClient:       client,
	}

	confirmed, err := provider.CreateUser(context.Background(), &dto.CreateUserDTO{
		Email:    "user@example.com",
		Password: "Password123!",
	})

	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}
	if confirmed {
		t.Fatal("CreateUser() confirmed = true, want false")
	}
}

func TestCreateUserInvalidPasswordException(t *testing.T) {
	invalidPasswordErr := &types.InvalidPasswordException{
		Message: aws.String("Password did not conform with policy"),
	}
	client := &fakeCognitoClient{
		signUpErr: invalidPasswordErr,
	}
	provider := &AwsProvider{
		CognitoClientID:     "test-client-id",
		CognitoClientSecret: "test-secret",
		cognitoClient:       client,
	}

	confirmed, err := provider.CreateUser(context.Background(), &dto.CreateUserDTO{
		Email:    "user@example.com",
		Password: "weak",
	})

	if err == nil {
		t.Fatal("CreateUser() expected error, got nil")
	}
	if confirmed {
		t.Fatal("CreateUser() confirmed = true, want false on error")
	}
	var typedErr *types.InvalidPasswordException
	if !errors.As(err, &typedErr) {
		t.Fatalf("CreateUser() error type = %T, want *types.InvalidPasswordException", err)
	}
}

func TestCreateUserGenericError(t *testing.T) {
	genericErr := errors.New("service unavailable")
	client := &fakeCognitoClient{
		signUpErr: genericErr,
	}
	provider := &AwsProvider{
		CognitoClientID:     "test-client-id",
		CognitoClientSecret: "test-secret",
		cognitoClient:       client,
	}

	confirmed, err := provider.CreateUser(context.Background(), &dto.CreateUserDTO{
		Email:    "user@example.com",
		Password: "Password123!",
	})

	if err == nil {
		t.Fatal("CreateUser() expected error, got nil")
	}
	if confirmed {
		t.Fatal("CreateUser() confirmed = true, want false on error")
	}
	if !errors.Is(err, genericErr) {
		t.Fatalf("CreateUser() error = %v, want %v", err, genericErr)
	}
}

func TestDeleteUserSuccess(t *testing.T) {
	client := &fakeCognitoClient{}
	provider := &AwsProvider{
		CognitoUserPoolID: "us-east-1_poolid",
		cognitoClient:     client,
	}

	err := provider.DeleteUser(context.Background(), "user@example.com")
	if err != nil {
		t.Fatalf("DeleteUser() error = %v", err)
	}
}

func TestDeleteUserError(t *testing.T) {
	deleteErr := errors.New("user not found in pool")
	client := &fakeCognitoClient{
		deleteUserErr: deleteErr,
	}
	provider := &AwsProvider{
		CognitoUserPoolID: "us-east-1_poolid",
		cognitoClient:     client,
	}

	err := provider.DeleteUser(context.Background(), "user@example.com")
	if err == nil {
		t.Fatal("DeleteUser() expected error, got nil")
	}
	if !errors.Is(err, deleteErr) {
		t.Fatalf("DeleteUser() error = %v, want %v", err, deleteErr)
	}
}
