import { useSesion } from '../sesion/useSesion.js'

// HU-01 Registrar producto (rol emprendedora).
// TODO HU-01:
//  - formulario con Campo: nombre, descripción, precio, categoría, imagen y texto alternativo
//  - validar en el cliente con las mismas reglas que producto.Validar en el backend
//  - enviar con FormData a api.crearProducto y mostrar los errores por campo que devuelva la API
//  - listar los productos propios con su estado (pendiente, publicado, borrador + motivo de rechazo)

export default function MisProductos() {
  const { usuario } = useSesion()
  return (
    <section>
      <h1>Mis productos</h1>
      {usuario?.negocio && <p className="texto-suave">{usuario.negocio}</p>}
      <p className="pendiente">El formulario de registro se completa en la HU-01.</p>
    </section>
  )
}
