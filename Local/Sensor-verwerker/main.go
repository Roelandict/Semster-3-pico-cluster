package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	PostgrestAPI = "http://192.168.1.10:30000/currentTemperature"
	TruckVIN     = "FC-TRUCK-2026-X99"
	TruckID      = 42
	SensorCount  = 30
)

type Sensor struct {
	ID    string
	Zone  string
	Value float64
}

type TemperaturePayload struct {
	SensorID       string  `json:"sensor_id"`
	TemperatureAvg float64 `json:"temperatureAvg"`
	Units          string  `json:"units"`
	TruckID        int     `json:"truckId"`
}

func main() {
	fmt.Printf("Pico-cluster Edge Processor gestart voor Truck %s\n", TruckVIN)
	rand.Seed(time.Now().UnixNano())

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// Stuur onmiddellijk het eerste bericht
	sendSensorData()

	for range ticker.C {
		sendSensorData()
	}
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
			temp = -17.8 + (rand.Float64() * 0.8)
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

func sendToPostgREST(p TemperaturePayload) {
	jsonData, err := json.Marshal(p)
	if err != nil {
		fmt.Printf("Fout bij JSON encoding: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", PostgrestAPI, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Fout bij aanmaken request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Verbindingsfout naar Hoofdserver: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		fmt.Printf("[%s] %d sensoren geaggregeerd. Gemiddelde: %.2f°C verzonden naar DB.\n",
			time.Now().Format("15:04:05"), SensorCount, p.TemperatureAvg)
	} else {
		fmt.Printf("[%s] API antwoord: %d\n", time.Now().Format("15:04:05"), resp.StatusCode)
	}
}
