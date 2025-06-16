# Todo Application

A full-stack Todo application built with React (TypeScript) and Go.

## Project Structure

- `frontend/` - React TypeScript application
- `main.go` - Go backend server
- `go.mod` - Go module file

## Prerequisites

- Node.js and npm
- Go 1.21 or later

## Setup and Running

### Backend

1. Install Go dependencies:
```bash
go mod tidy
```

2. Run the Go server:
```bash
go run main.go
```
The server will start on http://localhost:8080

### Frontend

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm start
```
The frontend will start on http://localhost:3000

## Features

- Create new todos
- Mark todos as complete/incomplete
- Delete todos
- Persistent storage (in-memory for now)
- Clean and responsive UI 