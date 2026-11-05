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

## Pruebas

```bash
cd backend && go test ./...                                   # unitarias + integración
cd backend && go test -coverpkg=./internal/... -cover ./...   # con cobertura total
cd backend && go test -short ./...                            # sin la prueba de rendimiento
cd frontend && npm test                                       # componentes y accesibilidad
cd frontend && npm run test:cobertura                         # con cobertura
```

Las pruebas de integración necesitan PostgreSQL encendido. Cada prueba crea
un esquema temporal (`prueba_xxxx`) que se borra al terminar, así que no
tocan los datos de desarrollo. Si la base no responde, esas pruebas aparecen
como omitidas (SKIP) en lugar de fallar.

## Estructura

```
backend/
  cmd/server/         punto de entrada de la API
  cmd/seed/           carga de datos sintéticos
  internal/auth/      claves (PBKDF2) y sesiones
  internal/config/    lectura de .env
  internal/db/        conexión y esquema PostgreSQL
  internal/dbprueba/  base aislada para pruebas de integración
  internal/httpapi/   rutas, middleware y handlers
  internal/producto/  reglas de validación y ciclo de publicación
  internal/seed/      datos de prueba
frontend/
  src/api.js          cliente de la API
  src/sesion/         contexto de sesión
  src/components/     componentes reutilizables (Campo accesible, RutaProtegida)
  src/pages/          páginas
  src/validaciones.js reglas de validación del cliente
docs/                 historias de usuario, DoR/DoD, bitácora, gestión de contenidos
scripts/              respaldo y restauración de la base de datos e imágenes
docker-compose.yml    PostgreSQL para desarrollo
```

## Estado del incremento

| Historia | Descripción | Estado |
|---|---|---|
| HU-01 | Registrar producto | Hecho |
| HU-02 | Aprobar o rechazar producto | Hecho |
| HU-03 | Ver catálogo público | Hecho |
| HU-04 | Archivar producto | Fuera de esta iteración (backlog) |


## Flujo de trabajo con Git

- `main` siempre estable; no se trabaja directo sobre ella.
- Una rama por historia: `feature/HU-01-registro-producto`.
- Commits con el código de la historia y el número del issue:
  `HU-01: valida precio en el formulario (#3)`.
- Cada historia entra a `main` por pull request con la plantilla de revisión.
- Al cerrar la iteración se crea la etiqueta `v1.0.0` con sus notas de versión.

## Respaldo y restauración

```bash
./scripts/respaldo.sh                                   # crea respaldos/bd-*.dump e imagenes-*.tar.gz
./scripts/restaurar.sh respaldos/bd-AAAAMMDD-HHMMSS.dump   # restaura (pide confirmación)
```

Se conservan los últimos 7 respaldos. El procedimiento completo está en
`docs/gestion-contenidos.md`.

## Seguridad y privacidad

- Todos los datos son sintéticos. No se usan datos personales reales.
- No hay secretos en el repositorio: la configuración va en `.env` (ignorado por Git).
- PostgreSQL solo escucha en `127.0.0.1`, no queda expuesto a la red.
- Las claves se guardan con PBKDF2-SHA256 y sal aleatoria, nunca en texto plano.
- La sesión usa una cookie `HttpOnly` y `SameSite=Strict`.
- Los permisos por rol se controlan en la API, no solo en la interfaz.
- El login responde igual si el correo no existe o si la clave es incorrecta.
- Los logs registran método, ruta y estado, sin cuerpos ni cookies.
