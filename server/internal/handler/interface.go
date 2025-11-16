package handler

import "net/http"

type ChatHandlerInterface interface {
	HandleChat(w http.ResponseWriter, r *http.Request)
}

type ProfileHandlerInterface interface {
	HandleListProfiles(w http.ResponseWriter, r *http.Request)
	HandleGetActiveProfile(w http.ResponseWriter, r *http.Request)
	HandleSwitchProfile(w http.ResponseWriter, r *http.Request)
}
