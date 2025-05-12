# Gunakan image Go terbaru dengan Alpine
FROM golang:1.24-alpine

# Install git dan curl (dibutuhkan oleh air)
RUN echo "https://dl-cdn.alpinelinux.org/alpine/v3.21/main" > /etc/apk/repositories && \
    echo "https://dl-cdn.alpinelinux.org/alpine/v3.21/community" >> /etc/apk/repositories

RUN apk add --no-cache git curl

# Install air (hot reload tool)
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b /usr/local/bin

# Buat direktori kerja
WORKDIR /app

# Salin semua file ke dalam container
COPY . .

# Unduh dependencies
RUN go mod tidy

# Buka port 8080
EXPOSE 8080

# Jalankan air untuk hot reload
CMD ["air"]
