# Notas de versión

Formato basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.1.0/).
Versionado según [Semantic Versioning](https://semver.org/lang/es/).

## [Sin publicar]

### Agregado

- HU-01: registro de productos con validación en cliente y servidor, imagen opcional (JPG/PNG, máx. 2 MB) con texto alternativo obligatorio.
- HU-02: revisión de productos por la administradora: aprobar o devolver con motivo; historial de cambios de estado.
- HU-03: catálogo público con filtro por categoría sin recargar la página.
- Inicio de sesión por rol (emprendedora, administradora) con cookie `HttpOnly`.
- Gestión de contenidos: matriz, procedimiento de publicación y scripts de respaldo y restauración.

### Corregido

- Al cerrar sesión desde una página protegida ahora se muestra el catálogo (DEF-03).
- El catálogo muestra los productos recién aprobados sin esperar a que venza la caché; usa validación con ETag (DEF-04).
- La consulta de sesión responde 204 cuando no hay sesión, sin errores en la consola (DEF-05).
- Se agrega `robots.txt` válido (DEF-06).

### Seguridad

- El directorio de imágenes no lista archivos (DEF-01).
