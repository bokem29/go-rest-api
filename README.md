# Go REST API - Character Management

Sebuah REST API sederhana yang dibangun dengan Go untuk mengelola data karakter game. API ini menyediakan operasi CRUD (Create, Read, Update, Delete) untuk karakter game dengan penyimpanan data berbasis file JSON.

## 🚀 Fitur

- **CRUD Operations**: Create, Read, Update, Delete karakter game
- **JSON Storage**: Data disimpan dalam file JSON (`characters.json`)
- **RESTful API**: Endpoint yang mengikuti standar REST
- **Auto-increment ID**: ID otomatis untuk karakter baru
- **Static File Serving**: Melayani file statis untuk frontend
- **Request Logging**: Log setiap request dengan timestamp

## 📁 Struktur Proyek

```
go-rest/
├── main.go                 # Entry point aplikasi
├── go.mod                  # Go module file
├── characters.json         # Database file (JSON)
├── index.html             # Frontend interface
├── handlers/
│   └── characterHandler.go # Handler untuk operasi karakter
├── models/
│   └── models.go          # Struktur data Character
└── utils/
    └── file.go            # Utility functions untuk file operations
```


## 📚 API Endpoints

### Base URL
```
http://localhost:8080
```

### Endpoints

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/api/characters` | Mendapatkan semua karakter |
| `GET` | `/api/characters/{id}` | Mendapatkan karakter berdasarkan ID |
| `POST` | `/api/characters` | Membuat karakter baru |
| `PUT` | `/api/characters/{id}` | Mengupdate karakter berdasarkan ID |
| `DELETE` | `/api/characters/{id}` | Menghapus karakter berdasarkan ID |

### Contoh Request/Response

#### 1. Mendapatkan Semua Karakter
```bash
GET /api/characters
```

**Response:**
```json
[
  {
    "id": 2,
    "name": "Mario",
    "role": "Main Char",
    "game": "Super Mario"
  },
  {
    "id": 4,
    "name": "Atma",
    "role": "Main Char",
    "game": "A Space for The Unbound"
  }
]
```

#### 2. Mendapatkan Karakter Berdasarkan ID
```bash
GET /api/characters/2
```

**Response:**
```json
{
  "id": 2,
  "name": "Mario",
  "role": "Main Char",
  "game": "Super Mario"
}
```

#### 3. Membuat Karakter Baru
```bash
POST /api/characters
Content-Type: application/json

{
  "name": "Link",
  "role": "Hero",
  "game": "The Legend of Zelda"
}
```

**Response:**
```json
{
  "id": 8,
  "name": "Link",
  "role": "Hero",
  "game": "The Legend of Zelda"
}
```

#### 4. Mengupdate Karakter
```bash
PUT /api/characters/2
Content-Type: application/json

{
  "name": "Mario Bros",
  "role": "Main Character",
  "game": "Super Mario Bros"
}
```

#### 5. Menghapus Karakter
```bash
DELETE /api/characters/2
```

**Response:** `204 No Content`

## 🏗️ Arsitektur

### Models
```go
type Character struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Role string `json:"role"`
    Game string `json:"game"`
}
```

### Handler Functions
- `CharacterHandler`: Menangani semua operasi CRUD untuk karakter
- Mendukung HTTP methods: GET, POST, PUT, DELETE
- Validasi input dan error handling
- Auto-increment ID untuk karakter baru

### Utils
- `LoadData()`: Memuat data dari file JSON
- `SaveData()`: Menyimpan data ke file JSON
- Global variables: `Characters` dan `LastID`

## 🔧 Pengembangan

### Menambahkan Fitur Baru
1. Tambahkan handler baru di folder `handlers/`
2. Update routing di `main.go`
3. Tambahkan model baru di folder `models/` jika diperlukan

### Testing API
Gunakan tools seperti:
- **Postman**
- **curl**
- **Insomnia**
- **Thunder Client** (VS Code extension)

### Contoh curl commands:
```bash
# Get all characters
curl http://localhost:8080/api/characters

# Get character by ID
curl http://localhost:8080/api/characters/2

# Create new character
curl -X POST http://localhost:8080/api/characters \
  -H "Content-Type: application/json" \
  -d '{"name":"Pikachu","role":"Pokemon","game":"Pokemon"}'

# Update character
curl -X PUT http://localhost:8080/api/characters/2 \
  -H "Content-Type: application/json" \
  -d '{"name":"Mario Updated","role":"Hero","game":"Super Mario Bros"}'

# Delete character
curl -X DELETE http://localhost:8080/api/characters/2
```

## 📝 Catatan

- Data disimpan dalam file `characters.json`
- Server berjalan di port `8080`
- Log request ditampilkan di console
- File statis dilayani dari root directory
- ID otomatis increment untuk karakter baru



## 👨‍💻 Author

Dibuat dengan ❤️ menggunakan Go
