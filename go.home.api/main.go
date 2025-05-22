package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// TemperatureReading represents a single temperature measurement with location
type TemperatureReading struct {
	ID        int       `json:"id"`
	Value     float64   `json:"value"`     // Temperature value in Fahrenheit
	Location  string    `json:"location"`  // e.g., "freezer", "living room"
	Timestamp time.Time `json:"timestamp"` // When the reading was taken
}

// PumpRunTime represents a single water pump run time record
type PumpRunTime struct {
	ID         int       `json:"id"`
	RunTime    int       `json:"run_time"`    // Run time in seconds
	Current    float64   `json:"current"`     // current in amps
	LowCurrent bool      `json:"low_current"` // 1 if low current, 0 if not
	Timestamp  time.Time `json:"timestamp"`   // When the pump ran
}

// DeviceHeartbeat represents a device's heartbeat with device identifier
type DeviceHeartbeat struct {
	ID        int       `json:"id"`
	DeviceID  string    `json:"device_id"` // Unique identifier for the device
	Pump      bool      `json:"pump"`      // True if this is a pump device
	Timestamp time.Time `json:"timestamp"` // When the heartbeat was received
}

var db *sql.DB

func main() {
	// Initialize the database using PostgreSQL/TimescaleDB
	var err error
	// Connection string:  host, user, password, dbname, and sslmode.  Adjust as needed.
	// db_connStr="host=<host> port=5432 user=<user> password=<pw> dbname=<dbname> sslmode=disable"
	dbhost := os.Getenv("dbhost")
	dbport := os.Getenv("dbport")
	dbuser := os.Getenv("dbuser")
	dbpass := os.Getenv("dbpass")
	dbname := os.Getenv("dbname")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpass, dbname)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Create the temperatures table as a hypertable in TimescaleDB
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS temperatures (
			id SERIAL,
			value REAL NOT NULL,
			location TEXT NOT NULL,
			timestamp TIMESTAMPTZ DEFAULT NOW(),
			PRIMARY KEY (id, timestamp)
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create temperature table: %v", err)
	}

	// Check if the table is already a hypertable
	var isHypertable bool
	err = db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM timescaledb_information.hypertables WHERE hypertable_name = 'temperatures')`).Scan(&isHypertable)
	if err != nil {
		log.Fatalf("Failed to check if table is a hypertable: %v", err)
	}

	// Create hypertable only if it's not already a hypertable
	if !isHypertable {
		_, err = db.ExecContext(ctx, `SELECT create_hypertable('temperatures', 'timestamp');`)
		if err != nil {
			log.Fatalf("Failed to create hypertable: %v", err)
		}
	}

	// Create the pump_run_times table as a hypertable in TimescaleDB
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS pump_run_times (
			id SERIAL,
			run_time INTEGER NOT NULL,
			current REAL NOT NULL,
			low_current BOOLEAN NOT NULL,
			timestamp TIMESTAMPTZ DEFAULT NOW(),
			PRIMARY KEY (id, timestamp)
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create pump_run_times table: %v", err)
	}
	// Check if the table is already a hypertable
	isHypertable = false
	err = db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM timescaledb_information.hypertables WHERE hypertable_name = 'pump_run_times')`).Scan(&isHypertable)
	if err != nil {
		log.Fatalf("Failed to check if table is a hypertable: %v", err)
	}

	// Create hypertable only if it's not already a hypertable
	if !isHypertable {
		_, err = db.ExecContext(ctx, `SELECT create_hypertable('pump_run_times', 'timestamp');`)
		if err != nil {
			log.Fatalf("Failed to create hypertable: %v", err)
		}
	}

	// Create the pump_run_times_critical table as a hypertable in TimescaleDB
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS pump_run_times_critical (
			id SERIAL,
			run_time INTEGER NOT NULL,
			current REAL NOT NULL,
			low_current BOOLEAN NOT NULL,
			timestamp TIMESTAMPTZ DEFAULT NOW(),
			PRIMARY KEY (id, timestamp)
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create pump_run_times_critical table: %v", err)
	}
	// Check if the table is already a hypertable
	isHypertable = false
	err = db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM timescaledb_information.hypertables WHERE hypertable_name = 'pump_run_times_critical')`).Scan(&isHypertable)
	if err != nil {
		log.Fatalf("Failed to check if table is a hypertable: %v", err)
	}

	// Create hypertable only if it's not already a hypertable
	if !isHypertable {
		_, err = db.ExecContext(ctx, `SELECT create_hypertable('pump_run_times_critical', 'timestamp');`)
		if err != nil {
			log.Fatalf("Failed to create critical hypertable: %v", err)
		}
	}

	// Create the device_heartbeats table as a hypertable in TimescaleDB
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS device_heartbeats (
			id SERIAL,
			device_id VARCHAR(255) NOT NULL,
			pump BOOLEAN NOT NULL,
			timestamp TIMESTAMPTZ DEFAULT NOW(),
			PRIMARY KEY (id, timestamp)
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create device_heartbeats table: %v", err)
	}
	// Check if the table is already a hypertable
	isHypertable = false
	err = db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM timescaledb_information.hypertables WHERE hypertable_name = 'device_heartbeats')`).Scan(&isHypertable)
	if err != nil {
		log.Fatalf("Failed to check if table is a hypertable: %v", err)
	}

	// Create hypertable only if it's not already a hypertable
	if !isHypertable {
		_, err = db.ExecContext(ctx, `SELECT create_hypertable('device_heartbeats', 'timestamp');`)
		if err != nil {
			log.Fatalf("Failed to create device_heartbeats hypertable: %v", err)
		}
	}

	// Define API endpoints
	http.HandleFunc("/tempmon", handleTemperatures)       // GET all, POST new
	http.HandleFunc("/tempmon/", handleSingleTemperature) // GET, DELETE by ID
	http.HandleFunc("/pumpmon", handlePumpRunTimes)       // GET all, POST new
	http.HandleFunc("/pumpmon/", handleSinglePumpRunTime) // GET, DELETE by ID
	http.HandleFunc("/heartbeat", handleDeviceHeartbeats)

	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Device heartbeat handlers
func handleDeviceHeartbeats(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createDeviceHeartbeat(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createDeviceHeartbeat(w http.ResponseWriter, r *http.Request) {
	var heartbeat DeviceHeartbeat
	if err := json.NewDecoder(r.Body).Decode(&heartbeat); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if heartbeat.DeviceID == "" {
		http.Error(w, "Device ID is required", http.StatusBadRequest)
		return
	}

	// Set the timestamp to now if not provided
	if heartbeat.Timestamp.IsZero() {
		heartbeat.Timestamp = time.Now()
	}

	// Set the pump to false if not provided
	if !heartbeat.Pump {
		heartbeat.Pump = false
	}

	// Set pump to true if device_id is "pump" or "well pump"
	if heartbeat.DeviceID == "pump" {
		heartbeat.Pump = true
	}

	// Insert the new heartbeat
	query := "INSERT INTO device_heartbeats (device_id, pump, timestamp) VALUES ($1, $2, $3) RETURNING id"
	var id int
	if err := db.QueryRow(query, heartbeat.DeviceID, heartbeat.Pump, heartbeat.Timestamp).Scan(&id); err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert heartbeat: %v", err), http.StatusInternalServerError)
		return
	}

	// Set the ID on the returned object
	heartbeat.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(heartbeat); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode JSON: %v", err), http.StatusInternalServerError)
	}
}

// --- Temperature Handlers ---

func handleTemperatures(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllTemperatures(w, r)
	case http.MethodPost:
		createTemperature(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleSingleTemperature(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tempmon/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getTemperature(w, r, id)
	case http.MethodDelete:
		deleteTemperature(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllTemperatures(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	ctx := context.Background()

	// Use a prepared statement to prevent SQL injection
	query := "SELECT id, value, location, timestamp FROM temperatures"
	var args []interface{}
	if location != "" {
		query += " WHERE location = $1" // Use $1 for prepared statement arguments in PostgreSQL
		args = append(args, location)
	}

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		http.Error(w, "Failed to fetch temperatures", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer rows.Close()

	var readings []TemperatureReading
	for rows.Next() {
		var reading TemperatureReading
		if err := rows.Scan(&reading.ID, &reading.Value, &reading.Location, &reading.Timestamp); err != nil {
			http.Error(w, "Failed to scan temperature reading", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		readings = append(readings, reading)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(readings)
}

func createTemperature(w http.ResponseWriter, r *http.Request) {
	var newReading TemperatureReading
	err := json.NewDecoder(r.Body).Decode(&newReading)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	ctx := context.Background()

	// Use a prepared statement
	stmt, err := db.PrepareContext(ctx, "INSERT INTO temperatures (value, location) VALUES ($1, $2) RETURNING id") // Use $1, $2 for prepared statement arguments in PostgreSQL and  RETURNING id
	if err != nil {
		http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRowContext(ctx, newReading.Value, newReading.Location).Scan(&id) // get the ID
	if err != nil {
		http.Error(w, "Failed to create temperature reading", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	newReading.ID = id

	// Fetch the created reading to get the timestamp
	row := db.QueryRowContext(ctx, "SELECT id, value, location, timestamp FROM temperatures WHERE id = $1", id)
	err = row.Scan(&newReading.ID, &newReading.Value, &newReading.Location, &newReading.Timestamp)
	if err != nil {
		log.Println("Failed to retrieve created temperature reading:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newReading)
}

func getTemperature(w http.ResponseWriter, r *http.Request, id int) {
	ctx := context.Background()

	// Use a prepared statement
	stmt, err := db.PrepareContext(ctx, "SELECT id, value, location, timestamp FROM temperatures WHERE id = $1") // Use $1 for prepared statement arguments in PostgreSQL
	if err != nil {
		http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	var reading TemperatureReading
	err = row.Scan(&reading.ID, &reading.Value, &reading.Location, &reading.Timestamp)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch temperature reading", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reading)
}

func deleteTemperature(w http.ResponseWriter, r *http.Request, id int) {
	ctx := context.Background()

	// Use a prepared statement
	stmt, err := db.PrepareContext(ctx, "DELETE FROM temperatures WHERE id = $1") // Use $1 for prepared statement arguments in PostgreSQL
	if err != nil {
		http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		http.Error(w, "Failed to delete temperature reading", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- Pump Run Time Handlers ---

func handlePumpRunTimes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllPumpRunTimes(w, r)
	case http.MethodPost:
		createPumpRunTime(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleSinglePumpRunTime(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/pumpmon/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getPumpRunTime(w, r, id)
	case http.MethodDelete:
		deletePumpRunTime(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllPumpRunTimes(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Use a prepared statement
	stmt, err := db.PrepareContext(ctx, "SELECT id, run_time, current, low_current, timestamp FROM pump_run_times") // Corrected query
	if err != nil {
		http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch pump run times", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer rows.Close()

	var runTimes []PumpRunTime
	for rows.Next() {
		var rt PumpRunTime
		var lowCurrentFromDB bool // Change to bool
		if err := rows.Scan(&rt.ID, &rt.RunTime, &rt.Current, &lowCurrentFromDB, &rt.Timestamp); err != nil {
			http.Error(w, "Failed to scan pump run time", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		rt.LowCurrent = lowCurrentFromDB // Use the boolean value directly
		runTimes = append(runTimes, rt)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(runTimes)
}

func createPumpRunTime(w http.ResponseWriter, r *http.Request) {
	var newRunTime PumpRunTime
	err := json.NewDecoder(r.Body).Decode(&newRunTime)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert bool to int for SQLite
	lowCurrentInt := 0
	if newRunTime.LowCurrent {
		lowCurrentInt = 1
	}
	ctx := context.Background()

	// Use a prepared statement
	stmt, err := db.PrepareContext(ctx, "INSERT INTO pump_run_times (run_time, current, low_current) VALUES ($1, $2, $3) RETURNING id") // Corrected query with RETURNING
	if err != nil {
		http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRowContext(ctx, newRunTime.RunTime, newRunTime.Current, lowCurrentInt).Scan(&id)
	if err != nil {
		http.Error(w, "Failed to create pump run time", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	newRunTime.ID = id

	// If this is a critical reading (low_current = true), also insert into the critical table
	if newRunTime.LowCurrent {
		criticalStmt, err := db.PrepareContext(ctx, "INSERT INTO pump_run_times_critical (run_time, current, low_current) VALUES ($1, $2, $3)")
		if err != nil {
			log.Printf("Failed to prepare critical statement: %v", err)
		} else {
			defer criticalStmt.Close()
			_, err = criticalStmt.ExecContext(ctx, newRunTime.RunTime, newRunTime.Current, lowCurrentInt)
			if err != nil {
				log.Printf("Failed to insert into critical table: %v", err)
			}
		}
	}

	// Fetch the created run time to get the timestamp
	row := db.QueryRowContext(ctx, "SELECT id, run_time, current, low_current, timestamp FROM pump_run_times WHERE id = $1", id)
	var lowCurrentFromDB bool // changed type to bool
	err = row.Scan(&newRunTime.ID, &newRunTime.RunTime, &newRunTime.Current, &lowCurrentFromDB, &newRunTime.Timestamp)
	if err != nil {
		log.Println("Failed to retrieve created pump run time:", err)
	}
	newRunTime.LowCurrent = lowCurrentFromDB

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newRunTime)
}

func getPumpRunTime(w http.ResponseWriter, r *http.Request, id int) {
	ctx := context.Background()

	// Use a prepared statement
	stmt, err := db.PrepareContext(ctx, "SELECT id, run_time, current, low_current, timestamp FROM pump_run_times WHERE id = $1") // Corrected query
	if err != nil {
		http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	var rt PumpRunTime
	var lowCurrentFromDB bool // Change to bool
	err = row.Scan(&rt.ID, &rt.RunTime, &rt.Current, &lowCurrentFromDB, &rt.Timestamp)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch pump run time", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	rt.LowCurrent = lowCurrentFromDB

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rt)
}

func deletePumpRunTime(w http.ResponseWriter, r *http.Request, id int) {
	ctx := context.Background()
	// Use a prepared statement
	stmt, err := db.PrepareContext(ctx, "DELETE FROM pump_run_times WHERE id = $1") // Corrected query
	if err != nil {
		http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		http.Error(w, "Failed to delete pump run time", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
