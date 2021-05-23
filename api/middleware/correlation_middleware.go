package middleware

import (
	"net/http"

	"github.com/google/uuid"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/utils"
)

func CorrelationMiddleware(log logger.LogInterface) func(http.Handler) http.Handler {
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
