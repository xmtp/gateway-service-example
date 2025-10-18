package main

import (
	"context"
	"fmt"
	"log"

	"github.com/xmtp/xmtpd/pkg/gateway"
)

func main() {
	gatewayService, err := gateway.NewGatewayServiceBuilder(gateway.MustLoadConfig()).
		WithIdentityFn(func(ctx context.Context) (gateway.Identity, error) {
			fmt.Println("Received request")
			return gateway.IPIdentityFn(ctx)
		}).
		WithAuthorizers(func(ctx context.Context, identity gateway.Identity, req gateway.PublishRequestSummary) (bool, error) {
			return true, nil
		}).
		Build() // This will gather all the config from environment variables and flags
	if err != nil {
		log.Fatalf("Failed to build gateway service: %v", err)
	}

	gatewayService.WaitForShutdown()
}
