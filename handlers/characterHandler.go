package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "strconv"
    "strings"

    "go-rest/config"
    "go-rest/models"
)

// ✅ GET All Characters
// @Summary      Ambil semua karakter game
// @Description  Mendapatkan list semua karakter dari database
// @Tags         characters
// @Produce      json
// @Success      200  {array}   models.Character
// @Failure      500  {object}  map[string]string
// @Router       /characters [get]
func GetCharacters(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    rows, err := config.DB.Query(ctx, "SELECT id, name, role, game FROM characters")
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var characters []models.Character
    for rows.Next() {
        var c models.Character
        if err := rows.Scan(&c.ID, &c.Name, &c.Role, &c.Game); err != nil {
            http.Error(w, "Error scanning data", http.StatusInternalServerError)
            return
        }
        characters = append(characters, c)
    }

    json.NewEncoder(w).Encode(characters)
}

// ✅ GET Character by ID
// @Summary      Ambil karakter berdasarkan ID
// @Tags         characters
// @Produce      json
// @Param        id   path      int  true  "Character ID"
// @Success      200  {object}  models.Character
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /characters/{id} [get]
func GetCharacterByID(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    idStr := strings.TrimPrefix(r.URL.Path, "/api/characters/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var c models.Character
    err = config.DB.QueryRow(ctx,
        "SELECT id, name, role, game FROM characters WHERE id=$1", id,
    ).Scan(&c.ID, &c.Name, &c.Role, &c.Game)

    if err != nil {
        http.Error(w, "Character not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(c)
}

// ✅ CREATE Character
// @Summary      Tambah karakter baru
// @Tags         characters
// @Accept       json
// @Produce      json
// @Param        character  body      models.Character  true  "Character Data"
// @Success      201  {object}  models.Character
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /characters [post]
func CreateCharacter(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    var character models.Character
    if err := json.NewDecoder(r.Body).Decode(&character); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    err := config.DB.QueryRow(ctx,
        "INSERT INTO characters (name, role, game) VALUES ($1, $2, $3) RETURNING id",
        character.Name, character.Role, character.Game,
    ).Scan(&character.ID)

    if err != nil {
        http.Error(w, "Failed to insert", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(character)
}

// ✅ UPDATE Character
// @Summary      Update karakter
// @Tags         characters
// @Accept       json
// @Produce      json
// @Param        id         path      int                true  "Character ID"
// @Param        character  body      models.Character  true  "Character Data"
// @Success      200  {object}  models.Character
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /characters/{id} [put]
func UpdateCharacter(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    idStr := strings.TrimPrefix(r.URL.Path, "/api/characters/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var character models.Character
    if err := json.NewDecoder(r.Body).Decode(&character); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    _, err = config.DB.Exec(ctx,
        "UPDATE characters SET name=$1, role=$2, game=$3 WHERE id=$4",
        character.Name, character.Role, character.Game, id,
    )
    if err != nil {
        http.Error(w, "Failed to update", http.StatusInternalServerError)
        return
    }

    character.ID = id
    json.NewEncoder(w).Encode(character)
}

// ✅ DELETE Character
// @Summary      Hapus karakter
// @Tags         characters
// @Param        id   path      int  true  "Character ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /characters/{id} [delete]
func DeleteCharacter(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    idStr := strings.TrimPrefix(r.URL.Path, "/api/characters/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    _, err = config.DB.Exec(ctx, "DELETE FROM characters WHERE id=$1", id)
    if err != nil {
        http.Error(w, "Failed to delete", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
