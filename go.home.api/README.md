# Go Home API

This project is a Go service that monitors temperature readings, pump run times, and device heartbeats. It uses a PostgreSQL database to store the data and provides a RESTful API for interaction.

## Project Structure

- `main.go`: Contains the main application logic, data structures, database connection, and HTTP handlers.
- `Dockerfile`: Used to build a Docker image for the Go service.
- `README.md`: Documentation for the project.

## Getting Started

### Prerequisites

- Go (1.16 or later)
- Docker

### Building the Docker Image

1. Navigate to the project directory:

   ```
   cd go.home.api
   ```

2. Build the Docker image:

   ```
   docker build -t go-home-api .
   ```

### Running the Docker Container

1. Run the Docker container:

   ```
   docker run -p 8080:8080 --env dbhost=<your_db_host> --env dbport=<your_db_port> --env dbuser=<your_db_user> --env dbpass=<your_db_password> --env dbname=<your_db_name> go-home-api
   ```

   Replace `<your_db_host>`, `<your_db_port>`, `<your_db_user>`, `<your_db_password>`, and `<your_db_name>` with your PostgreSQL database credentials.

### API Endpoints

- `POST /heartbeat`: Create a new device heartbeat.
- `GET /tempmon`: Retrieve all temperature readings.
- `POST /tempmon`: Create a new temperature reading.
- `GET /tempmon/{id}`: Retrieve a single temperature reading by ID.
- `DELETE /tempmon/{id}`: Delete a temperature reading by ID.
- `GET /pumpmon`: Retrieve all pump run times.
- `POST /pumpmon`: Create a new pump run time.
- `GET /pumpmon/{id}`: Retrieve a single pump run time by ID.
- `DELETE /pumpmon/{id}`: Delete a pump run time by ID.
