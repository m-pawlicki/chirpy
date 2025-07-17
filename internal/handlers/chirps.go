package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/m-pawlicki/chirpy/internal/auth"
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
	authorID := r.URL.Query().Get("author_id")
	var chirpList []Chirp
	if authorID != "" {
		uid, err := uuid.Parse(authorID)
		if err != nil {
			RespondWithError(w, 401, "Couldn't parse author_id")
			return
		}
		chirps, err := apiCfg.Config.DB.GetChirpsByAuthorID(r.Context(), uid)
		if err != nil {
			RespondWithError(w, 500, "Couldn't get chirps")
		}
		for _, chirp := range chirps {
			chirpList = append(chirpList,
				Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID})
		}
	} else {
		chirps, err := apiCfg.Config.DB.GetChirps(r.Context())
		if err != nil {
			RespondWithError(w, 500, "Couldn't get chirps")
			return
		}
		for _, chirp := range chirps {
			chirpList = append(chirpList,
				Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID})
		}
	}
	sortBy := r.URL.Query().Get("sort")
	switch sortBy {
	case "asc":
		RespondWithJSON(w, 200, chirpList)
		return
	case "desc":
		sort.Slice(chirpList, func(i, j int) bool {
			return chirpList[i].CreatedAt.After(chirpList[j].CreatedAt)
		})
		RespondWithJSON(w, 200, chirpList)
		return
	default:
		RespondWithJSON(w, 200, chirpList)
	}
}

func (apiCfg *APIHandler) GetChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	pathVal := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(pathVal)
	if err != nil {
		RespondWithError(w, 400, "Invalid ID")
		return
	}
	chirp, err := apiCfg.Config.DB.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		RespondWithError(w, 404, "Chirp not found")
		return
	}
	RespondWithJSON(w, 200, Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID})
}

func (apiCfg *APIHandler) PostChirpHandler(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Body string `json:"body"`
	}
	secret := apiCfg.Config.Secret
	decoder := json.NewDecoder(r.Body)
	pl := payload{}
	err := decoder.Decode(&pl)
	if err != nil {
		RespondWithError(w, 500, "Couldn't decode body")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 401, "Unauthorized")
		return
	}
	jwtUsr, err := auth.ValidateJWT(token, secret)
	if err != nil {
		RespondWithError(w, 401, "Unauthorized")
		return
	}

	chirpBody, isValidChirp := validateChirp(pl.Body)
	if isValidChirp {
		newChirp := database.CreateChirpParams{
			Body:   chirpBody,
			UserID: jwtUsr,
		}
		res, err := apiCfg.Config.DB.CreateChirp(r.Context(), newChirp)
		if err != nil {
			RespondWithError(w, 500, err.Error())
			return
		}
		RespondWithJSON(w, 201, Chirp{ID: res.ID, CreatedAt: res.CreatedAt, UpdatedAt: res.UpdatedAt, Body: res.Body, UserID: res.UserID})
		return
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

func (apiCfg *APIHandler) DeleteChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	secret := apiCfg.Config.Secret
	pathVal := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(pathVal)
	if err != nil {
		RespondWithError(w, 400, "Invalid ID")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 401, "Unauthorized")
		return
	}
	usr, err := auth.ValidateJWT(token, secret)
	if err != nil {
		RespondWithError(w, 401, "Unauthorized")
		return
	}

	chirp, err := apiCfg.Config.DB.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		RespondWithError(w, 404, "Chirp not found")
		return
	}

	if chirp.UserID != usr {
		RespondWithError(w, 403, "Unauthorized")
		return
	}

	err = apiCfg.Config.DB.DeleteChirpByID(r.Context(), chirpID)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	RespondWithMsg(w, 204, "")
}
