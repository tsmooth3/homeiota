package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type AlertPreference struct {
	GotifyToken      sql.NullString `db:"gotifyToken"`
	UserId           string         `db:"userId"`
	Location         string         `db:"location"`
	Threshold        float64        `db:"threshold"`
	Enabled          bool           `db:"enabled"`
	OfflineThreshold sql.NullFloat64 `db:"offlineThreshold"`
}

type PumpRunTime struct {
	Current   float64   `db:"current"`
	Timestamp time.Time `db:"timestamp"`
}

type Temperature struct {
	Value     float64   `db:"value"`
	Timestamp time.Time `db:"timestamp"`
}

type DeviceHeartbeat struct {
	Timestamp time.Time `db:"timestamp"`
}

func sendGotifyAlert(token, title, message string, priority int, logShort ...string) {
	if token == "" {
		log.Printf("No Gotify token provided, skipping alert: %s - %s", title, message)
		return
	}
	gotifyURL := os.Getenv("GOTIFY_URL")
	
	url := fmt.Sprintf("%s/message?token=%s", gotifyURL, token)
	payload := map[string]interface{}{
		"title":    title,
		"message":  message,
		"priority": priority,
	}
	jsonPayload, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("Failed to send Gotify alert: %v", err)
		return
	}
	defer resp.Body.Close()
	if len(logShort) > 0 && logShort[0] != "" {
		log.Printf("%s", logShort[0])
	} else {
		log.Printf("Sent Gotify alert: %s - %s (status %d)", title, message, resp.StatusCode)
	}
}

func main() {

	log.Printf("Go alert script triggered at %s", time.Now().Format(time.RFC3339))

	HOMEIOTA_URL := os.Getenv("HOMEIOTA_URL")
	GOHOME_DB_URL := os.Getenv("GOHOME_DB_URL")
	HOMEIOTA_DB_URL := os.Getenv("HOMEIOTA_DB_URL")
	gohomeDBConn, err := sqlx.Open("postgres", GOHOME_DB_URL)
	if err != nil {
		log.Fatalf("Failed to connect to gohome db: %v", err)
	}
	defer gohomeDBConn.Close()

	homeiotaDBConn, err := sqlx.Open("postgres", HOMEIOTA_DB_URL)
	if err != nil {
		log.Fatalf("Failed to connect to homeiota db: %v", err)
	}
	defer homeiotaDBConn.Close()

	alertPreferences := []AlertPreference{}
	alertPrefQuery := `select "User"."gotifyToken","AlertPreference".* from "User" join "AlertPreference" on "AlertPreference"."userId" = "User".id`
	err = homeiotaDBConn.Select(&alertPreferences, alertPrefQuery)
	if err != nil {
		log.Fatalf("Failed to fetch alert preferences: %v", err)
	}

	now := time.Now().UTC()
	timeDeltaAgo := now.Add(-60 * time.Minute)

	pumpResults := make(map[string][]PumpRunTime)
	tempResults := make(map[string][]Temperature)
	heartbeatResults := make(map[string]bool)
	gotifyTokens := make(map[string]string)
	enabledMap := make(map[string]bool)
	thresholdMap := make(map[string]float64)
	offlineThresholdMap := make(map[string]float64)

	for _, pref := range alertPreferences {
		location := pref.Location
		gotifyToken := ""
		if pref.GotifyToken.Valid {
			gotifyToken = pref.GotifyToken.String
		}
		gotifyTokens[location] = gotifyToken
		enabledMap[location] = pref.Enabled
		thresholdMap[location] = pref.Threshold
		if pref.OfflineThreshold.Valid {
			offlineThresholdMap[location] = pref.OfflineThreshold.Float64
		}

		if location == "wellpump" {
			pumpRows := []PumpRunTime{}
			pumpQuery := `SELECT current, timestamp
				FROM (
				  SELECT
					current,
					timestamp,
					LAG(current) OVER (ORDER BY timestamp) AS prev_current,
					LAG(timestamp) OVER (ORDER BY timestamp) AS prev_timestamp
				  FROM pump_run_times
				  WHERE timestamp > $2
				) t
				WHERE
				  current > 1 AND current < $1
				  AND prev_current > 1 AND prev_current < $1
				  AND current <> prev_current
				  AND timestamp <> prev_timestamp`
			err := gohomeDBConn.Select(&pumpRows, pumpQuery, pref.Threshold, timeDeltaAgo)
			if err != nil {
				log.Printf("Pump query error: %v", err)
			}
			pumpResults[location] = pumpRows

			if val, ok := offlineThresholdMap[location]; ok {
				heartbeatRows := []DeviceHeartbeat{}
				heartbeatTimeAgo := now.Add(-time.Duration(val) * time.Minute)
				heartbeatQuery := `SELECT timestamp FROM device_heartbeats WHERE pump = true AND timestamp > $1 LIMIT 1`
				err := gohomeDBConn.Select(&heartbeatRows, heartbeatQuery, heartbeatTimeAgo)
				if err != nil {
					log.Printf("Heartbeat query error: %v", err)
				}
				heartbeatResults[location] = len(heartbeatRows) > 0
			} else {
				heartbeatResults[location] = false
			}
		} else {
			tempRows := []Temperature{}
			tempQuery := `SELECT value, timestamp FROM temperatures WHERE location = $1 AND value > $2 AND timestamp > $3`
			err := gohomeDBConn.Select(&tempRows, tempQuery, location, pref.Threshold, timeDeltaAgo)
			if err != nil {
				log.Printf("Temp query error: %v", err)
			}
			tempResults[location] = tempRows

			if val, ok := offlineThresholdMap[location]; ok {
				offlineRows := []Temperature{}
				offlineTimeAgo := now.Add(-time.Duration(val) * time.Minute)
				offlineQuery := `SELECT value, timestamp FROM temperatures WHERE location = $1 AND timestamp > $2 LIMIT 1`
				err := gohomeDBConn.Select(&offlineRows, offlineQuery, location, offlineTimeAgo)
				if err != nil {
					log.Printf("Offline temp query error: %v", err)
				}
				heartbeatResults[location] = len(offlineRows) > 0
			} else {
				heartbeatResults[location] = false
			}
		}
	}

	// Send Gotify alert for each heartbeat result with no rows
	for location, hasHeartbeat := range heartbeatResults {
		if !hasHeartbeat && enabledMap[location] {
			token := gotifyTokens[location]
			title := fmt.Sprintf("Device Offline: %s", location)
			message := fmt.Sprintf("No heartbeat/reading for '%s' in the last offline threshold window. Device may be offline.\n\nView details: %s", location, HOMEIOTA_URL)
			shortLog := fmt.Sprintf("%s Sent Gotify alert: Device Offline: %s.", time.Now().Format(time.RFC3339), location)
			sendGotifyAlert(token, title, message, 7, shortLog)
		}
	}

	// Send Gotify alert for each temperature result with rows
	for location, rows := range tempResults {
		if len(rows) > 0 && enabledMap[location] {
			token := gotifyTokens[location]
			title := fmt.Sprintf("Temperature Alert: %s", location)
			message := fmt.Sprintf("Temperature for '%s' is above the threshold of %.2f°F at %s.\n\nView details: %s", location, rows[0].Value, rows[0].Timestamp.Format(time.RFC3339), HOMEIOTA_URL)
			shortLog := fmt.Sprintf("%s Sent Gotify alert: Temperature Alert: %s: %.2f.", time.Now().Format(time.RFC3339), location, rows[0].Value)
			sendGotifyAlert(token, title, message, 10, shortLog)
		}
	}

	// Send Gotify alert for each pump result with rows
	for location, rows := range pumpResults {
		if len(rows) > 0 && enabledMap[location] {
			token := gotifyTokens[location]
			title := fmt.Sprintf("Pump Alert: %s", location)
			message := fmt.Sprintf("Well may be low or dry. '%s' is running at %.2f Amps at %s.\n\nView details: %s", location, rows[0].Current, rows[0].Timestamp.Format(time.RFC3339), HOMEIOTA_URL)
			shortLog := fmt.Sprintf("%s Sent Gotify alert: Pump Alert: %s: %.2f.", time.Now().Format(time.RFC3339), location, rows[0].Current)
			sendGotifyAlert(token, title, message, 7, shortLog)
		}
	}

	log.Printf("Go alert script completed at %s", time.Now().Format(time.RFC3339))
} 