package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIResponse defines attributes for api Response
type APIResponse struct {
	HTTPCode int         `json:"-"`
	Code     string      `json:"code"`
	Message  interface{} `json:"message"`
	Data     interface{} `json:"data,omitempty"`
}

var (
	APIOK = APIResponse{
		HTTPCode: http.StatusOK,
		Code:     "SUCCESS",
		Message:  "Success",
	}

	APINotFoundHandler = APIResponse{
		HTTPCode: http.StatusNotFound,
		Code:     "ENDPOINT_NOT_FOUND",
		Message:  "Endpoint not found",
	}

	APIErrNotFound = APIResponse{
		HTTPCode: http.StatusNotFound,
		Code:     "DATA_NOT_FOUND",
	}

	APIInternalError = APIResponse{
		HTTPCode: http.StatusInternalServerError,
		Code:     "SERVER_ERROR",
	}

	APIErrorBadRequest = APIResponse{
		HTTPCode: http.StatusBadRequest,
		Code:     "BAD_REQUEST",
	}
)

// WriteAPIOK for write response as HTTP OK result
func WriteAPIOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(APIOK)
}

// WriteAPIOKWithData for write response as HTTP OK result with data
func WriteAPIOKWithData(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}

// WriteAPICreated for write response as HTTP Created result
func WriteAPICreated(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// WriteAPINoContent for write response as HTTP NoContent result
func WriteAPINoContent(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// WriteApplicationJSON ...
func WriteApplicationJSON(w http.ResponseWriter, httpCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	_ = json.NewEncoder(w).Encode(body)
}

// WriteAPIError for write response as HTTP Error result
func WriteAPIError(w http.ResponseWriter, response APIResponse, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.HTTPCode)
	_ = json.NewEncoder(w).Encode(response.WithMessage(fmt.Sprint(err)))
}

// WriteAPIErrorMessage for write response as HTTP Error result
func WriteAPIErrorMessage(w http.ResponseWriter, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.HTTPCode)
	_ = json.NewEncoder(w).Encode(response)
}

// WriteAPIErrorWithData for write response as HTTP Error result with data
func WriteAPIErrorWithData(w http.ResponseWriter, response APIResponse, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.HTTPCode)
	_ = json.NewEncoder(w).Encode(data)
}

// Write writes the data to http response writer
func Write(w http.ResponseWriter, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.HTTPCode)
	_ = json.NewEncoder(w).Encode(response)
}

// WithMessage modifies api response's message
func (a *APIResponse) WithMessage(message interface{}) APIResponse {
	new := new(APIResponse)
	new.HTTPCode = a.HTTPCode
	new.Code = a.Code
	new.Message = message
	new.Data = a.Data

	return *new
}

// WithHTTPCode modifies api response's http code
func (a *APIResponse) WithHTTPCode(httpCode int) APIResponse {
	new := new(APIResponse)
	new.HTTPCode = httpCode
	new.Code = a.Code
	new.Message = a.Message
	new.Data = a.Data

	return *new
}

// WithData if have data to response
func (a *APIResponse) WithData(data interface{}) APIResponse {
	new := new(APIResponse)
	new.HTTPCode = a.HTTPCode
	new.Code = a.Code
	new.Message = a.Message
	new.Data = data

	return *new
}
