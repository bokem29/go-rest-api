package main

import (
	"fmt"
	"net/http"
	"go-rest/handlers"
	"go-rest/utils"
)

// setupRoutes mengatur routing API
func setupRoutes() {
	http.HandleFunc("/api/characters", handlers.CharacterHandler)
	http.HandleFunc("/api/characters/", handlers.CharacterHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("."))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
}

func main() {
	utils.LoadData("characters.json")
	setupRoutes()

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
