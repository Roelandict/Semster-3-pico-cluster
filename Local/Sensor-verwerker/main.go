package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	PostgrestBase        string
	PostgrestAPI         string
	JWTSecret            string
	TruckVIN             string
	TruckID              int
	SensorCount          int
	tryMinutesUntilPanic int
)

func init() {
	PostgrestBase = getEnv("POSTGREST_BASE_URL", "http://postgrest-service.foodchain-db.svc.cluster.local:3000")
	PostgrestAPI = PostgrestBase + "/currenttemperature"
	JWTSecret = getEnv("JWT_SECRET", "super-secret-jwt-key-conform-plan-van-aanpak")
	TruckVIN = getEnv("TRUCK_VIN", "FC-TRUCK-2026-X99")

	var err error
	TruckID, err = strconv.Atoi(getEnv("TRUCK_ID", "42"))
	if err != nil {
		fmt.Printf("Warning: Invalid TRUCK_ID, using default 42. Error: %v\n", err)
		TruckID = 42
	}

	SensorCount, err = strconv.Atoi(getEnv("SENSOR_COUNT", "30"))
	if err != nil {
		fmt.Printf("Warning: Invalid SENSOR_COUNT, using default 30. Error: %v\n", err)
		SensorCount = 30
	}

	tryMinutesUntilPanic, err = strconv.Atoi(getEnv("TRY_MINUTES_UNTIL_PANIC", "10"))
	if err != nil {
		fmt.Printf("Warning: Invalid TRY_MINUTES_UNTIL_PANIC, using default 10. Error: %v\n", err)
		tryMinutesUntilPanic = 10
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type Sensor struct {
	ID    string
	Zone  string
	Value float64
}

type TemperaturePayload struct {
	SensorID       string  `json:"sensor_id"`
	TemperatureAvg float64 `json:"temperature_avg"`
	Units          string  `json:"units"`
	TruckID        int     `json:"truck_id"`
}

// JWTClaims represents the JWT payload
type JWTClaims struct {
	Role      string `json:"role"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

func main() {
	fmt.Printf("Pico-cluster Edge Processor gestart voor Truck %s\n", TruckVIN)
	rand.Seed(time.Now().UnixNano())

	// Test PostgREST connection first
	fmt.Println("Testing PostgREST connectivity...")

	// Probeer tot 10 minuten lang verbinding te maken
	deadline := time.Now().Add(time.Duration(tryMinutesUntilPanic) * time.Minute)

	for time.Now().Before(deadline) {
		if testPostgRESTConnection() {
			fmt.Println("PostgREST verbonden, start data verzending...")

			ticker := time.NewTicker(1 * time.Minute)
			defer ticker.Stop()

			// Stuur onmiddellijk het eerste bericht
			sendSensorData()

			for range ticker.C {
				sendSensorData()
			}
			return
		}

		// Nog geen verbinding, wacht 1 minuut voor volgende poging
		fmt.Printf("Opnieuw proberen in 1 minuut...\n")
		time.Sleep(1 * time.Minute)
	}

	panic("FATAL: PostgREST is niet bereikbaar na 10 minuten!")
}

// testPostgRESTConnection checks if PostgREST is reachable
func testPostgRESTConnection() bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Test basic connectivity
	resp, err := client.Get(PostgrestBase)
	if err != nil {
		fmt.Printf("[ERROR] Cannot reach PostgREST at %s: %v\n", PostgrestBase, err)
		return false
	}
	defer resp.Body.Close()

	fmt.Printf("[OK] PostgREST is reachable at %s (HTTP %d)\n", PostgrestBase, resp.StatusCode)

	// Try to list available tables/endpoints
	resp, err = client.Get(PostgrestBase + "/")
	if err == nil {
		defer resp.Body.Close()
		fmt.Printf("[OK] PostgREST OpenAPI available (HTTP %d)\n", resp.StatusCode)
	}

	return true
}

func sendSensorData() {
	rawSensors := simulateZoneSensors(SensorCount)
	avgTemp := calculateAverage(rawSensors)

	payload := TemperaturePayload{
		SensorID:       fmt.Sprintf("AGGR-%s", TruckVIN),
		TemperatureAvg: avgTemp,
		Units:          "Celsius",
		TruckID:        TruckID,
	}

	sendToPostgREST(payload)
}

func simulateZoneSensors(count int) []Sensor {
	var sensors []Sensor
	for i := 1; i <= count; i++ {
		var zone string
		var temp float64

		if i <= 10 {
			zone = "Front"
			temp = -18.5 + (rand.Float64() * 0.5)
		} else if i <= 20 {
			zone = "Middle"
			temp = -18.0 + (rand.Float64() * 0.4)
		} else {
			zone = "Back"
			temp = -17.8 + (rand.Float64()*0.8 + 6.7) // + 1 only for testing = test change for ayoub
		}

		sensors = append(sensors, Sensor{
			ID:    fmt.Sprintf("S-%03d", i),
			Zone:  zone,
			Value: temp,
		})
	}
	return sensors
}

func calculateAverage(sensors []Sensor) float64 {
	var total float64
	for _, s := range sensors {
		total += s.Value
	}
	return total / float64(len(sensors))
}

// generateJWT creates a JWT token signed with HMAC-SHA256
func generateJWT() (string, error) {
	now := time.Now()
	claims := JWTClaims{
		Role:      "sensor_admin",
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(1 * time.Hour).Unix(),
	}

	// Create header
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	// Encode header and claims to base64
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsEncoded := base64.RawURLEncoding.EncodeToString(claimsJSON)

	// Create signature
	message := headerEncoded + "." + claimsEncoded
	h := hmac.New(sha256.New, []byte(JWTSecret))
	h.Write([]byte(message))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return message + "." + signature, nil
}

func sendToPostgREST(p TemperaturePayload) {
	jsonData, err := json.Marshal(p)
	if err != nil {
		fmt.Printf("[ERROR] JSON encoding failed: %v\n", err)
		return
	}

	// Generate JWT token
	token, err := generateJWT()
	if err != nil {
		fmt.Printf("[ERROR] JWT generation failed: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", PostgrestAPI, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("[ERROR] Request creation failed: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("[ERROR] Connection failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read response body for better error messages
	bodyBytes, err := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	// Handle different status codes
	switch {
	case resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK:
		fmt.Printf("[%s] SUCCESS: %d sensoren geaggregeerd. Gemiddelde: %.2f°C verzonden naar DB.\n",
			time.Now().Format("15:04:05"), SensorCount, p.TemperatureAvg)

	case resp.StatusCode == http.StatusNotFound:
		fmt.Printf("[%s] ERROR 404: PostgREST endpoint niet gevonden!\n", time.Now().Format("15:04:05"))
		fmt.Printf("  API URL: %s\n", PostgrestAPI)
		fmt.Printf("  Response: %s\n", bodyStr)
		fmt.Println("  OPOSSING: Check dat de 'currentTemperature' tabel bestaat in PostgREST schema")

	case resp.StatusCode == http.StatusUnauthorized:
		fmt.Printf("[%s] ERROR 401: JWT Authentication failed\n", time.Now().Format("15:04:05"))
		fmt.Printf("  Response: %s\n", bodyStr)

	case resp.StatusCode == http.StatusForbidden:
		fmt.Printf("[%s] ERROR 403: Permission denied\n", time.Now().Format("15:04:05"))
		fmt.Printf("  Response: %s\n", bodyStr)

	case resp.StatusCode == http.StatusBadRequest:
		fmt.Printf("[%s] ERROR 400: Bad request\n", time.Now().Format("15:04:05"))
		fmt.Printf("  Payload: %s\n", string(jsonData))
		fmt.Printf("  Response: %s\n", bodyStr)

	default:
		fmt.Printf("[%s] HTTP %d: %s\n", time.Now().Format("15:04:05"), resp.StatusCode, bodyStr)
	}
}
