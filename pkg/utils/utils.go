package utils

import (
	"encoding/json"
	"net/http"
)

func Message(status bool, message interface{}) map[string]interface{} {
	if status {
		return map[string]interface{}{"status": "OK", "message": message}
	}
	return map[string]interface{}{"status": "error", "message": message}

}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func Json(w http.ResponseWriter, httpCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(&data)
}

func Text(w http.ResponseWriter, httpCode int, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(httpCode)
	w.Write([]byte(message))
}

func Err(w http.ResponseWriter, httpCode int, err error) {

	w.Header().Set("Content-Type", "application/json")
	//need more error status
	w.WriteHeader(httpCode)
	res := Message(false, err.Error())
	json.NewEncoder(w).Encode(res)
}
