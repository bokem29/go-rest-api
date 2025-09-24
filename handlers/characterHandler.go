package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"go-rest/models"
	"go-rest/utils"
)

func CharacterHandler(w http.ResponseWriter, r *http.Request) {
	// Log method + path + timestamp
	fmt.Printf("[%s] %s %s\n", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "application/json")
	path := strings.TrimPrefix(r.URL.Path, "/api/characters")

	switch r.Method {
	case "GET":
		if path == "" {
			json.NewEncoder(w).Encode(utils.Characters)
		} else {
			id, err := strconv.Atoi(strings.TrimPrefix(path, "/"))
			if err != nil {
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}
			for _, c := range utils.Characters {
				if c.ID == id {
					json.NewEncoder(w).Encode(c)
					return
				}
			}
			http.NotFound(w, r)
		}

	case "POST":
		var character models.Character
		if err := json.NewDecoder(r.Body).Decode(&character); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		utils.LastID++
		character.ID = utils.LastID
		utils.Characters = append(utils.Characters, character)

		utils.SaveData("characters.json")
		json.NewEncoder(w).Encode(character)

	case "PUT":
		id, err := strconv.Atoi(strings.TrimPrefix(path, "/"))
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		for i, c := range utils.Characters {
			if c.ID == id {
				var character models.Character
				if err := json.NewDecoder(r.Body).Decode(&character); err != nil {
					http.Error(w, "Invalid JSON", http.StatusBadRequest)
					return
				}
				character.ID = id
				utils.Characters[i] = character
				utils.SaveData("characters.json")
				json.NewEncoder(w).Encode(character)
				return
			}
		}
		http.NotFound(w, r)

	case "DELETE":
		id, err := strconv.Atoi(strings.TrimPrefix(path, "/"))
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		for i, c := range utils.Characters {
			if c.ID == id {
				utils.Characters = append(utils.Characters[:i], utils.Characters[i+1:]...)
				utils.SaveData("characters.json")
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.NotFound(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
