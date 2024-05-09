package services

import (
	"context"
	"fmt"
	"github.com/SeaCloudHub/notification-hub/pkg/config"
	kratos "github.com/ory/kratos-client-go"
)

type IdentityService struct {
	publicClient *kratos.APIClient
	adminClient  *kratos.APIClient
}

func NewIdentityService(cfg *config.Config) *IdentityService {
	return &IdentityService{
		publicClient: newKratosClient(cfg.Kratos.PublicURL, cfg.Debug),
		adminClient:  newKratosClient(cfg.Kratos.AdminURL, cfg.Debug),
	}
}

func newKratosClient(url string, debug bool) *kratos.APIClient {
	configuration := kratos.NewConfiguration()
	configuration.Servers = kratos.ServerConfigurations{{URL: url}}
	configuration.Debug = debug

	return kratos.NewAPIClient(configuration)
}

func (s *IdentityService) WhoAmI(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("token is empty")
	}
	session, _, err := s.publicClient.FrontendAPI.ToSession(ctx).XSessionToken(token).Execute()
	if err != nil {

		return "", fmt.Errorf("unexpected error: %w", err)
	}

	return session.Identity.Id, nil
}
