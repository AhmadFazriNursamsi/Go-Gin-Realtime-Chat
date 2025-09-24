# Gunakan image Go sebagai base
FROM golang:alpine

# Set working directory
WORKDIR /app

# Copy go.mod dan go.sum dulu biar cache build cepat
COPY go.mod go.sum ./
RUN go mod download

# Copy semua file source code
COPY . .

# Build aplikasi
RUN go build -o app .

# Expose port sesuai app (misal 8082)
EXPOSE 8082

# Jalankan aplikasi
CMD ["./app"]
