#!/usr/bin/env bash
# Restaura un respaldo de la base de datos (y, si se indica, de las imágenes).
# Uso:
#   ./scripts/restaurar.sh respaldos/bd-AAAAMMDD-HHMMSS.dump
#   ./scripts/restaurar.sh respaldos/bd-AAAAMMDD-HHMMSS.dump respaldos/imagenes-AAAAMMDD-HHMMSS.tar.gz
#
# ATENCIÓN: reemplaza los datos actuales por los del respaldo.

set -euo pipefail

RAIZ="$(cd "$(dirname "$0")/.." && pwd)"
cd "$RAIZ"

if [ $# -lt 1 ] || [ ! -f "$1" ]; then
  echo "Indique un archivo .dump existente. Ejemplo: ./scripts/restaurar.sh respaldos/bd-20261105-080000.dump"
  exit 1
fi

# Carga .env sin pisar variables que ya estén definidas en la terminal
# (igual que el backend: el entorno tiene prioridad sobre el archivo).
if [ -f .env ]; then
  while IFS='=' read -r clave valor; do
    case "$clave" in ''|\#*) continue ;; esac
    if [ -z "${!clave+x}" ]; then export "$clave=$valor"; fi
  done < .env
fi

read -r -p "Esto reemplaza los datos actuales. ¿Continuar? (s/N) " respuesta
[ "$respuesta" = "s" ] || { echo "Cancelado."; exit 0; }

if command -v pg_restore >/dev/null 2>&1 && [ -n "${DATABASE_URL:-}" ]; then
  pg_restore --clean --if-exists --no-owner --dbname="$DATABASE_URL" "$1"
else
  docker compose exec -T db pg_restore --clean --if-exists --no-owner \
    -U "$POSTGRES_USER" -d "$POSTGRES_DB" < "$1"
fi
echo "Base de datos restaurada desde $1"

if [ $# -ge 2 ] && [ -f "$2" ]; then
  UPLOADS="backend/${DIR_UPLOADS:-data/uploads}"
  mkdir -p "$UPLOADS"
  tar -xzf "$2" -C "$UPLOADS"
  echo "Imágenes restauradas en $UPLOADS"
fi
