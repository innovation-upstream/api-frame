package middleware

import (
  "net/http"
)

// CORS returns a middleware function that sets access control headers for all request methods.
// CORS-safelisted headers are allowed by default and need not be specified in allowedHeaders. See: https://developer.mozilla.org/en-US/docs/Glossary/CORS-safelisted_request_header.
// allowedHeaders may not be a wildcard.
// allowCredentials must be "true" or an empty string.
func CORS(allowedOrigin string, allowedHeaders string, allowCredentials string) func(http.Handler) http.Handler {
  if allowedHeaders == "*" {
    panic("allowedHeaders may not be a wildcard due to incompatibility with Safari and Internet Explorer. See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers.")
  }

  if allowCredentials != "true" && allowCredentials != "" {
    panic("Invalid value for allowCredentials. Must be either \"true\" or an empty string.")
  }

  return func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
      w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
      if allowCredentials == "true" {
        w.Header().Set("Access-Control-Allow-Credentials", allowCredentials)
      }

      if r.Method == "OPTIONS" {
        return
      }

      next.ServeHTTP(w, r)
    })
  }
}

