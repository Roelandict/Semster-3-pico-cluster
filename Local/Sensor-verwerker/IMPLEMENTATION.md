# Implementatie Sensor-verwerker JWT & Multi-Platform

## Wat is opgelost

### 1. JWT Token Ondersteuning
- **Probleem**: De applicatie stuurde data naar PostgREST zonder authentication
- **Oplossing**: 
  - Implementatie van HMAC-SHA256 JWT token generatie
  - Tokens worden automatisch gegenereerd met 1-uur expiratie
  - Authorization header wordt met Bearer token meegezonden

### 2. Multi-Platform Ondersteuning
- **AMD64** (Windows/Linux x86-64): Standaard build
- **ARM64** (Raspberry Pi): Cross-compilation support

### 3. Build Scripts
- `build.ps1`: PowerShell script voor Windows (AMD64/ARM64/Docker)
- `build.sh`: Bash script voor Linux/macOS
- Docker multi-platform builds via `docker buildx`

## Wijzigingen

### main.go
```
✓ JWT token generatie (HMAC-SHA256)
✓ Authorization Bearer header
✓ JWTClaims struct met role, iat, exp
✓ generateJWT() functie
✓ Cross-platform compatible code
```

### Dockerfile
```
✓ Multi-platform build support
✓ Automatische TARGETARCH detectie
✓ Alpine Linux basis voor kleine images
✓ ARM64 + AMD64 compatibiliteit
```

### Build Scripts
```
✓ build.ps1 - Windows PowerShell build
✓ build.sh - Linux/macOS bash build
✓ Ondersteuning voor beide architecturen
✓ Docker multi-platform builds
```

## JWT Token Details

**Gegenereerde tokens:**
- Algorithm: HS256 (HMAC-SHA256)
- Secret: `super-secret-jwt-key-conform-plan-van-aanpak`
- Rol: `sensor_admin`
- Expiratie: 1 uur na creatie
- Header: `Authorization: Bearer <token>`

**Voorbeeld JWT opbouw:**
```
Header.Claims.Signature
{
  "alg": "HS256",
  "typ": "JWT"
}.{
  "role": "sensor_admin",
  "iat": 1704371000,
  "exp": 1704374600
}.signature
```

## Compilatie Instructies

### Windows (PowerShell)
```powershell
# Navigeer naar de directory
cd C:\Users\Roela\code\pico\Local\Sensor-verwerker

# Build voor AMD64 (standaard)
.\build.ps1

# Build voor ARM64 (Raspberry Pi)
.\build.ps1 -TargetArch arm64

# Build voor beide
.\build.ps1 -TargetArch all

# Docker build
.\build.ps1 -TargetArch docker -Version 1.0.0
```

### Linux/macOS
```bash
cd ~/code/pico/Local/Sensor-verwerker

# Build voor AMD64
./build.sh amd64

# Build voor ARM64
./build.sh arm64

# Docker build
./build.sh docker
```

## Binaries

Na compilatie:
- `sensor-verwerker.exe` (AMD64, Windows buildable)
- `sensor-verwerker-arm64` (ARM64, Raspberry Pi compatible)

## Docker Deployment

### Multi-Platform Build
```bash
docker buildx build --platform linux/amd64,linux/arm64 \
  -t sensor-verwerker:latest \
  -t sensor-verwerker:1.0.0 \
  --push .
```

### Kubernetes Deployment (Raspberry Pi)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sensor-verwerker
  namespace: foodchain-db
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

## Testresultaten

✓ AMD64 Binary compilatie: **Succesvol** (8.46 MB)
✓ ARM64 Binary compilatie: **Succesvol** (7.86 MB)
✓ JWT Token generatie: **Succesvol**
✓ Applicatie startup: **Succesvol**

## Volgende Stappen

1. Docker image builden en pushen
2. Kubernetes deployment updaten met sensor-verwerker
3. Raspberry Pi node labels controleren
4. Sensor-verwerker pods deployen op Raspberry Pi's

## Troubleshooting

### "GOOS/GOARCH not recognized"
- Zorg dat CGO_ENABLED=0 ingesteld is
- Gebruik de build scripts (not direct `go build`)

### ARM64 Cross-compilation errors
- Zorg dat Go 1.21+ geïnstalleerd is
- Controleer GOARCH=arm64 ingesteld is

### JWT Connection Errors
- Controleer JWT secret in main.go matcht PostgREST `PGRST_JWT_SECRET`
- Controleer PostgREST API bereikbaar is
- Check token expiratie (standaard 1 uur)
