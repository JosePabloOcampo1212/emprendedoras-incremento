# Plan de pruebas

## 1. Alcance

Se prueba el incremento completo: inicio de sesión por rol, HU-01 (registrar producto),
HU-02 (aprobar o rechazar) y HU-03 (catálogo público). La HU-04 está fuera de alcance.

## 2. Entorno

| Elemento | Versión / detalle |
|---|---|
| Backend | Go 1.24, API en `localhost:8080` |
| Base de datos | PostgreSQL 16 (Docker `postgres:16-alpine`) |
| Frontend | React 19 + Vite, versión de producción con `npm run build` y `vite preview` (`localhost:4173`) |
| Navegador | Google Chrome (versión estable) |
| Herramientas | `go test`, Vitest + Testing Library, axe-core, Playwright (Python), Lighthouse 12 |
| Datos | Sintéticos, cargados con `go run ./cmd/seed`. Correos con dominio reservado `.test` |

## 3. Niveles de prueba

| Nivel | Qué cubre | Herramienta | Dónde está |
|---|---|---|---|
| Unitarias | Reglas de validación, ciclo de estados, hash de claves | `go test`, Vitest | `backend/internal/producto`, `backend/internal/auth`, `frontend/src/validaciones.test.js` |
| Componente | Formularios, mensajes de error, foco, anuncios a lectores de pantalla | Vitest + Testing Library | `frontend/src/**/*.test.jsx` |
| Integración | API completa contra PostgreSQL real (esquema temporal por prueba) | `go test` | `backend/internal/httpapi/*_test.go` |
| Funcionales (punta a punta) | Flujo completo en el navegador con la versión de producción | Playwright | `pruebas/e2e/flujo.py` |
| No funcionales | Accesibilidad, rendimiento, seguridad básica, eficiencia | axe-core, Lighthouse, `go test` | ver sección 5 |
| Aceptación | Cada criterio de aceptación de las historias | Combinación de las anteriores | ver sección 4 |

## 4. Casos de prueba por criterio de aceptación

| Caso | Historia / criterio | Descripción | Resultado esperado | Prueba automática |
|---|---|---|---|---|
| CP-01 | HU-01 CA2 | Registrar producto válido con imagen | 201, estado "pendiente", 1 registro en historial, imagen accesible | `TestCrearProducto_ValidoQuedaPendiente` |
| CP-02 | HU-01 CA1 | Nombre vacío y precio negativo | 422 con error por campo, no se guarda nada | `TestCrearProducto_InvalidoNoGuarda` |
| CP-03 | HU-01 CA1 | Precio con texto | 422, error en precio | `TestCrearProducto_PrecioNoNumerico` |
| CP-04 | HU-01 CA3 | Archivo GIF renombrado como .png | 422, error en imagen (se revisa el contenido) | `TestCrearProducto_ImagenTipoNoPermitido` |
| CP-05 | HU-01 CA3 | Imagen de más de 2 MB | 422, error en imagen | `TestCrearProducto_ImagenMuyPesada` |
| CP-06 | HU-01 CA3 | Imagen sin texto alternativo | 422, error en imagenAlt | `TestCrearProducto_ImagenSinAlt` |
| CP-07 | HU-01 CA1 | Categoría inexistente | 422, error en categoría | `TestCrearProducto_CategoriaInexistente` |
| CP-08 | Seguridad | Registrar sin sesión / como administradora | 401 / 403 | `TestCrearProducto_Permisos` |
| CP-09 | Privacidad | Cada emprendedora ve solo lo suyo | Lista sin productos ajenos | `TestMisProductos_SoloPropios` |
| CP-10 | HU-02 CA1 | Lista de pendientes según rol | Admin 200; emprendedora 403; sin sesión 401 | `TestPendientes_SoloAdministradora` |
| CP-11 | HU-02 CA2 | Aprobar | Estado "publicado" + historial | `TestAprobar_PublicaProducto` |
| CP-12 | HU-02 CA3 | Rechazar sin y con motivo | Sin motivo 422 y no cambia; con motivo pasa a "borrador" y la emprendedora lo ve | `TestRechazar_ConMotivo` |
| CP-13 | HU-02 | Revisar dos veces el mismo producto | 409 | `TestRevisar_DosVecesDaConflicto` |
| CP-14 | HU-02 / seguridad | Id inexistente, id inválido, emprendedora aprobando | 404 / 404 / 403 | `TestRevisar_NoEncontradoYPermisos` |
| CP-15 | HU-03 CA1 | Catálogo sin sesión | Solo publicados, sin estado ni motivo | `TestCatalogo_SoloPublicados` |
| CP-16 | HU-01 a HU-03 | Flujo registrar → aprobar → catálogo | Aparece solo después de aprobar | `TestCatalogo_FlujoCompleto` |
| CP-17 | HU-03 CA2 | Filtro por categoría | Solo la categoría elegida; inválida 400 | `TestCatalogo_FiltroCategoria` |
| CP-18 | No funcional | Rendimiento con 500 publicados | p95 < 300 ms | `TestCatalogo_Rendimiento` |
| CP-19 | HU-01 CA4 | Registro completo solo con teclado | Se llega a cada campo con Tab y se envía con Enter | `pruebas/e2e/flujo.py` paso 4 |
| CP-20 | HU-03 CA2 | El filtro no recarga la página | No hay navegación nueva; la URL guarda el filtro | `pruebas/e2e/flujo.py` paso 8 |
| CP-21 | HU-03 CA3 | Lighthouse accesibilidad | ≥ 90 | Lighthouse manual |

## 5. Pruebas no funcionales

| Aspecto | Prueba | Criterio |
|---|---|---|
| Accesibilidad | axe-core sobre las 4 pantallas (`accesibilidad.test.jsx`) y Lighthouse en Catálogo e Ingresar | 0 infracciones; Lighthouse ≥ 90 |
| Rendimiento | CP-18 y Lighthouse (LCP) | p95 < 300 ms; LCP < 2,5 s |
| Seguridad | Cabeceras HTTP, cookie `HttpOnly`, mensaje de login genérico, control por rol, tipo real de imagen, sin listado de `/uploads/` | Todas las pruebas pasan |
| Eficiencia | Peso del paquete JS, caché del catálogo, imágenes ≤ 2 MB | JS < 100 kB comprimido; 304 si no hay cambios |
| Adaptabilidad | Catálogo a 390 px de ancho | Sin desplazamiento horizontal |

## 6. Registro de resultados y defectos

- Cada ejecución se anota en el informe (tabla de comprobación) con fecha, resultado
  esperado y real.
- Todo defecto se registra en `docs/registro-defectos.md` y como issue en GitHub con la
  plantilla "Defecto". Se corrige en una rama `fix/DEF-XX-...`, se agrega una prueba de
  regresión y se vuelve a probar (re-prueba) antes de cerrarlo.

## 7. Cómo ejecutar

```bash
# Unitarias + integración (PostgreSQL encendido)
cd backend && go test -count=1 ./...

# Componentes y accesibilidad
cd frontend && npm test

# Funcionales en el navegador (ver pruebas/e2e/README.md)
python3 pruebas/e2e/flujo.py

# Lighthouse (Chrome DevTools > Lighthouse) sobre http://localhost:4173
```
