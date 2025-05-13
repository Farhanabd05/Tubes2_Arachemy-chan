# 🧪 Arachemy-chan: Little Alchemy Recipe Solver

A web application that implements BFS and DFS algorithms to find crafting paths in Little Alchemy 2. Built with React frontend and Go backend.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://golang.org/)
[![React Version](https://img.shields.io/badge/React-18+-61DAFB?logo=react)](https://reactjs.org/)
[![License](https://img.shields.io/badge/License-MIT-blue)](LICENSE)

## 🌐 Explore the Application

Check out the live demo: [Arachemy-chan Web App](https://frontend-production-c72f.up.railway.app/)

## ✨ Features

- 🕵️‍♀️ **Dual algorithm support** (BFS & DFS)
- 🚀 **Parallel processing** with worker pools
- 📊 **Performance metrics** (runtime & nodes visited)
- 🎨 **Interactive web interface**
- 🔍 **Multiple path finding** capabilities

## 🛠️ Technologies

### Frontend
- React + TypeScript
- Vite
- Tailwind CSS

### Backend
- Go 1.24.2
- Gin Web Framework
- Colly (Web Scraper)

## 📦 Installation

### Prerequisites
- WSL 2 (Windows) or Linux environment
- Go 1.22+
- Node.js 18+
- npm 9+
- Docker v2 (optional)

### Setup Instructions

1. **Clone Repository**
   ```bash
   git clone https://github.com/Farhanabd05/Tubes2_Arachemy-chan.git
   cd Tubes2_Arachemy-chan
   ```

2. **Choose Your Setup Method**

   ### Option A: Using Docker

   **If you need to install Docker first (Linux/Ubuntu):**
   ```bash
   # Update package list and install dependencies
   sudo apt-get update
   sudo apt-get install ca-certificates curl gnupg

   # Add Docker''s official GPG key
   sudo install -m 0755 -d /etc/apt/keyrings
   curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
   sudo chmod a+r /etc/apt/keyrings/docker.gpg

   # Add Docker repository to APT sources
   echo \
     "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
     $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
     sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

   # Update package list again
   sudo apt-get update

   # Install Docker Compose Plugin
   sudo apt-get install docker-compose-plugin

   # Verify installation
   sudo docker compose version
   ```

   **Running with Docker:**
   ```bash
   # Stop any running containers
   sudo docker compose down
   
   # Build and start the application
   sudo docker compose up --build
   ```

   ### Option B: Manual Setup

   **Frontend Setup:**
   ```bash
   cd frontend
   npm install
   ```

   **Backend Setup:**
   ```bash
   cd backend
   go mod tidy
   ```

3. **Environment Setup (For Manual Installation)**

   **For Windows Users:**
   ```bash
   wsl --install
   wsl --set-default-version 2
   ```

   **Install Development Tools:**
   ```bash
   sudo apt update && sudo apt upgrade -y
   sudo apt install curl snapd -y
   ```

   **Install Node.js (using nvm):**
   ```bash
   curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash
   source ~/.bashrc
   nvm install --lts
   nvm use --lts
   ```

   **Install Go:**
   ```bash
   sudo snap install go --classic
   ```
   
## 🚀 Running the Application

### Using Docker
The application will automatically start after running `sudo docker compose up --build`.
if you already build, just use this command: `sudo docker compose up`
### Manual Start

**Start Backend:**
```bash
cd backend
go build -o main
./main  # or just: go run main
```

**Start Frontend:**
```bash
cd frontend
npm run dev
```

The application will be available at:
- Frontend: `http://localhost:5173`
- Backend: `http://localhost:8080`
Note: if you use wsl, to run docker, you must click link http://localhost:5173 in:
frontend-1  |  INFO  Accepting connections at http://localhost:5173

## 🖥️ Usage

1. Enter target element (e.g., "gold")
2. Select algorithm (BFS or DFS)
3. Choose number of paths to find
4. Click "Cari" to start search

## 🧠 Algorithm Implementation

### BFS (Breadth-First Search)
```go
func bfsSinglePath(target string) ([]string, bool, time.Duration, int) {
    // Implementation details...
}
```

### DFS (Depth-First Search)
```go
func dfsSinglePath(element string, visited map[string]bool, trace []string, nodesVisited *int) ([]string, bool) {
    // Implementation details...
}
```

## 📂 Project Structure

```
├── backend
│   ├── main.go          # API server entrypoint
│   ├── bfsMultiple.go   # Parallel BFS implementation
│   ├── bfsSingle.go     # Single BFS implementation
│   ├── dfsMultiple.go   # Parallel DFS implementation
│   ├── dfsSingle.go     # Single DFS implementation
│   ├── scrape.go        # Scrape implementation
│   ├── utils.go         # Data loading utilities
│   └── data
│       └── recipes.json # Element combinations database
└── frontend
    ├── src
    │   ├── components   # React components
    │   └── App.tsx      # Main application logic
    └── package.json
```

## 🤝 Contributing

1. Fork the project
2. Create your feature branch:
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. Commit your changes:
   ```bash
   git commit -m 'Add some amazing feature'
   ```
4. Push to the branch:
   ```bash
   git push origin feature/amazing-feature
   ```
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Contributors

- [Abdullah Farhan](https://github.com/Farhanabd05)
- [Rafizan MZ](https://github.com/Rafizan46)
- [Muhammad Zahran](https://github.com/Muzaraar22)

---

Made with ❤️ by Arachemy-chan Team