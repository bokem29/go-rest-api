# Go REST API - Character Management

Sebuah REST API sederhana yang dibangun dengan Go untuk mengelola data karakter game. API ini menyediakan operasi CRUD (Create, Read, Update, Delete) untuk karakter game dengan penyimpanan data berbasis PostgreSQL.

API ini juga dilengkapi dengan Swagger API Documentation sehingga memudahkan pengembang untuk memahami dan menguji endpoint.

## 🚀 Fitur

- **CRUD Operations**: Create, Read, Update, Delete karakter game
- **Database Storage**: Data disimpan dalam database PostgreSQL (karakter_game)
- **RESTful API**: Endpoint yang mengikuti standar REST
- **Auto-increment ID**: ID otomatis untuk karakter baru
- **Static File Serving**: Melayani file statis untuk frontend
- **Request Logging**: Log setiap request dengan timestamp
// Fitur otentikasi & otorisasi
- **JWT Authentication**: Login menghasilkan access token (JWT)
- **Refresh Tokens**: Mendapatkan token baru tanpa login ulang
- **Logout (Token Revocation)**: Mencabut token via blacklist JTI sampai masa berlaku habis
- **Protected Routes**: Endpoint `/api/characters` diamankan dengan Bearer token
- **Config via YAML**: User di-load dari `config.yaml`; rahasia JWT via env `JWT_SECRET`
- **Cookie Fallback**: Server membaca token dari cookie `access_token` jika header Authorization tidak ada
- **API Fallback 404**: Rute `/api/*` yang tidak dikenali mengembalikan 404 JSON, bukan HTML

## 📁 Struktur Proyek

```
go-rest/
├── main.go                 # Entry point aplikasi
├── go.mod                  # Go module file
├── characters.json         # Database file (JSON)
├── config.yaml             # User accounts (demo/dev)
├── config/
│   ├── db.go               # Koneksi ke PostgreSQL
│   └── migration.go        # Script migrasi database
├── docs/
│   ├── docs.go             # Konfigurasi metadata dokumentasi Swagger ke aplikasi Go
│   ├── swagger.json        # Dokumentasi API (format JSON)
│   └── swagger.yaml        # Dokumentasi API (format YAML)
├── handlers/
│   ├── characterHandler.go # Handler untuk operasi karakter
│   ├── authHandler.go      # Handler untuk login/refresh/logout
│   └── apiFallback.go      # 404 JSON untuk rute /api/* yang tidak cocok
├── models/
│   └── models.go           # Struktur data Character
├── utils/
│   ├── file.go             # Utility functions untuk file operations
│   ├── auth.go             # Utilitas JWT, refresh store, extractor
│   └── middleware.go       # Middleware: Secure, RequestLogger, Recover
├── frontend/
│   ├── index.html          # Halaman utama frontend
│   ├── style.css           # File CSS
│   └── main.js             # File JavaScript

```


## 📚 API Endpoints

### Base URL
```
http://localhost:8080
```

### Endpoints

| Method | Endpoint | Deskripsi | Auth |
|--------|----------|-----------|------|
| `POST` | `/api/login` | Login, menghasilkan access + refresh token | No |
| `POST` | `/api/refresh` | Tukar refresh token untuk pasangan token baru | No |
| `POST` | `/api/logout` | Mencabut access token saat ini (blacklist JTI) | Bearer |
| `GET` | `/api/characters` | Mendapatkan semua karakter | Bearer |
| `GET` | `/api/characters/{id}` | Mendapatkan karakter berdasarkan ID | Bearer |
| `POST` | `/api/characters` | Membuat karakter baru | Bearer |
| `PUT` | `/api/characters/{id}` | Mengupdate karakter berdasarkan ID | Bearer |
| `DELETE` | `/api/characters/{id}` | Menghapus karakter berdasarkan ID | Bearer |

## 📜 Swagger API Documentation

API documentation otomatis tersedia menggunakan Swagger.
Setelah menjalankan server, akses dokumentasi di:

```bash
http://localhost:8080/docs/index.html
```

Swagger docs memuat:

- Daftar semua endpoint API

- Deskripsi endpoint

- Parameter request

- Response model

- Contoh request dan response


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

## 🔐 Otentikasi

### Login
```bash
$resp = Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/login" -ContentType "application/json" -Body (@{username="user";password="pass123"} | ConvertTo-Json)
$env:API_TOKEN = $resp.token
# Token khusus yang lebih panjang umurnya(Opsional)
$env:REFRESH_TOKEN = $resp.refresh
```
Cek Token yang disimpan dalam env:
```bash
$env:API_TOKEN
```

Get characters:
```bash
Invoke-RestMethod -Method Get -Uri "http://localhost:8080/api/characters" `
  -Headers @{ Authorization = "Bearer $env:API_TOKEN" }
```

Add characters:
```bash
Invoke-RestMethod "http://localhost:8080/api/characters" `
  -Method POST `
  -Headers @{ "Content-Type" = "application/json"; Authorization = "Bearer $env:API_TOKEN" } `
  -Body '{"name":"{new character name}","role":"{new role}","game":"{new game}"}'
```

Update characters:
```bash
Invoke-RestMethod -Method Put -Uri "http://localhost:8080/api/characters/1" `
  -Headers @{ Authorization = "Bearer $env:API_TOKEN"; "Content-Type" = "application/json" } `
  -Body '{"name":"{new name}","role":"{new role}","game":"{new game}"}'
    
```
Delete characters:
```bash
Invoke-RestMethod -Method Delete -Uri "http://localhost:8080/api/characters/{id}" `       
  -Headers @{ Authorization = "Bearer $env:API_TOKEN" }
    
```

### Logout
```bash
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/logout" -Headers @{Authorization="Bearer $env:API_TOKEN"}
```
Menggunakan blacklist berdasarkan `jti` pada JWT hingga masa berlaku habis.

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
- `Authenticate()`, `CreateToken()`, `Secure()`: Utilitas otentikasi JWT dan middleware

## 🔧 Pengembangan

### Menjalankan Aplikasi
1. (Opsional) Set rahasia JWT untuk produksi/deploy:
   - Windows PowerShell: `setx JWT_SECRET "your-strong-secret"`
   - Linux/macOS: `export JWT_SECRET="your-strong-secret"`
2. Pastikan `config.yaml` berisi user untuk login (contoh tersedia).
3. Jalankan server:
```bash
go run main.go
```
Server berjalan di `http://localhost:8080`.

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

- Data disimpan di PostgreSQL (karakter_game database, characters table)

- Pastikan koneksi ke database dikonfigurasi di config/db.go

- Dokumentasi API tersedia di /docs/index.html via Swagger

- Endpoint /api/characters memerlukan header Authorization: Bearer <ACCESS_JWT>


## 👨‍💻 Author

Dibuat dengan ❤️ menggunakan Go
