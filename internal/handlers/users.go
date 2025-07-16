package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/m-pawlicki/chirpy/internal/auth"
	"github.com/m-pawlicki/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
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
		RespondWithJSON(w, 201, User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email})
	}
}

func (apiCfg *APIHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Password  string `json:"password"`
		Email     string `json:"email"`
		ExpiresIn string `json:"expires_in_seconds"`
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
	var expiry time.Duration
	if resp.ExpiresIn == "" {
		expiry = time.Hour
	} else {
		conv, err := strconv.Atoi(resp.ExpiresIn)
		if err != nil {
			RespondWithError(w, 401, "Malformed expires_in_seconds")
			return
		}
		if conv > 3600 {
			expiry = time.Hour

		} else {
			expiry = (time.Second * time.Duration(conv))
		}
	}
	token, err := auth.MakeJWT(usr.ID, secret, expiry)
	if err != nil {
		RespondWithError(w, 401, "Authenication failure")
		return
	}
	RespondWithJSON(w, 200, User{ID: usr.ID, CreatedAt: usr.CreatedAt, UpdatedAt: usr.UpdatedAt, Email: usr.Email, Token: token})
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
