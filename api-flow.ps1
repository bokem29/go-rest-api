# Login
$resp = Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/login" -ContentType "application/json" -Body (@{username="user";password="pass123"} | ConvertTo-Json)
$env:API_TOKEN = $resp.token
# Token khusus yang lebih panjang umurnya(Opsional)
$env:REFRESH_TOKEN = $resp.refresh

# Get API
Invoke-RestMethod -Method Get -Uri "http://localhost:8080/api/characters" `
  -Headers @{ Authorization = "Bearer $env:API_TOKEN" }

# Tambah character
Invoke-RestMethod "http://localhost:8080/api/characters" `
  -Method POST `
  -Headers @{ "Content-Type" = "application/json"; Authorization = "Bearer $env:API_TOKEN" } `
  -Body '{"name":"Naruto Uzumaki","role":"Ninja","game":"Naruto Ultimate"}'

# Logout
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/logout" -Headers @{Authorization="Bearer $env:API_TOKEN"}

