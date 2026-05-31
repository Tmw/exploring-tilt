package httphelper

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func RespondError(w http.ResponseWriter, code int, msg string) {
	resp := ErrorResponse{Code: code, Msg: msg}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}
