apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: grafana-ingress
  namespace: monitoring
  annotations:
    # Ingress class (Traefik)
    kubernetes.io/ingress.class: "traefik"

    # Traefik: gebruik websecure entrypoint en TLS
    traefik.ingress.kubernetes.io/router.entrypoints: websecure
    traefik.ingress.kubernetes.io/router.tls: "true"

    # cert-manager: gebruik de cluster-issuer voor Let’s Encrypt
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
  - hosts:
    - dashboard.test.famivandenberg.nl
    secretName: dashboard-test-famivandenberg-nl-tls
  rules:
  - host: dashboard.test.famivandenberg.nl
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: kube-prometheus-grafana
            port:
              number: 80

