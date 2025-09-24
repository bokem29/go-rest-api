package models

//struktur data untuk menyimpan info karakter game
type Character struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	Game string `json:"game"`
}
