# Plataforma de la Asociación de Emprendedoras — Incremento 1

Proyecto No. 1 del curso 03637 Gestión de Sitios Web (UNED, III cuatrimestre 2026).

El incremento implementa un flujo completo de publicación de productos:
una emprendedora registra un producto, la administradora de la asociación lo
revisa y, si lo aprueba, el producto aparece en el catálogo público.

**Objetivo de la iteración:** al cierre de la iteración, una emprendedora puede
publicar un producto que, tras la aprobación de la administradora, queda visible
en el catálogo público.

## Tecnologías

| Parte | Tecnología |
|---|---|
| Backend | Go 1.24, biblioteca estándar `net/http`, driver `github.com/jackc/pgx/v5` |
| Base de datos | PostgreSQL 16 (en Docker) |
| Frontend | React 19 con Vite, React Router |
| Pruebas | `go test` (unitarias e integración), Vitest + Testing Library (componentes) |

## Requisitos

- Go 1.24 o superior
- Node.js 20 o superior
- Docker con Docker Compose (para PostgreSQL)

Si ya hay un PostgreSQL instalado, se puede usar ese en lugar de Docker
definiendo `DATABASE_URL` en el `.env`.

## Puesta en marcha

```bash
# 1. Configuración local
cp .env.example .env
#    Editar .env: poner claves inventadas en POSTGRES_PASSWORD y SEED_CLAVE

# 2. Base de datos
docker compose up -d   # PostgreSQL 16 en localhost:5432

# 3. Backend (terminal 1)
cd backend
go mod tidy            # descarga dependencias la primera vez
go run ./cmd/seed      # crea las tablas y carga datos sintéticos
go run ./cmd/server    # API en http://localhost:8080

# 4. Frontend (terminal 2)
cd frontend
npm install
npm run dev            # sitio en http://localhost:5173
```

En desarrollo, Vite redirige `/api` al backend, así que el sitio y la API
comparten origen y la cookie de sesión funciona sin configuración extra.

## Cuentas de prueba

Todas son ficticias y usan el dominio reservado `.test`. La clave de todas es
el valor de `SEED_CLAVE` en el `.env`.

| Correo | Rol |
|---|---|
| admin@asociacion.test | administradora |
| emprendedora1@asociacion.test | emprendedora |
| emprendedora2@asociacion.test | emprendedora |

