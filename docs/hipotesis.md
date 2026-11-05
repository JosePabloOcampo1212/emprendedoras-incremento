# Hipótesis y diseño de validación

## Hipótesis

Si el registro de productos valida cada campo en el navegador y explica el error
junto al campo, entonces una emprendedora podrá publicar un producto completo
(con imagen y texto alternativo) en menos de 3 minutos y sin que la interfaz envíe
al servidor datos inválidos. Además, el flujo registro → revisión → catálogo
funcionará sin errores críticos y con un nivel de accesibilidad AA.

## Métricas y criterios de éxito

| ID | Métrica | Cómo se mide | Criterio de éxito |
|---|---|---|---|
| M1 | Tiempo para registrar un producto completo | Prueba con 3 personas usando datos sintéticos; se cronometra desde que abren "Mis productos" hasta el mensaje de confirmación | Mediana menor a 3 minutos |
| M2 | Rechazos del servidor por validación (HTTP 422) originados desde la interfaz | Conteo en el log del backend durante las pruebas funcionales | 0 |
| M3 | Criterios de aceptación cumplidos | Casos de prueba del plan de pruebas | 100 % (10 de 10) |
| M4 | Accesibilidad | Lighthouse (categoría Accesibilidad) y axe-core | Lighthouse ≥ 90 y 0 infracciones en axe-core |
| M5 | Tiempo de respuesta del catálogo | Prueba automática con 500 productos publicados, 30 consultas | Percentil 95 menor a 300 ms |
| M6 | Defectos críticos abiertos al cierre | Registro de defectos | 0 |

## Regla de decisión

- La hipótesis se **acepta** si se cumplen M1, M2 y M3, y M4 a M6 no muestran problemas.
- Se **refuta** si M1 da 3 minutos o más, o si M2 es mayor que 0. En ese caso se
  registran las causas observadas y se proponen cambios para la siguiente iteración.

## Plan de medición

| Momento | Actividad | Métricas |
|---|---|---|
| Al cerrar cada historia | Correr pruebas unitarias e integración | M3 (parcial) |
| Jueves 05/11 | Pruebas funcionales de punta a punta y revisión del log | M2, M3 |
| Jueves 05/11 | Lighthouse y axe-core sobre la versión de producción | M4 |
| Jueves 05/11 | Prueba de rendimiento del catálogo | M5 |
| Jueves 05/11 (noche) | Prueba con 3 personas, cronometrada | M1 |
| Viernes 06/11 | Revisión del registro de defectos | M6 |
