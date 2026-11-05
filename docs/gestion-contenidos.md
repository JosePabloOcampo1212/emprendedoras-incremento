# Gestión de contenidos

Plataforma de la Asociación de Emprendedoras. Este documento define qué contenidos
maneja el sitio, quién se encarga de cada uno, qué permisos tiene cada rol y cómo
se publican, actualizan, archivan y respaldan.

## 1. Roles

| Rol | Quién es | Acceso |
|---|---|---|
| Visitante | Cualquier persona en internet | Solo lectura del catálogo público |
| Emprendedora | Socia de la asociación con cuenta | Registra sus productos y ve el estado de cada uno |
| Administradora | Persona designada por la junta directiva | Revisa, aprueba o devuelve productos; administra categorías y cuentas |

Los permisos se controlan en la API (`requiereRol`). Ocultar un botón en la
interfaz no se considera control de acceso.

## 2. Matriz de contenidos

| Tipo de contenido | Responsable de crear | Aprueba / publica | Permisos | Ciclo de vida | Actualización | Archivo y retención | Respaldo |
|---|---|---|---|---|---|---|---|
| Producto (nombre, descripción, precio, categoría) | Emprendedora dueña | Administradora | Emprendedora: crear y ver los suyos. Administradora: ver todos, aprobar, devolver. Visitante: ver publicados | borrador → pendiente → publicado → archivado | La emprendedora revisa sus precios una vez al mes | Archivado cuando deja de venderse; se conserva para historial (HU-04, en backlog) | Diario, base de datos |
| Imagen de producto | Emprendedora dueña | Se publica junto con el producto | Igual que el producto | Igual que el producto | Se reemplaza al cambiar la presentación | Se conserva mientras exista el producto | Diario, carpeta de imágenes |
| Texto alternativo de la imagen | Emprendedora dueña | Administradora lo revisa al aprobar | Igual que el producto | Igual que el producto | Junto con la imagen | Igual que la imagen | Diario, base de datos |
| Categoría | Administradora | Administradora | Solo administradora (por ahora vía seed/BD) | activa | Cuando la asociación agrega una línea de productos | No se borra si tiene productos | Diario, base de datos |
| Motivo de devolución | Administradora | No aplica (mensaje interno) | Lo ve la emprendedora dueña y la administración | Se guarda al devolver; se limpia al aprobar | No se edita | En el historial de estados | Diario, base de datos |
| Historial de estados | El sistema | No aplica | Solo administración (consulta en BD) | Solo se agrega, nunca se edita | No aplica | Se conserva completo para auditoría | Diario, base de datos |
| Cuenta de usuaria | Administradora | Administradora | La persona ve su propio perfil | activa / inactiva | Cuando cambia el contacto de la socia | Se desactiva al salir de la asociación | Diario, base de datos |

## 3. Procedimiento de publicación

1. **Registro.** La emprendedora completa el formulario de "Mis productos". El sistema
   valida los datos en el navegador y otra vez en el servidor. Si todo está bien, el
   producto queda en estado *pendiente* y se registra en el historial.
2. **Revisión.** La administradora abre "Por revisar", donde los productos aparecen en
   orden de llegada. Revisa con estos criterios editoriales:
   - el nombre y la descripción son claros y describen el producto real;
   - el precio está en colones y es razonable para el producto;
   - la imagen es propia de la emprendedora (no tomada de internet) y no muestra rostros
     ni datos personales;
   - el texto alternativo describe la imagen;
   - el lenguaje es respetuoso e inclusivo, sin estereotipos.
3. **Decisión.**
   - *Aprobar:* el producto pasa a *publicado* y aparece en el catálogo de inmediato.
   - *Devolver:* la administradora escribe un motivo (mínimo 10 caracteres) y el producto
     vuelve a *borrador*. La emprendedora ve el motivo en "Mis productos".
4. **Tiempo de respuesta.** Meta: revisar cada producto pendiente en un máximo de 2 días hábiles.
5. **Trazabilidad.** Cada cambio de estado queda en `historial_estados` con fecha, persona
   que lo hizo y nota.

## 4. Actualización y archivo

- La edición de productos y el archivo (HU-04) quedaron fuera de esta iteración. Mientras
  tanto, si una emprendedora necesita retirar o cambiar un producto, lo solicita a la
  administración, que lo gestiona directamente en la base de datos y deja nota en el historial.
- Una vez al mes la administración revisa el catálogo y consulta a las socias por productos
  que ya no venden.
- El archivo no borra datos: el producto pasa a *archivado*, deja de verse en el catálogo
  y se conserva para el historial de la asociación.

## 5. Respaldo y restauración

| Qué | Cómo | Frecuencia | Retención |
|---|---|---|---|
| Base de datos | `scripts/respaldo.sh` (usa `pg_dump` en formato comprimido) | Diaria | Últimos 7 respaldos |
| Imágenes subidas | El mismo script empaqueta `backend/data/uploads` | Diaria | Últimos 7 respaldos |
| Código fuente | Repositorio privado en GitHub | Cada commit | Historial completo |

- Los respaldos se guardan en `respaldos/`, que está fuera del control de versiones,
  porque contienen datos de la asociación.
- Para restaurar: `scripts/restaurar.sh respaldos/<archivo>.dump`.
- Una vez al mes se prueba restaurar el último respaldo en una base aparte, para
  confirmar que sirve.

## 6. Ejes transversales

**Accesibilidad.** Formularios con etiquetas asociadas, errores anunciados junto a cada
campo, foco visible, navegación completa con teclado, texto alternativo obligatorio en
imágenes y contraste AA. Meta verificable: Lighthouse accesibilidad ≥ 90 (resultado
obtenido: 100).

**Privacidad.** Solo se piden los datos necesarios para publicar un producto. El catálogo
público no muestra correos ni datos de contacto personales, solo el nombre del negocio.
Las imágenes se guardan con nombre aleatorio. Todos los datos de prueba son sintéticos.
Se toma como referencia la Ley 8968 de Protección de la Persona frente al Tratamiento de
sus Datos Personales.

**Inclusión y diversidad cultural.** La interfaz usa lenguaje inclusivo ("emprendedoras",
"la administración", "personas que usan lector de pantalla"), precios en colones con
formato local (es-CR) y categorías que reflejan productos de la zona.

**Sostenibilidad y eficiencia.** Imágenes limitadas a 2 MB, carga diferida de imágenes en
el catálogo (`loading="lazy"`), catálogo con validación de caché por ETag (si nada cambió,
el servidor responde 304 sin reenviar datos) y un índice en la base de datos para la
consulta del catálogo. Esto reduce datos móviles para visitantes y carga del servidor.
