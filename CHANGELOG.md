# Notas de versión

Formato basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.1.0/).
Versionado según [Semantic Versioning](https://semver.org/lang/es/).

## [Sin publicar]

## [1.0.0] - 2026-11-05

Primera versión del incremento: flujo completo de publicación de productos.

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

### Pendiente para la siguiente versión

- HU-04: archivar y restaurar productos.
- Edición de productos devueltos por la administración.
- Paginación del catálogo (hoy muestra hasta 200 productos).

[Sin publicar]: https://github.com/JosePabloOcampo1212/emprendedoras-incremento/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/JosePabloOcampo1212/emprendedoras-incremento/releases/tag/v1.0.0
