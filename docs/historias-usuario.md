# Historias de usuario

Flujo del incremento: registro → revisión → publicación de productos.

## HU-01 Registrar producto

Como emprendedora, quiero registrar un producto con nombre, descripción, precio,
categoría e imagen para darlo a conocer en la plataforma.

**Criterios de aceptación**

1. Si un campo obligatorio está vacío o el precio no es un número positivo, el formulario muestra el error junto al campo y no guarda nada.
2. Si los datos son válidos, el producto queda guardado con estado "pendiente".
3. La imagen solo acepta JPG o PNG de hasta 2 MB y pide texto alternativo.
4. El formulario se puede completar usando solo el teclado.

**Reglas de validación**

| Campo | Regla |
|---|---|
| Nombre | obligatorio, de 3 a 80 caracteres |
| Descripción | obligatoria, de 10 a 500 caracteres |
| Precio | entero en colones, mayor que 0 y hasta 10 000 000 |
| Categoría | obligatoria, de la lista |
| Imagen | opcional; JPG o PNG (se verifica el contenido, no la extensión), máximo 2 MB |
| Texto alternativo | obligatorio si hay imagen, máximo 150 caracteres |

## HU-02 Aprobar o rechazar producto

Como administradora de la asociación, quiero revisar los productos pendientes
para decidir si se publican.

**Criterios de aceptación**

1. Solo el rol administradora ve la lista de pendientes.
2. Al aprobar, el estado pasa a "publicado".
3. Al rechazar, el motivo es obligatorio y el producto vuelve a "borrador" con ese motivo visible para la emprendedora.

## HU-03 Ver catálogo público

Como visitante, quiero ver los productos publicados y filtrarlos por categoría
para encontrar lo que me interesa.

**Criterios de aceptación**

1. Solo aparecen productos con estado "publicado".
2. El filtro por categoría funciona sin recargar la página.
3. La página obtiene 90 o más en accesibilidad en Lighthouse.

## HU-04 Archivar producto (backlog)

Como emprendedora, quiero archivar un producto que ya no vendo para que deje de
verse sin perder su información.

**Criterios de aceptación**

1. Al archivarlo, desaparece del catálogo pero sigue en la base de datos.
2. La emprendedora puede restaurarlo, y en ese caso vuelve a "pendiente".

> Fuera del alcance de la iteración por el plazo de entrega (decisión del 03/11/2026).
> El ciclo de estados ya contempla el archivo (`producto.PuedeCambiar`).
