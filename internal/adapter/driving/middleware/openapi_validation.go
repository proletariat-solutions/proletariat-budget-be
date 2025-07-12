package middleware

import (
	"encoding/json"
	"ghorkov32/proletariat-budget-be/openapi"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/rs/zerolog/log"
)

// OpenAPIValidationMiddleware creates a middleware that validates requests against the OpenAPI spec
func OpenAPIValidationMiddleware(spec *openapi3.T) func(http.Handler) http.Handler {
	// Create validation options
	options := &nethttpmiddleware.Options{
		Options: openapi3filter.Options{
			// Include response body validation if needed
			IncludeResponseStatus: true,
			// Exclude certain operations from validation if needed
			ExcludeRequestBody:  false,
			ExcludeResponseBody: false,
			AuthenticationFunc:  openapi3filter.NoopAuthenticationFunc,
		},
		// Custom error handler
		ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
			log.Error().
				Int("status_code", statusCode).
				Str("validation_error", message).
				Msg("OpenAPI validation failed")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			errorObj := openapi.N400JSONResponse{
				Message: message,
			}
			err := json.NewEncoder(w).Encode(errorObj)
			if err != nil {
				log.Error().Err(err).Msg("Failed to encode error response")
			}
		},
	}

	// Create the middleware
	return nethttpmiddleware.OapiRequestValidatorWithOptions(spec, options)
}
