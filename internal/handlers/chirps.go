package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/m-pawlicki/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (apiCfg *APIHandler) GetChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := apiCfg.Config.DB.GetChirps(r.Context())
	if err != nil {
		RespondWithError(w, 500, "Couldn't get chirps")
		return
	}
	var chirpList []Chirp
	for _, chirp := range chirps {
		chirpList = append(chirpList,
			Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID})
	}
	RespondWithJSON(w, 200, chirpList)
}

func (apiCfg *APIHandler) PostChirpHandler(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	pl := payload{}
	err := decoder.Decode(&pl)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	chirpBody, isValid := validateChirp(pl.Body)
	if isValid {
		newChirp := database.CreateChirpParams{
			Body:   chirpBody,
			UserID: pl.UserID,
		}
		res, err := apiCfg.Config.DB.CreateChirp(r.Context(), newChirp)
		if err != nil {
			RespondWithError(w, 500, err.Error())
			return
		} else {
			RespondWithJSON(w, 201, Chirp{ID: res.ID, CreatedAt: res.CreatedAt, UpdatedAt: res.UpdatedAt, Body: res.Body, UserID: res.UserID})
			return
		}

	} else {
		RespondWithError(w, 400, "Chirp is too long")
		return
	}
}

func checkProfanity(body string) string {
	invalidWords := [...]string{"kerfuffle", "sharbert", "fornax"}
	splitStr := strings.Split(body, " ")
	for i, word := range splitStr {
		for _, val := range invalidWords {
			if strings.ToLower(word) == val {
				splitStr[i] = "****"
			}
		}
	}
	cleaned := strings.Join(splitStr, " ")
	return cleaned
}

func validateChirp(body string) (string, bool) {

	if len(body) > 140 {
		return "", false
	}

	cleanedBody := checkProfanity(body)

	return cleanedBody, true
}
