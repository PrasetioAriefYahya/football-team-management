## Cara Menjalankan Aplikasi (How to Run)

1. Pastikan PostgreSQL sudah terinstal dan berjalan.
2. Buat database bernama `football`.
3. Sesuaikan konfigurasi koneksi pada variabel `dsn` di file `main.go`:
   ```go
   dsn := "host=localhost user=postgres password=dawad123321 dbname=football port=5432 sslmode=disable"
   ```
4. Jalankan aplikasi dengan perintah:
   ```bash
   go run main.go
   ```
5. Server akan berjalan pada `http://localhost:8080`.

---

## Autentikasi (Authentication)

### Endpoint
```
POST /auth/login
```

### Deskripsi
Login sebagai admin untuk mendapatkan token JWT.

### Contoh Request
```json
{
  "username": "admin",
  "password": "admin123"
}
```

### Contoh Response
```json
{
  "token": "your_jwt_token_here"
}
```

Gunakan token ini untuk mengakses endpoint yang membutuhkan autentikasi dengan menambahkan header:
```
Authorization: Bearer your_jwt_token_here
```

---

## Endpoint API

### 1. Teams (Tim)

#### GET /teams
Menampilkan seluruh data tim. (Tidak memerlukan autentikasi)

**Response:**
```json
[
  {
    "ID": 1,
    "name": "Garuda FC",
    "logo": "garuda.png",
    "founded": 2005,
    "address": "Jl. Merdeka 10",
    "city": "Jakarta",
    "Players": []
  }
]
```

#### POST /teams
Menambahkan tim baru. (Memerlukan JWT)

**Request:**
```json
{
  "name": "Garuda FC",
  "logo": "garuda.png",
  "founded": 2005,
  "address": "Jl. Merdeka 10",
  "city": "Jakarta"
}
```

**Response:**
```json
{
  "ID": 1,
  "name": "Garuda FC",
  "logo": "garuda.png",
  "founded": 2005,
  "address": "Jl. Merdeka 10",
  "city": "Jakarta"
}
```

#### PUT /teams/:id
Mengubah data tim berdasarkan ID.

**Request:**
```json
{
  "name": "Garuda United"
}
```

**Response:**
```json
{
  "ID": 1,
  "name": "Garuda United"
}
```

#### DELETE /teams/:id
Menghapus tim (soft delete).

**Response:**
```json
{
  "message": "Team soft deleted"
}
```

---

### 2. Players (Pemain)

#### POST /players
Menambahkan pemain ke dalam tim.

**Request:**
```json
{
  "name": "Budi Santoso",
  "height": 175,
  "weight": 68,
  "position": "Penyerang",
  "number": 10,
  "team_id": 1
}
```

**Response:**
```json
{
  "ID": 1,
  "name": "Budi Santoso",
  "height": 175,
  "weight": 68,
  "position": "Penyerang",
  "number": 10,
  "team_id": 1
}
```

---

### 3. Matches (Pertandingan)

#### POST /matches
Menambahkan jadwal pertandingan antar tim.

**Request:**
```json
{
  "date": "2025-10-20",
  "time": "15:00",
  "home_team_id": 1,
  "away_team_id": 2
}
```

**Response:**
```json
{
  "ID": 1,
  "date": "2025-10-20",
  "time": "15:00",
  "home_team_id": 1,
  "away_team_id": 2
}
```

---

### 4. Match Results (Hasil Pertandingan)

#### POST /matches/result
Melaporkan hasil pertandingan.

**Request:**
```json
{
  "match_id": 1,
  "home_score": 2,
  "away_score": 1,
  "scorer_name": "Budi Santoso",
  "goal_minute": 55
}
```

**Response:**
```json
{
  "ID": 1,
  "match_id": 1,
  "home_score": 2,
  "away_score": 1,
  "scorer_name": "Budi Santoso",
  "goal_minute": 55
}
```

---

### 5. Reports (Laporan Hasil Pertandingan)

#### GET /reports
Menampilkan laporan hasil pertandingan, termasuk skor, status akhir, pencetak gol terbanyak, dan total kemenangan tim.

**Response:**
```json
[
  {
    "match_id": 1,
    "home_team": "Garuda FC",
    "away_team": "Macan FC",
    "home_score": 2,
    "away_score": 1,
    "status": "Home Win",
    "top_scorer": "Budi Santoso",
    "home_total_wins": 1,
    "away_total_wins": 0
  }
]
```

---

## English Summary

### How to Run
1. Ensure PostgreSQL is running and create a database named `football`.
2. Edit the connection string in `main.go` if needed.
3. Run the application:
   ```bash
   go run main.go
   ```
4. Server runs on `http://localhost:8080`.

### Authentication
`POST /auth/login` → Returns a JWT token.

Use this token in every protected endpoint:
```
Authorization: Bearer <token>
```

### Available Endpoints
| Entity | Method | Endpoint | Description |
|---------|---------|-----------|-------------|
| Teams | GET | /teams | List all teams |
| Teams | POST | /teams | Create a new team |
| Teams | PUT | /teams/:id | Update team info |
| Teams | DELETE | /teams/:id | Soft delete a team |
| Players | POST | /players | Add a player to a team |
| Matches | POST | /matches | Create match schedule |
| Results | POST | /matches/result | Submit match result |
| Reports | GET | /reports | Get match summary report |

---

**Default admin login:**  
Username: `admin`  
Password: `admin123`
