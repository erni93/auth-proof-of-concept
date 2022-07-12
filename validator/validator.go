package validator

import "net/http"

type Validator struct {
	Writer   http.ResponseWriter
	Request  *http.Request
	Services *Services
}
