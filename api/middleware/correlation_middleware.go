package middleware

import (
	"net/http"
	"our-expenses-server/logger"
	"our-expenses-server/utils"

	"github.com/google/uuid"
)

func CorrelationMiddleware(log logger.AppLoggerInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			correlationIDHeader := utils.CorrelationIDHeader
			correlationId := r.Header.Get(correlationIDHeader)
			if correlationId == "" {
				id := uuid.New()
				correlationId = id.String()
				r.Header.Set(correlationIDHeader, correlationId)
				log.Infof(r.Context(), "No %s HTTP header, using a new correlation id %s", correlationIDHeader, id.String())
			}
			ctx := utils.SetContextStringValue(r.Context(), utils.ContextKeyCorrelationID, correlationId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
