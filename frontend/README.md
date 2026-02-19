# GoWeather - SvelteKit Frontend

A beautiful SvelteKit frontend for the GoWeather application.

## Setup & Running

### 1. Start the Go Backend (Terminal 1)

```bash
go run main.go
```

The backend will run on `http://localhost:8080`

### 2. Start the SvelteKit Frontend (Terminal 2)

```bash
cd frontend
npm run dev
```

The frontend will run on `http://localhost:5173`

## Usage

1. Open your browser to `http://localhost:5173`
2. Enter a US address (e.g., "1600 Pennsylvania Ave NW, Washington, DC")
3. Click "Get Weather" to see the forecast

## Features

- ✨ Modern, responsive UI with gradient background
- 🎨 Beautiful card-based design
- 🔄 Loading states
- ❌ Error handling
- 📱 Mobile-friendly
- 🌐 Proxy configuration to communicate with Go backend

## Tech Stack

- **Frontend**: SvelteKit 2 with Svelte 5
- **Backend**: Go with Gin framework
- **APIs**: OpenStreetMap Nominatim & National Weather Service API
