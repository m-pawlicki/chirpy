package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (apiCfg *APIHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	user, err := apiCfg.Config.DB.CreateUser(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, 500, err.Error())
	} else {
		RespondWithJSON(w, 201, User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email})
	}
}

func (apiCfg *APIHandler) DeleteUsersHandler(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Body string `json:"body"`
	}

	if apiCfg.Config.Platform == "dev" {
		err := apiCfg.Config.DB.DeleteUsers(r.Context())
		if err != nil {
			RespondWithError(w, 500, err.Error())
		} else {
			RespondWithJSON(w, 200, response{Body: "Users reset."})
		}
	} else {
		RespondWithError(w, 403, "Forbidden")
	}
}
