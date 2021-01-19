package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/coreos/go-oidc"
)

func VerifyCloudflareAccessJWT() func(http.Handler) http.Handler {
	var (
		ctx        = context.TODO()
		authDomain = os.Getenv("CLOUDFLARE_AUTH_DOMAIN")
		certsURL   = fmt.Sprintf("%s/cdn-cgi/access/certs", authDomain)
		policyAUD  = os.Getenv("CLOUDFLARE_AUDIENCE")
		config     = &oidc.Config{
			ClientID: policyAUD,
		}
		keySet   = oidc.NewRemoteKeySet(ctx, certsURL)
		verifier = oidc.NewVerifier(authDomain, keySet, config)
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookies := r.Cookies()
			var accessJWT string
			for _, cookie := range cookies {
				if cookie.Name == "CF_Authorization" {
					accessJWT = cookie.Value
				}
			}

			if accessJWT == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthenticated"))
				return
			}

			ctx := r.Context()
			_, err := verifier.Verify(ctx, accessJWT)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(fmt.Sprintf("Invalid token: %s", err.Error())))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
