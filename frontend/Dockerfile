# Stage 1: build
FROM node:23-alpine AS builder
WORKDIR /app

# Terima build‐arg VITE_BACKEND_URL
ARG VITE_BACKEND_URL
# Set sebagai ENV agar vite bisa membacanya saat `npm run build`
ENV VITE_BACKEND_URL=$VITE_BACKEND_URL

COPY package*.json ./
RUN npm ci
COPY . .
# Build aplikasi—VITE_BACKEND_URL sudah tersedia di import.meta.env
RUN npm run build

# Stage 2: runtime
FROM node:23-alpine
WORKDIR /app
RUN npm install -g serve

# Salin hasil build
COPY --from=builder /app/dist ./dist

EXPOSE 8080
CMD ["sh", "-c", "serve -s dist -l $PORT"]
