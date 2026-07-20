package providers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/carlosEA28/smartcondo/internal/dto"
)

func (a *AwsProvider) computeSecretHash(username string) string {
	mac := hmac.New(sha256.New, []byte(a.CognitoClientSecret))
	mac.Write([]byte(username + a.CognitoClientID))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (a *AwsProvider) CreateUser(ctx context.Context, user *dto.CreateUserDTO) (bool, error) {
	confirmed := false

	output, err := a.GetCognitoClient().SignUp(ctx, &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(a.CognitoClientID),
		SecretHash: aws.String(a.computeSecretHash(user.Email)),
		Username:   aws.String(user.Email),
		Password:   aws.String(user.Password),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(user.Email),
			},
		},
	})

	if err != nil {
		var invalidPassword *types.InvalidPasswordException
		if errors.As(err, &invalidPassword) {
			log.Println(*invalidPassword.Message)
		} else {
			log.Printf("Couldn't sign up user %v. Here's why: %v\n", user.Email, err)
		}
	} else {
		confirmed = output.UserConfirmed
	}

	return confirmed, err
}
