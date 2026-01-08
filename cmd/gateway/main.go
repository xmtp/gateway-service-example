package main

import (
	"log"

	auth "github.com/xmtp/gateway-service-example/internal/auth"
	"github.com/xmtp/gateway-service-example/internal/config"
	"github.com/xmtp/xmtpd/pkg/gateway"
)

func main() {
	cfg := config.Load()

	builder := gateway.NewGatewayServiceBuilder(
		gateway.MustLoadConfig(),
	)

	switch cfg.AuthMode {
	case config.AuthAllowAll:
		log.Println("✓ identity mode: allowall")
		builder = builder.WithIdentityFn(auth.AllowAll())

	case config.AuthJWTES256:
		log.Println("✓ identity mode: jwt_es256")
		builder = builder.WithIdentityFn(
			auth.JWTES256(
				cfg.JWT.PublicKeyPEM,
				cfg.JWT.ExpectedIssuer,
			),
		)

	default:
		log.Fatalf("unknown auth mode: %s", cfg.AuthMode)
	}

	log.Printf("auth mode: %s", cfg.AuthMode)
	if cfg.AuthMode == config.AuthJWTES256 {
		log.Printf("jwt issuer: %s", cfg.JWT.ExpectedIssuer)
	}

	// Rate limiting can be added here later without touching identity code

	service, err := builder.Build()
	if err != nil {
		log.Fatalf("failed to build gateway: %v", err)
	}

	log.Println("✓ gateway started")
	service.WaitForShutdown()
}
