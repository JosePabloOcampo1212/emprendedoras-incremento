# Bitácora de ejecución

Iteración única: 03/11/2026 al 06/11/2026.

| Fecha | Actividad | Estado | Dificultad / decisión |
|---|---|---|---|
| 03/11/2026 | Análisis del enunciado y selección del caso | Hecho | Se retoma la plataforma de la Asociación de Emprendedoras de la Tarea 1 |
| 03/11/2026 | Definición del incremento y de 4 historias de usuario | Hecho | Flujo elegido: registro, aprobación y publicación de productos |
| 03/11/2026 | Ajuste de alcance | Hecho | Por el plazo de 3 días, la HU-04 (archivar) pasa al backlog; el flujo central sigue completo con HU-01 a HU-03 |
| 03/11/2026 | Definición de DoR, DoD y objetivo de iteración | Hecho | Iteración única del 03/11 al 06/11 |
| 03/11/2026 | Estructura base del proyecto | Hecho | Backend Go + SQLite, frontend React + Vite, login por rol, validación de producto y primeras pruebas |
| 03/11/2026 | Cambio de base de datos: SQLite a PostgreSQL | Hecho | Cambio respecto a la línea base. PostgreSQL es más cercano a un entorno real de producción y soporta concurrencia de varias usuarias. Se agrega docker-compose para levantarlo y las pruebas de integración pasan a usar esquemas temporales |
| 04/11/2026 | Configuración de repo, issues y tablero | Hecho | Repositorio privado con plantillas de issue y PR |
| 04/11/2026 | Hipótesis y plan de medición | Hecho | Se definen seis métricas (M1 a M6) |
| 04/11/2026 | HU-01 Registrar producto | Hecho | En la revisión del PR se detecta DEF-01 (/uploads/ listaba archivos) y se corrige |
| 04/11/2026 | HU-02 Aprobar o rechazar | Hecho | Se agrega bloqueo de fila para que dos revisiones no se pisen |
| 04/11/2026 | HU-03 Catálogo público | Hecho | Se agregan pruebas de accesibilidad con axe-core |
| 04/11/2026 | Gestión de contenidos y respaldos | Hecho | Restauración probada |
| 04/11/2026 | Plan de pruebas | Hecho | 21 casos vinculados a criterios de aceptación |
| 04/11/2026 | Prueba de punta a punta n.º 1 | Hecho | Falla al cerrar sesión: DEF-03 |
| 04/11/2026 | Corrección DEF-03 | Hecho | Prueba de regresión falla antes y pasa después |
| 04/11/2026 | Prueba de punta a punta n.º 2 | Hecho | Falla al volver al catálogo: no muestra el producto aprobado (DEF-04) |
| 04/11/2026 | Corrección DEF-04 | Hecho | Caché del catálogo con ETag |
| 04/11/2026 | Prueba de punta a punta n.º 3 | Hecho | 12 de 12 pasos; M2 = 0 |
| 04/11/2026 | Lighthouse n.º 1 | Hecho | Buenas prácticas 96 y SEO 92: DEF-05 y DEF-06 |
| 04/11/2026 | Correcciones DEF-05 y DEF-06 | Hecho | Conflicto en CHANGELOG.md resuelto al integrar |
| 04/11/2026 | Lighthouse n.º 2 | Hecho | 100 en accesibilidad, buenas prácticas y SEO |
| 04/11/2026 | Prueba con 3 personas (M1) | Hecho | Mediana: 1 min 50 s |
| 04/11/2026 | Versión 1.0.0 | Hecho | Notas de versión y release publicados |

<!-- Agregar una fila por cada actividad, problema o cambio de plan. -->
