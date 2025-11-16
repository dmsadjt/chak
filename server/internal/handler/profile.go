package handler

import (
	"chak-server/internal/config"
	"encoding/json"
	"net/http"
)

type ProfileHandler struct {
	configManager config.ConfigInterface
}

type SwitchProfileRequest struct {
	SzProfileName string `json:"profile_name"`
}

func NewProfileHandler(cfgMgr config.ConfigInterface) *ProfileHandler {
	return &ProfileHandler{
		configManager: cfgMgr,
	}
}

func (prfHandler *ProfileHandler) HandleListProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	profiles := prfHandler.configManager.ListProfile()
	json.NewEncoder(w).Encode(profiles)
}

func (prfHandler *ProfileHandler) GetActiveProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	profile := prfHandler.configManager.GetActiveProfile()
	json.NewEncoder(w).Encode(profile)
}

func (prfHandler *ProfileHandler) HandleSwitchProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var req SwitchProfileRequest 
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if err := prfHandler.configManager.SwitchProfile(req.SzProfileName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile := prfHandler.configManager.GetActiveProfile()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":"success",
		"active_profile":profile,
	})
}
