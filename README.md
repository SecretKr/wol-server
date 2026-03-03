# Wake-on-LAN Server

This is a simple server to manage devices and send Wake-on-LAN magic packets to wake them up.
It provides a web interface to add/remove devices and wake them up.

## Features

- **Device Management**: Add and remove devices with Name and MAC address.
- **Wake-on-LAN**: Send magic packets to wake devices up.
- **Web Interface**: Clean, dark-themed UI.
- **Data Persistence**: Uses a local SQLite database (`data/devices.db`).

## Setup

1. **Install Go**: Make sure you have Go installed (1.20+ recommended).
2. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

## Running the Server

Since the project is split into multiple files (`main.go`, `db.go`, `wol.go`), you must run it using the package import path `.` or build it first.

**Do NOT run `go run main.go`**, as it will fail to find functions defined in other files.

### Option 1: Run directly
```bash
go run .
```

### Option 2: Build and Run
```bash
go build -o wol-server .
./wol-server
```

The server will start on port `8090`.
Open your browser and navigate to: [http://localhost:8090](http://localhost:8090)

## Docker Usage

### Quick Start (Docker Compose)

The easiest way to run the server is with Docker Compose:

```bash
docker compose up -d
```

This will start the server on port `8090` and persist data in the `data/` directory.

### Manual Build

Build the image:

```bash
docker build -t wol-server .
```

Run the container:

```bash
docker run -d \
  -p 8090:8090 \
  -v $(pwd)/data:/app/data \
  --name wol-server \
  wol-server
```

## Structure

- `main.go`: Main server entry point and API routes.
- `db.go`: Database operations.
- `wol.go`: Wake-on-LAN logic.
- `static/`: Frontend assets (HTML, CSS, JS).
- `data/`: Database storage.
