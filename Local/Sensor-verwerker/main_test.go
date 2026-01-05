package main

import (
	"encoding/json"
	"math"
	"strings"
	"testing"
	"time"
)

// TestGenerateJWT verifies JWT token generation and format
func TestGenerateJWT(t *testing.T) {
	token, err := generateJWT()
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	// JWT should have 3 parts separated by dots
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("Expected JWT with 3 parts, got %d", len(parts))
	}

	// Token should not be empty
	if token == "" {
		t.Error("JWT token should not be empty")
	}
}

// TestTemperaturePayload verifies payload marshaling
func TestTemperaturePayload(t *testing.T) {
	payload := TemperaturePayload{
		SensorID:       "TEST-SENSOR",
		TemperatureAvg: -18.5,
		Units:          "Celsius",
		TruckID:        42,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	// Verify JSON can be unmarshaled back
	var unmarshaled TemperaturePayload
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal payload: %v", err)
	}

	if unmarshaled.SensorID != payload.SensorID {
		t.Errorf("SensorID mismatch: expected %s, got %s", payload.SensorID, unmarshaled.SensorID)
	}
	if unmarshaled.TemperatureAvg != payload.TemperatureAvg {
		t.Errorf("Temperature mismatch: expected %f, got %f", payload.TemperatureAvg, unmarshaled.TemperatureAvg)
	}
	if unmarshaled.TruckID != payload.TruckID {
		t.Errorf("TruckID mismatch: expected %d, got %d", payload.TruckID, unmarshaled.TruckID)
	}
}

// TestCalculateAverage verifies temperature average calculation
func TestCalculateAverage(t *testing.T) {
	tests := []struct {
		name     string
		sensors  []Sensor
		expected float64
	}{
		{
			name:     "empty sensors",
			sensors:  []Sensor{},
			expected: 0,
		},
		{
			name: "single sensor",
			sensors: []Sensor{
				{ID: "S-001", Zone: "Front", Value: -18.5},
			},
			expected: -18.5,
		},
		{
			name: "multiple sensors",
			sensors: []Sensor{
				{ID: "S-001", Zone: "Front", Value: 10.0},
				{ID: "S-002", Zone: "Front", Value: 20.0},
				{ID: "S-003", Zone: "Front", Value: 30.0},
			},
			expected: 20.0,
		},
		{
			name: "negative values",
			sensors: []Sensor{
				{ID: "S-001", Zone: "Front", Value: -10.0},
				{ID: "S-002", Zone: "Front", Value: -20.0},
			},
			expected: -15.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateAverage(tt.sensors)

			// Handle division by zero case
			if len(tt.sensors) == 0 {
				if !math.IsNaN(result) {
					t.Errorf("Expected NaN for empty sensors, got %f", result)
				}
				return
			}

			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

// TestSimulateZoneSensors verifies sensor simulation
func TestSimulateZoneSensors(t *testing.T) {
	tests := []struct {
		name  string
		count int
	}{
		{
			name:  "10 sensors",
			count: 10,
		},
		{
			name:  "30 sensors",
			count: 30,
		},
		{
			name:  "100 sensors",
			count: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sensors := simulateZoneSensors(tt.count)

			if len(sensors) != tt.count {
				t.Errorf("Expected %d sensors, got %d", tt.count, len(sensors))
			}

			// Verify all sensors have valid data
			for i, sensor := range sensors {
				if sensor.ID == "" {
					t.Errorf("Sensor %d has empty ID", i)
				}
				if sensor.Zone == "" {
					t.Errorf("Sensor %d has empty Zone", i)
				}
				if sensor.Value > 0 {
					t.Errorf("Sensor %d temperature should be cold (negative or very low), got %f", i, sensor.Value)
				}
			}

			// Verify zone distribution
			frontCount := 0
			middleCount := 0
			backCount := 0

			for _, sensor := range sensors {
				switch sensor.Zone {
				case "Front":
					frontCount++
				case "Middle":
					middleCount++
				case "Back":
					backCount++
				}
			}

			if tt.count >= 10 && frontCount == 0 {
				t.Error("Expected Front zone sensors")
			}
			if tt.count >= 20 && middleCount == 0 {
				t.Error("Expected Middle zone sensors")
			}
			if tt.count >= 21 && backCount == 0 {
				t.Error("Expected Back zone sensors")
			}
		})
	}
}

// TestJWTClaimsStructure verifies JWT claims format
func TestJWTClaimsStructure(t *testing.T) {
	now := time.Now()
	claims := JWTClaims{
		Role:      "sensor_admin",
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(1 * time.Hour).Unix(),
	}

	if claims.Role != "sensor_admin" {
		t.Errorf("Expected role 'sensor_admin', got '%s'", claims.Role)
	}

	if claims.ExpiresAt <= claims.IssuedAt {
		t.Error("ExpiresAt should be after IssuedAt")
	}

	// Verify claims can be marshaled
	jsonData, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("Failed to marshal claims: %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("JSON marshaled claims should not be empty")
	}
}

// TestSensorStructure verifies Sensor data structure
func TestSensorStructure(t *testing.T) {
	sensor := Sensor{
		ID:    "S-001",
		Zone:  "Front",
		Value: -18.5,
	}

	if sensor.ID != "S-001" {
		t.Errorf("Expected ID 'S-001', got '%s'", sensor.ID)
	}
	if sensor.Zone != "Front" {
		t.Errorf("Expected Zone 'Front', got '%s'", sensor.Zone)
	}
	if sensor.Value != -18.5 {
		t.Errorf("Expected Value -18.5, got %f", sensor.Value)
	}
}

// TestConstants verifies configuration constants
func TestConstants(t *testing.T) {
	if TruckID <= 0 {
		t.Errorf("TruckID should be positive, got %d", TruckID)
	}

	if TruckVIN == "" {
		t.Error("TruckVIN should not be empty")
	}

	if SensorCount <= 0 {
		t.Errorf("SensorCount should be positive, got %d", SensorCount)
	}

	if PostgrestAPI == "" {
		t.Error("PostgrestAPI should not be empty")
	}

	if JWTSecret == "" {
		t.Error("JWTSecret should not be empty")
	}
}
