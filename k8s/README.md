# Kubernetes Deployment

Deploy BigPanda Super Agent to Kubernetes.

## Prerequisites

- Kubernetes 1.20+
- kubectl configured
- BigPanda API credentials

## Quick Start

### 1. Create Namespace

```bash
kubectl create namespace bigpanda
```

### 2. Update Credentials

Edit `secret.yaml` and add your BigPanda credentials:

```bash
kubectl edit -f secret.yaml
```

Or create from command line:

```bash
kubectl create secret generic bigpanda-credentials \
  --from-literal=token=YOUR_TOKEN \
  --from-literal=app_key=YOUR_APP_KEY \
  -n bigpanda
```

### 3. Deploy

Using kubectl:

```bash
kubectl apply -f rbac.yaml -n bigpanda
kubectl apply -f secret.yaml -n bigpanda
kubectl apply -f configmap.yaml -n bigpanda
kubectl apply -f pvc.yaml -n bigpanda
kubectl apply -f deployment.yaml -n bigpanda
kubectl apply -f service.yaml -n bigpanda
```

Or using kustomize:

```bash
kubectl apply -k .
```

### 4. Verify

```bash
# Check pods
kubectl get pods -n bigpanda

# Check logs
kubectl logs -f deployment/bigpanda-agent -n bigpanda

# Check health
kubectl exec -it deployment/bigpanda-agent -n bigpanda -- \
  wget -q -O- http://localhost:8443/health
```

## Configuration

### ConfigMap

Edit `configmap.yaml` to customize:
- BigPanda API endpoint
- Queue settings
- Logging level
- Module configuration

### Resources

Adjust resource limits in `deployment.yaml`:

```yaml
resources:
  requests:
    cpu: 500m
    memory: 256Mi
  limits:
    cpu: 2000m
    memory: 1Gi
```

### Persistent Storage

Modify `pvc.yaml` for storage size:

```yaml
resources:
  requests:
    storage: 10Gi
```

## Exposing Services

### SNMP Traps (UDP 162)

For SNMP traps, use NodePort or LoadBalancer:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: bigpanda-agent-snmp
spec:
  type: NodePort
  selector:
    app: bigpanda-agent
  ports:
  - name: snmp
    port: 162
    targetPort: 162
    nodePort: 30162
    protocol: UDP
```

### Web UI (HTTPS 8443)

Use Ingress for Web UI:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: bigpanda-agent
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt
spec:
  tls:
  - hosts:
    - agent.example.com
    secretName: bigpanda-tls
  rules:
  - host: agent.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: bigpanda-agent
            port:
              number: 8443
```

## Monitoring

### Prometheus

The agent exposes metrics on port 9090:

```yaml
apiVersion: v1
kind: ServiceMonitor
metadata:
  name: bigpanda-agent
spec:
  selector:
    matchLabels:
      app: bigpanda-agent
  endpoints:
  - port: metrics
    interval: 30s
```

## Scaling

For high availability, increase replicas:

```bash
kubectl scale deployment bigpanda-agent --replicas=3 -n bigpanda
```

Note: Each replica needs its own SNMP listening address or use a LoadBalancer.

## Updating

### Rolling Update

```bash
kubectl set image deployment/bigpanda-agent \
  agent=reggiejtech/super-agent:v1.1.0 \
  -n bigpanda
```

### Update Configuration

```bash
kubectl edit configmap bigpanda-config -n bigpanda
kubectl rollout restart deployment/bigpanda-agent -n bigpanda
```

## Troubleshooting

### Pod Not Starting

```bash
kubectl describe pod -l app=bigpanda-agent -n bigpanda
kubectl logs -l app=bigpanda-agent -n bigpanda --previous
```

### Permission Issues

Check RBAC:

```bash
kubectl auth can-i get configmaps \
  --as=system:serviceaccount:bigpanda:bigpanda-agent \
  -n bigpanda
```

### Network Issues

Test from within pod:

```bash
kubectl exec -it deployment/bigpanda-agent -n bigpanda -- sh
wget -q -O- http://localhost:8443/health
```

## Cleanup

```bash
kubectl delete -k . -n bigpanda
kubectl delete namespace bigpanda
```

Or:

```bash
kubectl delete -f deployment.yaml -n bigpanda
kubectl delete -f service.yaml -n bigpanda
kubectl delete -f pvc.yaml -n bigpanda
kubectl delete -f configmap.yaml -n bigpanda
kubectl delete -f secret.yaml -n bigpanda
kubectl delete -f rbac.yaml -n bigpanda
```
