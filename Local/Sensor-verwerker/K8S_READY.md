# Sensor-verwerker Kubernetes Deployment - Final Summary

## 📋 Kubernetes Resources

De `k8s-deployment.yaml` bevat **9 resources**:

✅ **Namespace** - foodchain-db
✅ **ServiceAccount** - sensor-verwerker  
✅ **Role** - sensor-verwerker (RBAC)
✅ **RoleBinding** - sensor-verwerker
✅ **ConfigMap** - sensor-verwerker-config
✅ **Deployment** - sensor-verwerker (met security context)
✅ **Service** - sensor-verwerker-service (ClusterIP)
✅ **PodDisruptionBudget** - sensor-verwerker-pdb
✅ **HorizontalPodAutoscaler** - sensor-verwerker-hpa

## 🔐 Security Features

### Container Security
```
✅ Non-root user (UID 65534 - nobody)
✅ Read-only root filesystem (except /tmp, /var/run)
✅ No privilege escalation allowed
✅ All capabilities dropped
✅ seccomp profile: RuntimeDefault
```

### Pod Security
```
✅ Security context at pod level
✅ Security context at container level
✅ Service account with minimal RBAC
✅ Network policies supported
✅ Resource limits enforced
```

## 📦 Deployment Configuration

### Node Selection
```
nodeSelector:
  kubernetes.io/arch: arm64    # Raspberry Pi
```

**Wijzigen naar AMD64:**
```yaml
nodeSelector:
  kubernetes.io/arch: amd64
```

### Resource Allocation
```
Requests:
  Memory: 64 MB
  CPU:    100m

Limits:
  Memory: 128 MB  
  CPU:    250m
```

### Health Monitoring
```
Liveness Probe:  Every 30s (starts after 15s)
Startup Probe:   Every 5s (allows 50s total)
```

## 🚀 Quick Deployment

```bash
# 1. Build Docker image for multi-platform
docker buildx build --platform linux/amd64,linux/arm64 \
  -t sensor-verwerker:latest \
  --push .

# 2. Deploy to Kubernetes
kubectl apply -f k8s-deployment.yaml

# 3. Verify
kubectl get pods -n foodchain-db -l app=sensor-verwerker
kubectl logs -n foodchain-db -l app=sensor-verwerker -f
```

## 📊 Auto-Scaling (HPA)

```
Min Replicas: 1
Max Replicas: 3

Triggers:
  - Memory > 80% → Scale up
  - CPU > 80% → Scale up
  - Memory/CPU < 80% for 5min → Scale down
```

```bash
# View HPA status
kubectl get hpa -n foodchain-db -w
```

## 🔗 Service Discovery

**Internal DNS:**
```
postgrest-service.foodchain-db.svc.cluster.local:3000
```

**Environment Variables in Pod:**
```
POSTGREST_HOST=postgrest-service
POSTGREST_PORT=3000
```

## 📝 Configuration

Via **ConfigMap** `sensor-verwerker-config`:
```
TZ=Europe/Amsterdam
LOG_LEVEL=INFO
TRUCK_VIN=FC-TRUCK-2026-X99
TRUCK_ID=42
SENSOR_COUNT=30
```

**Update config:**
```bash
kubectl edit configmap sensor-verwerker-config -n foodchain-db
kubectl rollout restart deployment sensor-verwerker -n foodchain-db
```

## 🛡️ Pod Disruption Budget

```
maxUnavailable: 0
→ Minimum 1 pod always running
→ Prevents accidental eviction
```

## 📈 Monitoring (Optional)

**Install Prometheus Operator first:**
```bash
kubectl apply -f k8s-monitoring.yaml
```

**Includes:**
- ServiceMonitor for Prometheus scraping
- PrometheusRules for alerting
- Alerts for: Down, HighMemory, HighCPU

## 🔄 Rolling Updates

```bash
# Update image
kubectl set image deployment/sensor-verwerker \
  sensor-verwerker=sensor-verwerker:v1.1 \
  -n foodchain-db

# Check status
kubectl rollout status deployment/sensor-verwerker -n foodchain-db

# Rollback if needed
kubectl rollout undo deployment/sensor-verwerker -n foodchain-db
```

## 🧹 Cleanup

```bash
# Delete all resources
kubectl delete -f k8s-deployment.yaml

# Or just the deployment
kubectl delete deployment sensor-verwerker -n foodchain-db
```

## ✨ Best Practices Implemented

✅ Non-root containers
✅ Immutable infrastructure (read-only FS)
✅ Resource limits
✅ Health checks
✅ RBAC
✅ Service accounts
✅ ConfigMaps for config
✅ Auto-scaling
✅ Pod disruption budgets
✅ Affinity rules
✅ Graceful shutdown
✅ Proper logging
✅ Security context
✅ Capability dropping

## 📚 Related Files

- `k8s-deployment.yaml` - Main Kubernetes manifests
- `k8s-monitoring.yaml` - Optional monitoring setup
- `KUBERNETES.md` - Detailed Kubernetes guide
- `Dockerfile` - Multi-platform container image
- `docker-compose.yml` - Local development (optional)

## 🎯 Deployment Workflow

1. **Build & Push Image**
   ```bash
   docker buildx build --platform linux/amd64,linux/arm64 \
     -t <registry>/sensor-verwerker:latest --push .
   ```

2. **Apply Kubernetes Manifests**
   ```bash
   kubectl apply -f k8s-deployment.yaml
   ```

3. **Verify Deployment**
   ```bash
   kubectl get all -n foodchain-db -l app=sensor-verwerker
   ```

4. **Monitor Logs**
   ```bash
   kubectl logs -f -n foodchain-db -l app=sensor-verwerker
   ```

5. **Troubleshoot if Needed**
   ```bash
   kubectl describe pod -n foodchain-db sensor-verwerker-xxxxx
   kubectl debug pod -n foodchain-db sensor-verwerker-xxxxx
   ```

## 🔍 Common Commands

```bash
# Status
kubectl get pods -n foodchain-db -l app=sensor-verwerker -o wide

# Logs
kubectl logs -n foodchain-db -l app=sensor-verwerker -f --tail=50

# Port forward (testing)
kubectl port-forward -n foodchain-db svc/sensor-verwerker-service 8080:8080

# Execute command in pod
kubectl exec -it -n foodchain-db pod/sensor-verwerker-xxxxx -- /bin/sh

# Describe pod
kubectl describe pod -n foodchain-db sensor-verwerker-xxxxx

# Delete pod (will restart via deployment)
kubectl delete pod -n foodchain-db sensor-verwerker-xxxxx
```

---

**Status**: ✅ Production Ready
**Version**: 1.0.0
**Date**: 2026-01-04
