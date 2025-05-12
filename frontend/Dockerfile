# frontend/Dockerfile
FROM node:23-alpine

WORKDIR /app

COPY . .

RUN npm install

RUN npm run build

# Gunakan serve untuk menyajikan hasil build secara statis
RUN npm install -g serve

EXPOSE 5173

CMD ["serve", "-s", "dist", "-l", "5173"]
