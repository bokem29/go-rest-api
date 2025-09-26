package main

import (
	"fmt"
	"go-rest/handlers"
	"go-rest/utils"
	"net/http"
)

func setupRoutes() {

	// Static files & frontend
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("."))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// Auth endpoints
	http.HandleFunc("/api/login", handlers.LoginHandler)
	http.HandleFunc("/api/logout", utils.Secure(handlers.LogoutHandler))
	http.HandleFunc("/api/refresh", handlers.RefreshHandler)

	// API endpoints (secured)
	http.HandleFunc("/api/characters", utils.Secure(handlers.CharacterHandler))
	http.HandleFunc("/api/characters/", utils.Secure(handlers.CharacterHandler))

	// API not found fallback: must be registered AFTER specific /api routes
	http.HandleFunc("/api/", handlers.ApiNotFoundHandler)
}

func main() {
	// load dummy data (characters.json)
	utils.LoadData("characters.json")

	// load users from YAML config
	if err := utils.LoadUsersFromYAML("config.yaml"); err != nil {
		fmt.Println("Gagal baca config.yaml:", err)
	}

	setupRoutes()

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
