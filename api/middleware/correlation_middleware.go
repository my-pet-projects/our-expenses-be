package middleware

import (
	"net/http"

	"github.com/google/uuid"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/utils"
)

// CorrelationMiddleware puts correlation id into the context.
func CorrelationMiddleware(log logger.LogInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			correlationIDHeader := utils.CorrelationIDHeader
			correlationID := r.Header.Get(correlationIDHeader)
			if correlationID == "" {
				id := uuid.New()
				correlationID = id.String()
				r.Header.Set(correlationIDHeader, correlationID)
				log.Infof(r.Context(), "No %s HTTP header, using a new correlation id %s", correlationIDHeader, id.String())
			}
			ctx := utils.SetContextStringValue(r.Context(), utils.ContextKeyCorrelationID, correlationID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
