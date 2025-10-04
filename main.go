package main

import (
	"fmt"
	"net/http"

	"go-rest/config"
	"go-rest/handlers"
	"go-rest/utils"
	

	httpSwagger "github.com/swaggo/http-swagger"
	_ "go-rest/docs" // hasil generate swag
)

func setupRoutes() {
	// 🔹 Swagger docs
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	// 🔹 Serve static files (CSS, JS, gambar, dll)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))
	// Serve HTML
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/index.html")
	})
	// 🔹 Auth endpoints
	http.HandleFunc("/api/login", handlers.LoginHandler)
	http.HandleFunc("/api/logout", utils.Secure(handlers.LogoutHandler))
	http.HandleFunc("/api/refresh", handlers.RefreshHandler)

	// 🔹 Characters CRUD (secured)
	// GET all & POST
	http.HandleFunc("/api/characters", utils.Secure(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetCharacters(w, r)
		case http.MethodPost:
			handlers.CreateCharacter(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// GET by ID, PUT, DELETE
	http.HandleFunc("/api/characters/", utils.Secure(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetCharacterByID(w, r)
		case http.MethodPut:
			handlers.UpdateCharacter(w, r)
		case http.MethodDelete:
			handlers.DeleteCharacter(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// 🔹 API not found fallback
	http.HandleFunc("/api/", handlers.ApiNotFoundHandler)
}

// @title           Game Characters REST API
// @version         1.0
// @description     Dokumentasi API untuk sistem karakter game
// @termsOfService  http://example.com/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	pool, err := config.InitDB()
	if err != nil {
		fmt.Println("❌ Gagal konek ke database:", err)
		return
	}
	defer pool.Close()

	if err := config.RunMigration(pool); err != nil {
		fmt.Println("❌ Gagal menjalankan migration:", err)
		return
	}

	if err := utils.LoadUsersFromYAML("config.yaml"); err != nil {
		fmt.Println("⚠️ Gagal baca config.yaml:", err)
	}

	setupRoutes()

	fmt.Println("🚀 Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("❌ Server error:", err)
	}
}
