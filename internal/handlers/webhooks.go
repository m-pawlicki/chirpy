package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/m-pawlicki/chirpy/internal/auth"
)

func (apiCfg *APIHandler) UpgradeUserToRedHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}
	if key != apiCfg.Config.APIKey {
		RespondWithError(w, 401, "Unauthorized")
		return
	}
	decoder := json.NewDecoder(r.Body)
	resp := response{}
	err = decoder.Decode(&resp)
	if err != nil {
		RespondWithMsg(w, 500, "")
		return
	}
	if resp.Event != "user.upgraded" {
		RespondWithMsg(w, 204, "")
		return
	}
	id, err := uuid.Parse(resp.Data.UserID)
	if err != nil {
		RespondWithMsg(w, 500, "")
		return
	}
	usr, err := apiCfg.Config.DB.FindUserByID(r.Context(), id)
	if err != nil {
		RespondWithError(w, 404, "User not found")
	}
	err = apiCfg.Config.DB.UpgradeUserToRed(r.Context(), usr.ID)
	if err != nil {
		RespondWithMsg(w, 500, "")
		return
	}
	RespondWithMsg(w, 204, "")
}
