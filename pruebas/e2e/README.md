# Prueba funcional de punta a punta

Recorre el flujo completo en Chrome: catálogo → ingreso → registro (con teclado) →
revisión → aprobación y rechazo → catálogo filtrado → vista móvil. Guarda una
captura de cada paso en `capturas/`.

## Requisitos

- Python 3.10 o superior
- Google Chrome instalado

```bash
pip install -r pruebas/e2e/requirements.txt
```

## Cómo correrla

La prueba espera los datos del seed recién cargados (3 productos publicados y 1 pendiente).

```bash
# 1. Base de datos limpia
docker compose down -v && docker compose up -d
cd backend && go run ./cmd/seed && go run ./cmd/server      # terminal 1

# 2. Versión de producción del frontend
cd frontend && npm run build && npm run preview                # terminal 2

# 3. Prueba
python3 pruebas/e2e/flujo.py                                   # terminal 3
```

Para usar Chromium de Playwright en lugar de Chrome:
`playwright install chromium` y luego `E2E_NAVEGADOR=chromium python3 pruebas/e2e/flujo.py`.

El resultado esperado es `12/12 pasos correctos`.
