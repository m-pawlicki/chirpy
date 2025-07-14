package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

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

func (apiCfg *APIHandler) ValidateHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {

		tooLongResp := struct {
			Err string `json:"error"`
		}{
			Err: "Chirp is too long",
		}

		RespondWithJSON(w, 400, tooLongResp)
		return
	}

	cleanedBody := checkProfanity(params.Body)

	cleanChirp := struct {
		Cleaned string `json:"cleaned_body"`
	}{
		Cleaned: cleanedBody,
	}

	RespondWithJSON(w, 200, cleanChirp)
}
