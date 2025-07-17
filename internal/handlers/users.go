package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/m-pawlicki/chirpy/internal/auth"
	"github.com/m-pawlicki/chirpy/internal/database"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (apiCfg *APIHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	type payload struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	pl := payload{}
	err := decoder.Decode(&pl)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	hash, err := auth.HashPassword(pl.Password)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	userParams := database.CreateUserParams{
		Email:          pl.Email,
		HashedPassword: hash,
	}

	user, err := apiCfg.Config.DB.CreateUser(r.Context(), userParams)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	} else {
		RespondWithJSON(w, 201, User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email, IsChirpyRed: user.IsChirpyRed})
	}
}

func (apiCfg *APIHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	godotenv.Load(".env")
	secret := os.Getenv("SECRET")
	decoder := json.NewDecoder(r.Body)
	resp := response{}
	err := decoder.Decode(&resp)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	usr, err := apiCfg.Config.DB.FindUserByEmail(r.Context(), resp.Email)
	if err != nil {
		RespondWithError(w, 401, "Incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(resp.Password, usr.HashedPassword)
	if err != nil {
		RespondWithError(w, 401, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(usr.ID, secret)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	refresh, err := auth.MakeRefreshToken()
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	rt, err := apiCfg.Config.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{Token: refresh, UserID: usr.ID})
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	RespondWithJSON(w, 200, User{ID: usr.ID, CreatedAt: usr.CreatedAt, UpdatedAt: usr.UpdatedAt, Email: usr.Email, Token: token, RefreshToken: rt.Token, IsChirpyRed: usr.IsChirpyRed})
}

func (apiCfg *APIHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	godotenv.Load(".env")
	secret := os.Getenv("SECRET")
	type response struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}
	decoder := json.NewDecoder(r.Body)
	resp := response{}
	err = decoder.Decode(&resp)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	uid, err := auth.ValidateJWT(token, secret)
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}
	newPw, err := auth.HashPassword(resp.Password)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	refreshToken, err := apiCfg.Config.DB.GetRefreshTokenFromUser(r.Context(), uid)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	updateParams := database.UpdateUserParams{
		ID:             uid,
		Email:          resp.Email,
		HashedPassword: newPw,
	}
	updatedUser, err := apiCfg.Config.DB.UpdateUser(r.Context(), updateParams)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	RespondWithJSON(w, 200, User{ID: updatedUser.ID, CreatedAt: updatedUser.CreatedAt, UpdatedAt: updatedUser.UpdatedAt, Email: updatedUser.Email, Token: token, RefreshToken: refreshToken.Token, IsChirpyRed: updatedUser.IsChirpyRed})
}

func (apiCfg *APIHandler) DeleteUsersHandler(w http.ResponseWriter, r *http.Request) {

	if apiCfg.Config.Platform == "dev" {
		err := apiCfg.Config.DB.DeleteUsers(r.Context())
		if err != nil {
			RespondWithError(w, 500, err.Error())
		} else {
			RespondWithMsg(w, 200, "Users reset")
		}
	} else {
		RespondWithError(w, 403, "Forbidden")
	}
}
