package cors

import (
	"github.com/gorilla/handlers"
	"net/http"
)

func HandleCores(allowedOrigin []string) func(handler http.Handler) http.Handler {
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete})
	headers := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin", "Content-Type", "Authorization"})
	exposeHeaders := handlers.ExposedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin", "Content-Type", "Authorization"})
	allowedOriginHeaders := handlers.AllowedOrigins(allowedOrigin)
	return handlers.CORS(credentials, methods, headers, exposeHeaders, allowedOriginHeaders)
}
