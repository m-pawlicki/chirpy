package handlers

import (
	"net/http"
	"time"

	"github.com/m-pawlicki/chirpy/internal/auth"
)

func (apiCfg *APIHandler) RefreshAccessTokenHandler(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	secret := apiCfg.Config.Secret
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}
	usr, err := apiCfg.Config.DB.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}
	tkn, err := apiCfg.Config.DB.LookupRefreshToken(r.Context(), usr.Token)
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}
	if tkn.RevokedAt.Valid {
		RespondWithError(w, 401, "revoked token")
		return
	}
	if tkn.ExpiresAt.Before(time.Now().UTC()) {
		RespondWithError(w, 401, "expired token")
		return
	}
	newToken, err := auth.MakeJWT(usr.ID, secret)
	if err != nil {
		RespondWithError(w, 401, err.Error())
	}
	RespondWithJSON(w, 200, response{Token: newToken})

}

func (apiCfg *APIHandler) RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}
	err = apiCfg.Config.DB.SetRefreshTokenAsRevoked(r.Context(), token)
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}
	RespondWithMsg(w, 204, "")
}
