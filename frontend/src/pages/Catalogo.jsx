// HU-03 Ver catálogo público.
// TODO HU-03:
//  - pedir api.catalogo(categoriaId) y mostrar solo productos publicados
//  - filtro por categoría sin recargar la página (select con api.categorias())
//  - cada tarjeta: imagen con su texto alternativo, nombre, negocio, precio en colones
//  - meta de accesibilidad: Lighthouse >= 90

export default function Catalogo() {
  return (
    <section>
      <h1>Catálogo de productos</h1>
      <p className="texto-suave">Productos hechos por las emprendedoras de la asociación.</p>
      <p className="pendiente">Esta sección se completa en la HU-03.</p>
    </section>
  )
}
