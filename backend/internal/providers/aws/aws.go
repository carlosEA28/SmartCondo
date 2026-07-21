package providers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	appconfig "github.com/carlosEA28/smartcondo/internal/config"
)

type AwsProvider struct {
	client              *aws.Config
	CognitoClientID     string
	CognitoClientSecret string
	CognitoUserPoolID   string
}

func NewAwsProvider(cfg *appconfig.Config) *AwsProvider {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(cfg.AWS.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWS.AccessKeyID,
			cfg.AWS.SecretAccessKey,
			"",
		)))

	if err != nil {
		panic("failed to create AWS config " + err.Error())
	}

	return &AwsProvider{
		client:              &awsCfg,
		CognitoClientID:     cfg.AWS.CognitoClientId,
		CognitoClientSecret: cfg.AWS.CognitoClientSecret,
		CognitoUserPoolID:   cfg.AWS.CognitoUserPoolID,
	}
}

func (a *AwsProvider) GetS3Client() *s3.Client {
	return s3.NewFromConfig(*a.client)
}

func (a *AwsProvider) GetCognitoClient() *cognitoidentityprovider.Client {
	return cognitoidentityprovider.NewFromConfig(*a.client)
}
