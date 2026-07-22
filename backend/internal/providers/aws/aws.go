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

type CognitoClient interface {
	SignUp(ctx context.Context, params *cognitoidentityprovider.SignUpInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.SignUpOutput, error)
	AdminDeleteUser(ctx context.Context, params *cognitoidentityprovider.AdminDeleteUserInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDeleteUserOutput, error)
}

type AwsProvider struct {
	client              *aws.Config
	cognitoClient       CognitoClient
	CognitoClientID     string
	CognitoClientSecret string
	CognitoUserPoolID   string
	S3Bucket            string
	s3Endpoint          string
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
		S3Bucket:            cfg.AWS.S3Bucket,
		s3Endpoint:          cfg.AWS.S3Endpoint,
	}
}

func (a *AwsProvider) GetS3Client() *s3.Client {
	return s3.NewFromConfig(*a.client, func(o *s3.Options) {
		if a.s3Endpoint != "" {
			o.BaseEndpoint = aws.String(a.s3Endpoint)
			o.UsePathStyle = true
		}
	})
}

func (a *AwsProvider) GetCognitoClient() CognitoClient {
	if a.cognitoClient != nil {
		return a.cognitoClient
	}
	return cognitoidentityprovider.NewFromConfig(*a.client)
}
