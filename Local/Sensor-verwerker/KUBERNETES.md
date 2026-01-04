# Kubernetes Deployment Guide

## Quick Start

```bash
# Deploy naar foodchain-db namespace
kubectl apply -f k8s-deployment.yaml

# Verify deployment
kubectl get deployment -n foodchain-db sensor-verwerker
kubectl get pods -n foodchain-db -l app=sensor-verwerker
kubectl logs -n foodchain-db -l app=sensor-verwerker -f
```

## Architecture

### Node Selection
- **ARM64**: Raspberry Pi nodes
- **AMD64**: Wijzig nodeSelector in deployment

### Pod Replicas
- **Min**: 1 (via HPA)
- **Max**: 3 (automatic scaling)

### Resources
```
Memory Request: 64 MB
Memory Limit:   128 MB
CPU Request:    100m (0.1 cores)
CPU Limit:      250m (0.25 cores)
```

## Security Context

✅ **Non-root user** (UID 65534 - nobody)
✅ **Read-only filesystem** (except /tmp)
✅ **No privilege escalation**
✅ **Dropped all capabilities**
✅ **seccomp profile**: RuntimeDefault

## Health Checks

### Liveness Probe
- Interval: 30 seconds
- Timeout: 5 seconds
- Failure threshold: 3
- Starts after: 15 seconds

### Startup Probe
- Allows up to 50 seconds for initial startup
- More lenient than liveness probe

## Service Discovery

**Internal DNS:**
```
http://postgrest-service.foodchain-db.svc.cluster.local:3000
```

**Environment Variables:**
```
POSTGREST_HOST=postgrest-service
POSTGREST_PORT=3000
```

## Deployment Commands

### Apply Configuration
```bash
kubectl apply -f k8s-deployment.yaml
```

### View Status
```bash
# Deployment status
kubectl describe deployment sensor-verwerker -n foodchain-db

# Pod status
kubectl get pods -n foodchain-db -l app=sensor-verwerker -o wide

# Pod logs
kubectl logs -n foodchain-db pod/sensor-verwerker-xxxxx -f

# Pod events
kubectl describe pod -n foodchain-db sensor-verwerker-xxxxx
```

### Scaling
```bash
# Manual scaling (disables HPA)
kubectl scale deployment sensor-verwerker -n foodchain-db --replicas=2

# View HPA status
kubectl get hpa -n foodchain-db sensor-verwerker-hpa -w

# Get HPA details
kubectl describe hpa sensor-verwerker-hpa -n foodchain-db
```

### Update Image
```bash
# Rolling update
kubectl set image deployment/sensor-verwerker \
  sensor-verwerker=sensor-verwerker:v1.1 \
  -n foodchain-db

# Check rollout status
kubectl rollout status deployment/sensor-verwerker -n foodchain-db

# Rollback if needed
kubectl rollout undo deployment/sensor-verwerker -n foodchain-db
```

### Port Forwarding (Local Testing)
```bash
# Forward local port to pod
kubectl port-forward -n foodchain-db \
  pod/sensor-verwerker-xxxxx 8080:8080

# Or via service
kubectl port-forward -n foodchain-db \
  svc/sensor-verwerker-service 8080:8080
```

## Affinity & Distribution

### Pod Anti-Affinity
- **Prefers** spreading pods across nodes
- Won't fail if not possible
- Weight: 100 (high preference)

### Node Affinity
- **Prefers** nodes with label `node-type=sensor-node`
- Optional (will schedule elsewhere if needed)

## Pod Disruption Budget

```bash
# View PDB
kubectl get pdb -n foodchain-db

# Check budget
kubectl describe pdb sensor-verwerker-pdb -n foodchain-db
```

**Policy**: maxUnavailable: 0
- Ensures minimum 1 pod always running
- Prevents voluntary disruptions

## Horizontal Pod Autoscaling

### Trigger Conditions
- Memory > 80% of limit
- CPU > 80% of limit

### Scale Up
- Immediate (stabilization: 0 seconds)
- 100% increase or +1 pod per 15 seconds

### Scale Down
- After 5 minutes stable low usage
- 50% reduction per 60 seconds

### View Metrics
```bash
# If metrics-server installed:
kubectl top pod -n foodchain-db -l app=sensor-verwerker
kubectl top pod -n foodchain-db --all-namespaces
```

## ConfigMap Management

```bash
# View ConfigMap
kubectl get configmap -n foodchain-db sensor-verwerker-config -o yaml

# Update ConfigMap
kubectl edit configmap sensor-verwerker-config -n foodchain-db

# Restart pods after ConfigMap update
kubectl rollout restart deployment sensor-verwerker -n foodchain-db
```

## RBAC (Role-Based Access Control)

The deployment includes:
- ServiceAccount: `sensor-verwerker`
- Role: `sensor-verwerker` (read pods, read configmaps)
- RoleBinding: Connects Role to ServiceAccount

```bash
# Verify permissions
kubectl auth can-i get pods --as=system:serviceaccount:foodchain-db:sensor-verwerker -n foodchain-db
```

## Monitoring & Observability

### ServiceMonitor (Prometheus)
```bash
# Requires Prometheus Operator
kubectl get servicemonitor -n foodchain-db
```

### Metrics Endpoint
```
http://sensor-verwerker-service.foodchain-db.svc.cluster.local:8080/metrics
```

### Logging
```bash
# Stream logs
kubectl logs -n foodchain-db -l app=sensor-verwerker -f

# View last 100 lines
kubectl logs -n foodchain-db -l app=sensor-verwerker --tail=100

# Previous pod logs
kubectl logs -n foodchain-db -l app=sensor-verwerker -p
```

## Troubleshooting

### Pod stuck in CrashLoopBackOff
```bash
# Check logs
kubectl logs -n foodchain-db pod/sensor-verwerker-xxxxx --previous

# Check events
kubectl describe pod -n foodchain-db sensor-verwerker-xxxxx

# Check resource availability
kubectl describe node <node-name>
```

### Pod pending
```bash
# Check if node exists
kubectl get nodes -L kubernetes.io/arch

# Check scheduling constraints
kubectl describe pod -n foodchain-db sensor-verwerker-xxxxx

# Check node affinity
kubectl get nodes --show-labels | grep arm64
```

### Connection issues to PostgREST
```bash
# Check DNS resolution
kubectl exec -it -n foodchain-db pod/sensor-verwerker-xxxxx -- \
  nslookup postgrest-service.foodchain-db.svc.cluster.local

# Check connectivity
kubectl exec -it -n foodchain-db pod/sensor-verwerker-xxxxx -- \
  nc -zv postgrest-service 3000
```

## Debug Pod

```bash
# Create debug pod in namespace
kubectl debug -it -n foodchain-db pod/sensor-verwerker-xxxxx --image=busybox:latest

# Or use debugging shell
kubectl run -it --rm debug --image=alpine:latest --restart=Never -n foodchain-db -- sh
```

## Backup & Restore

### Backup Deployment Config
```bash
kubectl get deployment sensor-verwerker -n foodchain-db -o yaml > sensor-verwerker-backup.yaml
```

### Restore from Backup
```bash
kubectl apply -f sensor-verwerker-backup.yaml
```

## Resource Limits Explanation

**Requests**: Minimum guaranteed resources
- Memory: 64 MB - voor normale operatie
- CPU: 100m - voor continuous background work

**Limits**: Maximum allowed resources
- Memory: 128 MB - 2x request (buffer voor peaks)
- CPU: 250m - geeft ruimte voor brief CPU spikes

## Network Policies

Current configuration allows:
- ✅ Egress naar PostgREST (port 3000)
- ✅ DNS queries (port 53)
- ❌ Ingress alleen van monitoring namespace
- ❌ All other traffic blocked

## Best Practices Applied

✅ Non-root user
✅ Read-only filesystem
✅ Resource limits
✅ Health probes
✅ Security context
✅ Pod disruption budgets
✅ RBAC
✅ Service accounts
✅ ConfigMaps for config
✅ Environment variable injection
✅ Pod anti-affinity
✅ Graceful shutdown
✅ Proper logging

## Production Checklist

- [ ] Image pushed to registry
- [ ] Registry credentials configured (if private)
- [ ] Node labels set (kubernetes.io/arch=arm64)
- [ ] PostgREST service running
- [ ] PostgreSQL database accessible
- [ ] Networking policies in place
- [ ] Monitoring/logging configured
- [ ] Resource quotas set for namespace
- [ ] RBAC policies reviewed
- [ ] Backup/restore tested
