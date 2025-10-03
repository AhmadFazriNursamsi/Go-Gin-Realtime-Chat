# 💬 Go Gin Realtime Chat

Proyek ini adalah aplikasi **Realtime Chat** berbasis **Go (Gin Framework)** dengan dukungan:

- 🔑 **JWT Authentication**  
- 🗄 **PostgreSQL + GORM** untuk manajemen user, role, permission, dan pesan  
- ⚡ **Redis (Pub/Sub)** untuk distribusi pesan lintas instance  
- 🌐 **WebSocket** untuk komunikasi real-time  
- 📑 **Swagger Docs** untuk dokumentasi API  
- 🛡 **Graceful shutdown** & signal handling  

---

## 🚀 Fitur Utama

- **Autentikasi dengan JWT** (login menghasilkan token, digunakan untuk akses API & WebSocket)  
- **Manajemen User, Role, Permission** via API  
- **Chat Real-time** dengan WebSocket  
- **Redis Pub/Sub** untuk broadcasting pesan ke semua client di berbagai server  
- **Swagger UI** untuk eksplorasi API di `http://localhost:8082/swagger/index.html`  

---

## 📂 Struktur Project

```bash
Go-Gin-Realtime-Chat/
├── database/         # Koneksi DB
├── docs/             # Swagger docs
├── models/           # User, Role, Permission, Message, dll
├── routes/           # Routing untuk Auth, User, Role, Permission
├── main.go           # Entry point
├── hub.go            # Hub untuk WebSocket
├── client.go         # Client handler
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

---

## ⚙️ Installation

### 1. Clone Repo
```bash
git clone https://github.com/AhmadFazriNursamsi/Go-Gin-Realtime-Chat.git
cd Go-Gin-Realtime-Chat
```

### 2. Setup `.env`
Buat file `.env` berdasarkan `.env.example`:
```env
APP_PORT=8082
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=chatdb
JWT_SECRET=supersecret
REDIS_ADDR=localhost:6379
```

### 3. Jalankan dengan Docker Compose
```bash
docker-compose up --build
```

---

## 🔌 API Endpoints

### Auth
- `POST /auth/login` → Login user, return JWT  
- `POST /auth/register` → Register user  

### User
- `GET /users` → List users  
- `GET /users/me` → Info user dari token  

### Role & Permission
- `POST /roles` → Tambah role  
- `POST /permissions` → Tambah permission  
- `POST /role-permissions` → Assign permission ke role  

### WebSocket
- `GET /ws?token=<JWT>` → Connect WebSocket untuk chat  

---

## 📡 Contoh WebSocket

Connect dengan query param token JWT:

```javascript
let ws = new WebSocket("ws://localhost:8082/ws?token=Bearer <your_jwt>");
ws.onmessage = (msg) => console.log("📩", msg.data);
ws.send(JSON.stringify({ room_id: 1, content: "Halo semua!" }));
```

---

## 📑 Swagger Docs

Setelah server berjalan, buka:  
👉 [http://localhost:8082/swagger/index.html](http://localhost:8082/swagger/index.html)  

---

## 🧪 Testing

Gunakan **Postman** atau **cURL**:

```bash
# Register user
curl -X POST http://localhost:8082/auth/register   -H "Content-Type: application/json"   -d '{"name":"Fazri","email":"fazri@example.com","password":"123456"}'

# Login
curl -X POST http://localhost:8082/auth/login   -H "Content-Type: application/json"   -d '{"email":"fazri@example.com","password":"123456"}'
```

Gunakan JWT yang didapat untuk connect ke WebSocket.

---

## 📜 License

Proyek ini dirilis di bawah lisensi [MIT License](LICENSE).  
