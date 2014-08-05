package server

import (
	"encoding/json"
	"net/http"
)

type Response map[string]interface{}

// ResponseSerializer is responsible for serializing REST responses and sending
// them back to the client.
type ResponseSerializer interface {
	SendSuccessResponse(http.ResponseWriter, Response, int)
	SendErrorResponse(http.ResponseWriter, error, int)
}

type JsonSerializer struct{}

func (j JsonSerializer) SendErrorResponse(w http.ResponseWriter, err error, errorCode int) {
	response := makeErrorResponse(err)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)
	w.Write(jsonResponse)
}

func (j JsonSerializer) SendSuccessResponse(w http.ResponseWriter, r Response, status int) {
	jsonResponse, err := json.Marshal(r)
	if err != nil {
		j.SendErrorResponse(w, err, 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonResponse)
}

func makeSuccessResponse(r interface{}, cursor string) Response {
	response := Response{
		"success": true,
		"result":  r,
	}

	if cursor != "" {
		response["next"] = cursor
	}

	return response
}

func makeErrorResponse(err error) Response {
	return Response{
		"success": false,
		"error":   err.Error(),
	}
}

func sendResponse(s ResponseSerializer, w http.ResponseWriter, r interface{}, cursor string, err error, status int) {
	if s == nil {
		// Fall back to json serialization.
		s = JsonSerializer{}
	}

	if err != nil {
		if status < 400 {
			status = http.StatusInternalServerError
		}
		s.SendErrorResponse(w, err, status)
		return
	}

	s.SendSuccessResponse(w, makeSuccessResponse(r, cursor), status)
}
