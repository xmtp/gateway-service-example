package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xmtp/xmtpd/pkg/gateway"
)

var (
	ErrMissingToken     = errors.New("missing JWT token")
	ErrInvalidToken     = errors.New("invalid JWT token")
	ErrInvalidSignature = errors.New("invalid token signature")
	ErrTokenExpired     = errors.New("token has expired")
)

func JWTES256(publicKeyPEM string, expectedIssuer string) gateway.IdentityFn {
	pubKey := mustParsePublicKey(publicKeyPEM)

	return func(ctx context.Context) (gateway.Identity, error) {
		raw := gateway.AuthorizationHeaderFromContext(ctx)
		if raw == "" {
			return gateway.Identity{}, gateway.NewUnauthenticatedError(
				"missing JWT token",
				ErrMissingToken,
			)
		}

		tokenStr := stripBearer(raw)

		claims := jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(
			tokenStr,
			&claims,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, ErrInvalidSignature
				}
				return pubKey, nil
			},
			jwt.WithIssuer(expectedIssuer),
			jwt.WithExpirationRequired(),
			jwt.WithValidMethods([]string{"ES256"}),
		)

		if err != nil || !token.Valid {
			return gateway.Identity{}, gateway.NewPermissionDeniedError(
				"failed to validate token",
				err,
			)
		}

		if time.Now().After(claims.ExpiresAt.Time) {
			return gateway.Identity{}, gateway.NewPermissionDeniedError(
				"token expired",
				ErrTokenExpired,
			)
		}

		sub, err := claims.GetSubject()
		if err != nil || sub == "" {
			return gateway.Identity{}, gateway.NewPermissionDeniedError(
				"missing subject",
				ErrInvalidToken,
			)
		}

		return gateway.NewUserIdentity(sub), nil
	}
}

func mustParsePublicKey(pemStr string) *ecdsa.PublicKey {
	pemStr = strings.ReplaceAll(pemStr, "\\n", "\n")

	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		panic("invalid JWT_PUBLIC_KEY PEM")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		panic("JWT public key is not ECDSA")
	}

	return ecdsaPub
}

func stripBearer(v string) string {
	v = strings.TrimSpace(v)
	if strings.HasPrefix(strings.ToLower(v), "bearer ") {
		return strings.TrimSpace(v[7:])
	}
	return v
}
