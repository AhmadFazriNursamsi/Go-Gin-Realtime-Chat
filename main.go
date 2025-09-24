package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"myapi/database"
	"myapi/docs"
	"myapi/models"
	"myapi/routes"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// TODO: ganti sesuai kebutuhan (hanya izinkan origin tertentu di production)
			return true
		},
	}
	// rdb       *redis.Client
	// ctx       = context.Background()
	jwtSecret         []byte // ✅ ini jadi source of truth
	subscribedRooms   = make(map[uint]bool)
	subscribedRoomsMu sync.Mutex
)

type Claims struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var ctx = context.Background()
var rdb *redis.Client

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"), // "localhost:6379"
	})
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("❌ Redis connection failed:", err)
	}
	log.Println("✅ Redis connected:", pong)

}

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Masukkan token JWT dalam format 'Bearer <token>'
func main() {
	// 1️⃣ Connect database & migrate
	database.Connect()
	models.Migrate(database.DB)

	// 2️⃣ Init Redis
	initRedis()

	// Jalankan subscriber di background
	go subscribeMessages("chatroom-1")

	// Tidak perlu publish pesan di startup kecuali untuk test
	// publishMessage("chatroom-1", `{"sender_id":1,"content":"Halo semua!"}`)

	// 3️⃣ Ambil JWT secret dari env (atau fallback)
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		log.Println("⚠️  JWT_SECRET tidak ditemukan di env, gunakan default secret.")
		jwtSecret = []byte("0WvRtY6h9V7qCrMm6KDxjD3c6nQFlQ0gTK9r4ggh7LM=")
	}
	log.Printf("🔑 JWT secret loaded, length=%d\n", len(jwtSecret))

	// 4️⃣ Init Gin
	r := gin.Default()

	// Swagger setup
	docs.SwaggerInfo.Title = "My API"
	docs.SwaggerInfo.Description = "API documentation with Swagger"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8082"
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ✅ Register routes
	routes.RegisterAuthRoutes(r, database.DB, jwtSecret)
	routes.RegisterUserRoutes(r, database.DB)
	routes.RegisterRoleRoutes(r, database.DB)
	routes.RegisterRoleChildRoutes(r, database.DB)
	routes.RegisterPermissionRoutes(r, database.DB)
	routes.RegisterRolePermissionRoutes(r, database.DB)

	// WebSocket route
	r.GET("/ws", func(c *gin.Context) {
		wsHandler(c.Writer, c.Request)
	})

	// Graceful shutdown support
	srv := &http.Server{
		Addr:    ":8082",
		Handler: r,
	}

	go func() {
		log.Println("🚀 Server running on port 8082")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %s", err)
	}

	log.Println("✅ Server exited cleanly")
}

func ensureSubscribedRoom(roomID uint) {
	subscribedRoomsMu.Lock()
	defer subscribedRoomsMu.Unlock()
	if subscribedRooms[roomID] {
		return
	}
	channel := fmt.Sprintf("chatroom-%d", roomID)
	go subscribeMessages(channel) // fungsi kamu yang sudah ada: subscribeMessages(channel)
	subscribedRooms[roomID] = true
	log.Printf("🔔 Subscribed to Redis channel '%s'", channel)
}

// type Client struct {
// 	ID    uint
// 	Rooms []uint // support multi-room
// 	Conn  *websocket.Conn
// 	Send  chan []byte
// }

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("🌐 Incoming WebSocket connection...")

	// 🔑 Ambil token dari query
	rawToken := r.URL.Query().Get("token")
	if rawToken == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	token := strings.TrimSpace(rawToken)
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		token = token[7:]
	}

	// 🔍 Parse & validasi token
	claims, err := validateAndParseClaims(token)
	if err != nil {
		log.Println("❌ invalid token:", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID := claims.ID
	roomIDs := claims.RoomsId // ✅ langsung ambil dari JWT

	// ✅ Jika user belum punya room, buatkan room baru
	if len(roomIDs) == 0 {
		newRoom := models.Rooms{Name: fmt.Sprintf("Room-%d", userID)}
		if err := database.DB.Create(&newRoom).Error; err != nil {
			log.Printf("❌ gagal membuat room baru: %v", err)
			http.Error(w, "gagal membuat room baru", http.StatusInternalServerError)
			return
		}

		roomIDs = append(roomIDs, newRoom.ID)
		log.Printf("🆕 Room baru dibuat: ID=%d oleh user %d", newRoom.ID, userID)
	}

	// Upgrade ke WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ upgrade error:", err)
		return
	}

	client := &Client{
		ID:    userID,
		Rooms: roomIDs,
		Conn:  conn,
		Send:  make(chan []byte, 256),
	}

	hub.register(client)
	go writePump(client)

	// --- Ambil history untuk setiap room ---
	type MessageWithName struct {
		ID         uint      `json:"id"`
		RoomID     uint      `json:"room_id"`
		SenderID   uint      `json:"sender_id"`
		Content    string    `json:"content"`
		Type       string    `json:"type"`
		CreatedAt  time.Time `json:"created_at"`
		SenderName string    `json:"sender_name"`
	}

	for _, roomID := range roomIDs {
		var history []MessageWithName
		err := database.DB.Table("messages m").
			Select("m.id, m.room_id, m.sender_id, m.content, m.type, m.created_at, u.name AS sender_name").
			Joins("LEFT JOIN users u ON u.id = m.sender_id").
			Where("m.room_id = ?", roomID).
			Order("m.created_at ASC").
			Scan(&history).Error

		if err != nil {
			log.Printf("❌ DB error ambil history room %d: %v", roomID, err)
			continue
		}

		log.Printf("📜 Kirim %d pesan history untuk room %d ke user %d", len(history), roomID, userID)
		for _, m := range history {
			payload, _ := json.Marshal(m)
			client.Send <- payload
		}
	}

	readPump(client)
}
func validateAndParseClaims(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}

func publishMessage(channel string, message string) {
	err := rdb.Publish(ctx, channel, message).Err()
	if err != nil {
		log.Println("❌ Failed to publish message:", err)
	}
}
func subscribeMessages(channel string) {
	pubsub := rdb.PSubscribe(ctx, "chatroom-*")
	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			log.Printf("📩 Received via Redis channel '%s': %s", msg.Channel, msg.Payload)
			hub.broadcast([]byte(msg.Payload))
		}
	}()
}

func validateJWT(tokenString string) (string, error) {
	log.Printf("📥 Validating token (len=%d)", len(tokenString))
	log.Printf("🔑 Using JWT secret (len=%d)", len(jwtSecret))

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		log.Printf("🔍 Token header: %+v", token.Header)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		log.Println("❌ JWT parse error:", err)
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", fmt.Errorf("token expired")
		}
		return "", fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("token not valid")
	}

	log.Printf("✅ Claims parsed: %+v", claims)

	idVal, ok := claims["id"]
	if !ok {
		return "", fmt.Errorf("id claim missing")
	}

	switch v := idVal.(type) {
	case float64:
		return strconv.Itoa(int(v)), nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("unsupported id type: %T", v)
	}
}
