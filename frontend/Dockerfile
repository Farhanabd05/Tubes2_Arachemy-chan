# Gunakan image Node ringan
FROM node:23-alpine AS builder
WORKDIR /app
COPY . .
RUN npm install
RUN npm run build

# Stage runtime: hanya butuh serve
FROM node:23-alpine
WORKDIR /app
RUN npm install -g serve
COPY --from=builder /app/dist ./dist

# Biarkan Railway menentukan port
EXPOSE 8080

# Bind ke semua interface dan baca env PORT
CMD ["sh", "-c", "serve -s dist -l $PORT"]
