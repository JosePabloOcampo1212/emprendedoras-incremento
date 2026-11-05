# Registro de defectos

Cada defecto se corrige en su propia rama, con una prueba de regresión que falla
antes de la corrección y pasa después. Severidad: Crítica / Alta / Media / Baja.

## Resumen al cierre de la iteración

| ID | Descripción | Severidad | Encontrado en | Rama de corrección | Estado |
|---|---|---|---|---|---|
| DEF-01 | `/uploads/` listaba los archivos | Alta | Revisión de código | `feature/HU-01-registro-producto` | Cerrado |
| DEF-02 | Prueba de login fallaba por la limpieza del mock | Baja | `npm test` | `main` (estructura inicial) | Cerrado |
| DEF-03 | Salir desde página protegida abría "Ingresar" | Media | Prueba punta a punta | `fix/DEF-03-salir-redirige` | Cerrado |
| DEF-04 | Catálogo desactualizado después de aprobar | Media | Prueba punta a punta | `fix/DEF-04-cache-catalogo` | Cerrado |
| DEF-05 | Error 401 en consola en cada visita | Baja | Lighthouse | `fix/DEF-05-06-lighthouse` | Cerrado |
| DEF-06 | robots.txt inválido | Baja | Lighthouse | `fix/DEF-05-06-lighthouse` | Cerrado |

Defectos críticos abiertos: 0.

---

### DEF-01 El directorio de imágenes mostraba la lista de archivos

- **Historia:** HU-01
- **Encontrado en:** revisión de código del pull request de HU-01
- **Entorno:** backend local, Go 1.24
- **Pasos:** abrir `http://localhost:8080/uploads/`
- **Esperado:** 404, sin listar archivos
- **Real:** el servidor mostraba el nombre de todas las imágenes subidas
- **Severidad:** Alta (exposición de información)
- **Corrección:** función `sinListado` en `servidor.go`
- **Prueba de regresión:** `TestUploads_SinListado`
- **Re-prueba:** pasa
- **Estado:** Cerrado

---

### DEF-02 La prueba del login fallaba por un error no controlado

- **Historia:** estructura inicial
- **Encontrado en:** primera ejecución de `npm test`
- **Entorno:** Vitest 5, jsdom
- **Pasos:** ejecutar `npm test`
- **Esperado:** la prueba "muestra el error de la API en una alerta" pasa
- **Real:** falla con "Correo o clave incorrectos." aunque el componente sí captura el error
- **Causa:** `beforeEach(() => iniciar.mockReset())` devolvía el mock, y Vitest trata el valor devuelto como función de limpieza; al terminar la prueba lo ejecutaba y lanzaba el error
- **Severidad:** Baja (código de prueba, no afecta el producto)
- **Corrección:** usar llaves en el `beforeEach` para no devolver nada
- **Re-prueba:** pasa
- **Estado:** Cerrado

---

### DEF-03 Al cerrar sesión desde una página protegida se abría "Ingresar"

- **Historia:** inicio de sesión (afecta HU-01 y HU-02)
- **Encontrado en:** prueba funcional de punta a punta (`pruebas/e2e/flujo.py`, paso 5)
- **Entorno:** versión de producción en `localhost:4173`, Chrome
- **Pasos:** ingresar como emprendedora → "Mis productos" → "Salir"
- **Esperado:** se muestra el catálogo
- **Real:** se muestra la página "Ingresar"
- **Causa:** al borrar la sesión antes de navegar, `RutaProtegida` detecta que no hay usuario y redirige a `/ingresar` antes de que se ejecute la navegación al catálogo
- **Severidad:** Media
- **Corrección:** en `App.jsx`, navegar al catálogo primero y cerrar la sesión después
- **Prueba de regresión:** `App.test.jsx` ("desde Mis productos lleva al catálogo"); se verificó que falla sin la corrección
- **Re-prueba:** pasa; prueba punta a punta 12/12
- **Estado:** Cerrado

---

### DEF-04 El catálogo seguía mostrando datos viejos después de aprobar

- **Historia:** HU-03
- **Encontrado en:** prueba funcional de punta a punta (paso 8)
- **Entorno:** versión de producción en `localhost:4173`, Chrome
- **Pasos:** abrir el catálogo → aprobar un producto como administradora → volver al catálogo
- **Esperado:** el producto aprobado aparece
- **Real:** el navegador mostraba la lista anterior hasta por 30 segundos
- **Causa:** la respuesta del catálogo tenía `Cache-Control: public, max-age=30`
- **Severidad:** Media
- **Corrección:** `Cache-Control: no-cache` con `ETag`; si el catálogo no cambió, el servidor responde 304 sin cuerpo
- **Prueba de regresión:** `TestCatalogo_ETag`
- **Re-prueba:** pasa; prueba punta a punta 12/12
- **Estado:** Cerrado

---

### DEF-05 Error 401 en la consola en cada visita sin sesión

- **Historia:** inicio de sesión
- **Encontrado en:** Lighthouse (Buenas prácticas: 96, "errores en la consola")
- **Entorno:** versión de producción, Lighthouse 12
- **Pasos:** abrir el catálogo sin haber iniciado sesión y revisar la consola
- **Esperado:** sin errores
- **Real:** `GET /api/sesion` responde 401 y el navegador lo muestra como error
- **Severidad:** Baja
- **Corrección:** `GET /api/sesion` responde 204 cuando no hay sesión
- **Prueba de regresión:** `TestSesion_SinCookieDa204`
- **Re-prueba:** Lighthouse Buenas prácticas 100
- **Estado:** Cerrado

---

### DEF-06 robots.txt inválido

- **Historia:** HU-03 (catálogo público)
- **Encontrado en:** Lighthouse (SEO: 92)
- **Pasos:** abrir `/robots.txt`
- **Esperado:** archivo de texto con reglas válidas
- **Real:** el servidor devuelve la página HTML de la aplicación
- **Severidad:** Baja
- **Corrección:** `frontend/public/robots.txt` con reglas válidas (excluye páginas privadas)
- **Re-prueba:** Lighthouse SEO 100
- **Estado:** Cerrado
