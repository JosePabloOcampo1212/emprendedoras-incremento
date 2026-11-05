#!/usr/bin/env bash
# Respaldo de la base de datos y de las imágenes subidas.
# Uso (desde la raíz del repositorio):  ./scripts/respaldo.sh
#
# Si pg_dump está instalado y existe DATABASE_URL, lo usa directo.
# Si no, usa el contenedor de docker-compose.
# Conserva los últimos 7 respaldos de cada tipo.

set -euo pipefail

RAIZ="$(cd "$(dirname "$0")/.." && pwd)"
cd "$RAIZ"

# Carga .env sin pisar variables que ya estén definidas en la terminal
# (igual que el backend: el entorno tiene prioridad sobre el archivo).
if [ -f .env ]; then
  while IFS='=' read -r clave valor; do
    case "$clave" in ''|\#*) continue ;; esac
    if [ -z "${!clave+x}" ]; then export "$clave=$valor"; fi
  done < .env
fi

DESTINO="${DESTINO_RESPALDOS:-respaldos}"
CONSERVAR=7
FECHA="$(date +%Y%m%d-%H%M%S)"
mkdir -p "$DESTINO"

ARCHIVO_BD="$DESTINO/bd-$FECHA.dump"
if command -v pg_dump >/dev/null 2>&1 && [ -n "${DATABASE_URL:-}" ]; then
  pg_dump --format=custom --no-owner --file="$ARCHIVO_BD" "$DATABASE_URL"
else
  docker compose exec -T db pg_dump --format=custom --no-owner \
    -U "$POSTGRES_USER" -d "$POSTGRES_DB" > "$ARCHIVO_BD"
fi
echo "Base de datos respaldada en $ARCHIVO_BD ($(du -h "$ARCHIVO_BD" | cut -f1))"

UPLOADS="backend/${DIR_UPLOADS:-data/uploads}"
if [ -d "$UPLOADS" ]; then
  ARCHIVO_IMG="$DESTINO/imagenes-$FECHA.tar.gz"
  tar -czf "$ARCHIVO_IMG" -C "$UPLOADS" .
  echo "Imágenes respaldadas en $ARCHIVO_IMG ($(du -h "$ARCHIVO_IMG" | cut -f1))"
fi

# Retención: borra los respaldos más viejos y deja los últimos $CONSERVAR.
for patron in "bd-*.dump" "imagenes-*.tar.gz"; do
  # shellcheck disable=SC2012
  { ls -1t "$DESTINO"/$patron 2>/dev/null || true; } | tail -n +$((CONSERVAR + 1)) | while read -r viejo; do
    rm -f "$viejo"
    echo "Eliminado respaldo antiguo: $viejo"
  done
done
