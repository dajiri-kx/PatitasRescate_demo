package shared

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	OK      bool        `json:"ok"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, status int, resp Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func JSONOk(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, Response{OK: true, Data: data})
}

func JSONMsg(w http.ResponseWriter, msg string) {
	JSON(w, http.StatusOK, Response{OK: true, Message: msg})
}

func JSONErr(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, Response{OK: false, Error: msg})
}

func DecodeBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
