package models

import "time"

// Struktur data untuk menyimpan info karakter game
type Character struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Role      string     `json:"role"`
	Game      string     `json:"game"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
