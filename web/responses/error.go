package responses

type ErrorResponse struct {
	Message interface{} `json:"message"`
}
