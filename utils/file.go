package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"go-rest/models"
)

var Characters []models.Character
var LastID int

// LoadData membaca data karakter dari file JSON
func LoadData(filename string) {
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Gagal baca file:", err)
		return
	}
	err = json.Unmarshal(file, &Characters)
	if err != nil {
		fmt.Println("Gagal decode JSON:", err)
		return
	}

	// Update LastID untuk auto-increment
	for _, c := range Characters {
		if c.ID > LastID {
			LastID = c.ID
		}
	}
}

// SaveData menyimpan data karakter ke file JSON
func SaveData(filename string) {
	data, err := json.MarshalIndent(Characters, "", "  ")
	if err != nil {
		fmt.Println("Gagal encode JSON:", err)
		return
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println("Gagal simpan file:", err)
	}
}
