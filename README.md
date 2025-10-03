# Go Gin Realtime Chat

Aplikasi chat real-time berbasis **Go** dan **Gin Framework**, memanfaatkan WebSocket untuk komunikasi langsung antar klien. Dilengkapi dengan backend, routing, dan contoh client minimal.

---

## 🚀 Fitur Utama

- Komunikasi real-time dengan WebSocket  
- Routing dan pengelolaan pesan dengan Gin  
- Struktur modular: `controllers`, `routes`, `middlewares`, `models`, `utils`  
- Mendukung deployment via Docker + `docker-compose`  
- Contoh client HTML sederhana untuk testing  

---

## 📁 Struktur Proyek

```
Go-Gin-Realtime-Chat/
├── controllers/      # Logic pengelolaan WebSocket & request
├── middlewares/      # Middleware (misal auth, CORS)
├── models/           # Model data (pesan, user, etc)
├── routes/           # Definisi route API / WebSocket
├── utils/            # Helper, utilitas
├── docs/             # Dokumentasi / spesifikasi API (jika ada)
├── docker-compose.yml
├── Dockerfile
├── main.go
├── hub.go             # manajemen WebSocket hub
├── client.go          # contoh client Go (opsional)
├── index.html         # contoh client berbasis browser
└── .gitignore
```

---

## ⚙️ Instalasi & Jalankan

### 1. Clone Repo
```bash
git clone https://github.com/AhmadFazriNursamsi/Go-Gin-Realtime-Chat.git
cd Go-Gin-Realtime-Chat
```

### 2. Install dependencies
```bash
go mod tidy
```

### 3. (Opsional) Setup `.env` kalau ada variabel konfigurasi (misalnya port, DB, dsb)  
Buat `.env` berdasarkan contoh jika ada file `.env.example`.

### 4. Jalankan aplikasi
```bash
go run main.go
```

Atau jika kamu menggunakan Docker + docker-compose:
```bash
docker-compose up --build
```

---

## 🔌 Endpoint & Contoh Penggunaan

- `GET /` → Menampilkan halaman client (HTML) untuk chat  
- WebSocket endpoint (misalnya `/ws`) → Untuk koneksi real-time  
- Kirim / terima pesan antar klien  

Contoh client sederhana disertakan: `index.html`  
Buka di browser dan sambungkan ke server WebSocket untuk mencoba.

---

## 🧪 Testing

Coba buka beberapa tab browser ke `index.html`, kirim pesan dari satu tab → pesan muncul di semua tab lain.

Bisa juga menggunakan client websocket Go (`client.go`) sebagai simulasi klien.

---

## 📝 Catatan Pengembangan

- Pastikan port WebSocket tidak bertabrakan  
- Kelola hub/client connection dengan baik agar tidak ada memory leak  
- Jika ingin menambahkan auth, cukup tambahkan middleware JWT sebelum upgrade ke WebSocket  
- Untuk skala besar, pertimbangkan Redis Pub/Sub agar WebSocket bisa horizontal scale  

---

## 📜 Lisensi

Aplikasi ini dirilis di bawah lisensi **MIT License**.  
(Masukkan file `LICENSE` di root repo)
