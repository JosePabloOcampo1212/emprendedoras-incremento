// HU-02 Aprobar o rechazar productos (rol administradora).
// TODO HU-02:
//  - listar api.pendientes()
//  - botón Aprobar -> api.aprobar(id)
//  - botón Rechazar -> pide motivo obligatorio -> api.rechazar(id, motivo)
//  - anunciar el resultado con un mensaje role="status"

export default function Pendientes() {
  return (
    <section>
      <h1>Productos por revisar</h1>
      <p className="pendiente">La revisión de productos se completa en la HU-02.</p>
    </section>
  )
}
