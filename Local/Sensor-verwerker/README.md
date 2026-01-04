# Sensor-verwerker Edge Processor

Go-applicatie voor het aggregeren en verzenden van sensordata naar de PostgREST API via JWT-gevalideerde requests.

## Features

- Simuleert 30 temperatuursensoren over 3 zones (Front, Middle, Back)
- Aggregeert sensordata elke minuut
- JWT-ondersteuning voor veilige API communicatie
- Multi-platform ondersteuning: AMD64 (Windows/Linux) en ARM64 (Raspberry Pi)
- Kubernetes-native configuratie

## Architecturen

- **AMD64**: x86-64 processors (standaard laptops/servers)
- **ARM64**: ARM64 processors (Raspberry Pi 3B+, 4, 5)

## Building

### Windows (PowerShell)

```powershell
# Build voor AMD64 (huidige systeem)
.\build.ps1

# Build voor ARM64 (Raspberry Pi)
.\build.ps1 -TargetArch arm64

# Build voor beide architecturen
.\build.ps1 -TargetArch all

# Docker build en push
.\build.ps1 -TargetArch docker -Version 1.0.0
```

### Linux/macOS (Bash)

```bash
# Build voor AMD64
./build.sh amd64

# Build voor ARM64
./build.sh arm64

# Build voor beide architecturen
./build.sh all

# Docker multi-platform build
./build.sh docker
```

### Go direct

```bash
# AMD64
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o sensor-verwerker-amd64 .

# ARM64
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o sensor-verwerker-arm64 .
```

## Configuratie

De applicatie gebruikt de volgende instellingen (zie `main.go`):

- **PostgrestAPI**: `http://postgrest-service.foodchain-db.svc.cluster.local:3000/currentTemperature`
- **JWT Secret**: `super-secret-jwt-key-conform-plan-van-aanpak` (komt overeen met PostgREST configuratie)
- **Truck VIN**: `FC-TRUCK-2026-X99`
- **Truck ID**: `42`
- **Sensoren**: `30` (gespreid over 3 zones)

## JWT Token

De applicatie genereert automatisch JWT tokens met:
- **Algorithm**: HS256 (HMAC-SHA256)
- **Role**: `sensor_admin`
- **Geldigheid**: 1 uur
- **Headers**: Authorization: Bearer {token}

## Docker

### Dockerfile

Gebruik `docker buildx` voor multi-platform builds:

```bash
# Build voor beide architecturen
docker buildx build --platform linux/amd64,linux/arm64 -t sensor-verwerker:latest --push .

# Builden zonder push
docker buildx build --platform linux/amd64,linux/arm64 -t sensor-verwerker:latest --load .
```

### Kubernetes Deployment

Voor Raspberry Pi:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sensor-verwerker
spec:
  replicas: 1
  selector:
    matchLabels:
      app: sensor-verwerker
  template:
    metadata:
      labels:
        app: sensor-verwerker
    spec:
      nodeSelector:
        kubernetes.io/arch: arm64
      containers:
      - name: sensor-verwerker
        image: sensor-verwerker:latest
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "250m"
```

Voor AMD64:

```yaml
nodeSelector:
  kubernetes.io/arch: amd64
```

## Vereisten

- Go 1.21+
- Docker (voor containerization)
- Docker buildx (voor multi-platform builds)
- kubectl (voor Kubernetes deployment)

## Dependencies

- Standaard Go libraires: no external dependencies!

## Troubleshooting

### JWT Token Errors

Als je `403 Unauthorized` krijgt, controleer:
1. JWT Secret in `main.go` matcht de `PGRST_JWT_SECRET` in PostgREST deployment
2. Token expiration (standaard 1 uur)
3. PostgREST JWT configuratie

### Arm64 Build Errors

Bij cross-compilation, zorg ervoor dat `CGO_ENABLED=0` altijd is ingesteld.

### Connection Refused

Controleer of PostgREST API bereikbaar is:

```bash
# Kubernetes
kubectl get svc -n foodchain-db postgrest-service
kubectl port-forward -n foodchain-db svc/postgrest-service 3000:3000
```

## Development

```bash
# Dependencies downloaden
go mod download

# Lokaal runnen
go run main.go

# Tests runnen (als beschikbaar)
go test ./...
```
