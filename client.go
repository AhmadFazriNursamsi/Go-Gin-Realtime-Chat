package main

import (
	"encoding/json"
	"fmt"
	"log"
	"myapi/database"
	"myapi/models"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID    uint
	Rooms []uint // sekarang slice, bukan 1 room saja
	Conn  *websocket.Conn
	Send  chan []byte
}

// readPump: read messages from client
func readPump(c *Client) {
	defer func() {
		hub.unregister(c)
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("❌ readPump error: %v", err)
			break
		}

		log.Printf("📩 Received message from User %d: %s", c.ID, string(msg))

		// 1️⃣ Unmarshal JSON
		var incoming models.Messages
		if err := json.Unmarshal(msg, &incoming); err != nil {
			log.Println("❌ JSON unmarshal error:", err)
			continue
		}

		// 2️⃣ Lengkapi data (SenderID, RoomID, Timestamp)
		if incoming.RoomID == 0 {
			if len(c.Rooms) > 0 {
				incoming.RoomID = c.Rooms[0]
			} else {
				log.Println("⚠️ Tidak ada RoomID untuk pesan ini, skip")
				continue
			}
		}

		incoming.CreatedAt = time.Now()

		// 3️⃣ Simpan ke database
		if err := database.DB.Create(&incoming).Error; err != nil {
			log.Println("❌ DB insert error:", err)
			continue
		}

		// 4️⃣ Ambil nama pengirim dari tabel users
		var senderName string
		database.DB.Table("users").Select("name").Where("id = ?", c.ID).Scan(&senderName)

		// 5️⃣ Bungkus jadi response lengkap
		type MessageWithName struct {
			ID         uint      `json:"id"`
			RoomID     uint      `json:"room_id"`
			SenderID   uint      `json:"sender_id"`
			Content    string    `json:"content"`
			Type       string    `json:"type"`
			CreatedAt  time.Time `json:"created_at"`
			SenderName string    `json:"sender_name"`
		}

		payloadStruct := MessageWithName{
			ID:         incoming.ID,
			RoomID:     incoming.RoomID,
			SenderID:   incoming.SenderID,
			Content:    incoming.Content,
			Type:       incoming.Type,
			CreatedAt:  incoming.CreatedAt,
			SenderName: senderName,
		}

		// Marshal the payload to JSON
		payload, err := json.Marshal(payloadStruct)
		if err != nil {
			log.Println("❌ JSON marshal error:", err)
			continue
		}

		// 7️⃣ Publish ke Redis
		channel := fmt.Sprintf("chatroom-%d", incoming.RoomID)
		publishMessage(channel, string(payload))

		log.Printf("📤 Published message to Redis channel '%s': %s", channel, payload)
	}
}

func writePump(c *Client) {
	ticker := time.NewTicker(25 * time.Second) // kirim ping setiap 25 detik
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
