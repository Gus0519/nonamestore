# NoNameStore

Tienda en línea desplegada sobre Kubernetes con un pipeline CI/CD completo. El frontend está hecho en Svelte y el backend en Go (Gin), con MariaDB como base de datos y una pila de observabilidad basada en Prometheus, Grafana y OpenSearch.

---

## Arquitectura general

```
GitHub Push
    │
    ▼
Tekton (CI)
  ├─ git-clone
  ├─ kaniko → build frontend + backend → Docker Hub
  ├─ kubectl rollout restart (staging)
  └─ git push → actualiza tags en k8s/staging/

ArgoCD (CD)
  ├─ nonamestore-staging  →  namespace: staging
  └─ nonamestore-prod     →  namespace: prod
```

El flujo completo es: push a `main` → webhook a Tekton → build de imágenes con el hash del commit como tag → ArgoCD detecta el cambio en el repo y sincroniza el clúster.

---

## Stack

| Capa | Tecnología |
|---|---|
| Frontend | Svelte + Vite, servido con Nginx |
| Backend | Go + Gin Gonic |
| Base de datos | MariaDB 10.6 (StatefulSet) |
| Contenedores | Kaniko (builds sin Docker daemon) |
| CI | Tekton Pipelines + Tekton Triggers |
| CD | ArgoCD |
| Métricas | Prometheus + Grafana + kube-state-metrics |
| Logs | Fluentd → OpenSearch → OpenSearch Dashboards |
| Alertas | Alertmanager → Slack |
| Autoscaling | VPA (Vertical Pod Autoscaler) en modo Auto para el backend |
| Ingress | nginx-ingress-controller |

---

## Estructura del repositorio

```
.
├── src/
│   ├── frontend/          # App Svelte
│   │   └── Containerfile
│   └── backend/           # API Go
│       └── Containerfile
├── k8s/
│   ├── staging/           # Manifiestos del entorno de pruebas
│   │   ├── backend/
│   │   ├── frontend/
│   │   ├── mariadb/
│   │   └── misc/          # Ingress
│   └── prod/              # Manifiestos de producción (misma estructura)
├── infra/
│   ├── monitoring/        # Prometheus, Grafana, OpenSearch, Alertmanager, Fluentd
│   ├── metrics-server/    # Metrics Server + VPA
│   └── vpa/
└── tekton/
    ├── active/            # Pipeline y triggers en uso
    └── old_tests/         # Runs manuales de prueba (archivados)
```

---

## Namespaces

| Namespace | Contenido |
|---|---|
| `staging` | Frontend, Backend, MariaDB, Ingress |
| `prod` | Misma estructura que staging |
| `monitoring` | Prometheus, Grafana, Alertmanager, Fluentd, OpenSearch |
| `default` | Tekton, ServiceAccounts, Tekton Triggers |
| `argocd` | ArgoCD |
| `kube-system` | Metrics Server |

---

## Setup inicial

### 1. Secretos necesarios

Antes de aplicar cualquier manifiesto hay que crear los secretos a mano. Tekton y el backend los esperan en el clúster.

```bash
# Credenciales de Docker Hub para Kaniko
kubectl create secret generic docker-credentials \
  --from-file=.dockerconfigjson=$HOME/.docker/config.json \
  --type=kubernetes.io/dockerconfigjson

# Token de GitHub para que el pipeline haga push al repo
kubectl create secret generic github-auth-secret \
  --from-literal=token=<TU_GITHUB_PAT>

# Webhook de Slack para Alertmanager
kubectl create secret generic alertmanager-slack-webhook \
  --from-literal=webhook-url=<TU_WEBHOOK_URL> \
  -n monitoring
```

La contraseña de MariaDB está en `k8s/staging/mariadb/mariadb-secret.yaml` en base64. El valor por defecto es `Admin123`. Cámbialo antes de ir a producción.

### 2. Almacenamiento

El PV de Tekton usa `hostPath` en `/mnt/data/tekton`. El PV de MariaDB usa NFS apuntando a `192.168.122.1` (host de libvirt). Ajusta ambas rutas y la IP según tu entorno antes de aplicar.

```bash
kubectl apply -f tekton/active/00-pv.yaml
kubectl apply -f tekton/active/01-pvc.yaml
kubectl apply -f k8s/staging/mariadb/mariadb-storage.yaml
```

### 3. Tekton

```bash
# Instalar Tekton Pipelines y Triggers (si no están instalados)
kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
kubectl apply --filename https://storage.googleapis.com/tekton-releases/triggers/latest/release.yaml

# Aplicar la configuración del proyecto
kubectl apply -f tekton/active/tekton-setup.yaml
kubectl apply -f tekton/active/trigger-rbac.yaml
kubectl apply -f tekton/active/deploy-rbac.yaml
kubectl apply -f tekton/active/nonamestore-pipeline.yaml
kubectl apply -f tekton/active/github-trigger.yaml
```

### 4. ArgoCD

```bash
kubectl apply -f tekton/active/04-argo-staging.yaml
kubectl apply -f tekton/active/argo-prod.yaml
```

ArgoCD leerá `k8s/staging/` y `k8s/prod/` y sincronizará automáticamente cuando el pipeline empuje cambios.

### 5. Monitoreo

```bash
kubectl apply -f infra/monitoring/
```

Esto despliega Prometheus, Grafana, Alertmanager, Fluentd y OpenSearch en el namespace `monitoring`.

---

## Puertos expuestos (NodePort)

| Servicio | Puerto |
|---|---|
| Prometheus | 30090 |
| Grafana | 30030 |
| Alertmanager | 30093 |
| OpenSearch Dashboards | 30561 |

---

## Pipeline CI/CD en detalle

El pipeline `build-and-push-nonamestore` tiene cinco etapas:

1. **fetch-repository** — clona el repo con `git-clone`
2. **build-frontend** / **build-backend** — construye las imágenes con Kaniko usando el hash del commit como tag (`gusgus890/nonamestore-frontend:<hash>`)
3. **auto-deploy** — hace `kubectl rollout restart` en staging para forzar el pull de la nueva imagen
4. **update-gitops-manifests** — actualiza el tag de imagen en los YAMLs de `k8s/staging/` y hace push a `main` para que ArgoCD lo detecte

El webhook de GitHub apunta al `EventListener` en el clúster. Para exponerlo al exterior en un entorno local se puede usar `kubectl port-forward` o `ngrok`.

---

## Acceso a la aplicación

La tienda se sirve bajo el host `nonamestore.local`. Agrega la IP del nodo a tu `/etc/hosts`:

```
<IP_DEL_NODO>  nonamestore.local
```

Las rutas configuradas en el Ingress son:

- `/` → frontend (Svelte)
- `/products` → backend (Go)
- `/health` → healthcheck del backend

---

## Variables de entorno del backend

El backend recibe la configuración de la base de datos a través de un ConfigMap y un Secret:

| Variable | Fuente | Valor por defecto |
|---|---|---|
| `DB_HOST` | ConfigMap `nonamestore-config` | `mariadb-service` |
| `DB_PORT` | ConfigMap `nonamestore-config` | `3306` |
| `DB_NAME` | ConfigMap `nonamestore-config` | `prueba_persistencia` |
| `DB_USER` | hardcoded | `root` |
| `DB_PASSWORD` | Secret `mariadb-secret` | — |

---

## Alertas configuradas

| Alerta | Condición |
|---|---|
| `HighCPUUsage` | CPU > 80% sostenido durante 1 minuto en pods de nonamestore |
| `PodCrashLooping` | Cualquier pod en estado `CrashLoopBackOff` |

Las notificaciones van al canal `#alertas-nonamestore` en Slack.


# Prueba pipeline
