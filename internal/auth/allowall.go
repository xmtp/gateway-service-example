package auth

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/gateway"
)

func AllowAll() gateway.IdentityFn {
	return func(ctx context.Context) (gateway.Identity, error) {
		// Anonymous but allowed
		return gateway.NewUserIdentity("anonymous"), nil
	}
}

