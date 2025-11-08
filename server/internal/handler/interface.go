package handler

import "net/http"

type ChatHandlerInterface interface {
	HandleChat(w http.ResponseWriter, r *http.Request)
}
