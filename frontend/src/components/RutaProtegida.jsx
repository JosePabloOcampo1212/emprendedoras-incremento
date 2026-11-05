import { Navigate, useLocation } from 'react-router-dom'
import { useSesion } from '../sesion/useSesion.js'

// Muestra la página solo si hay sesión y el rol coincide.
// Esto es comodidad de interfaz: el control real de permisos está en la API.
export default function RutaProtegida({ rol, children }) {
  const { usuario, cargando } = useSesion()
  const ubicacion = useLocation()

  if (cargando) return <p role="status">Cargando…</p>
  if (!usuario) return <Navigate to="/ingresar" replace state={{ desde: ubicacion.pathname }} />
  if (rol && usuario.rol !== rol) {
    return (
      <section>
        <h1>Sin acceso</h1>
        <p>Su cuenta no tiene permiso para ver esta página.</p>
      </section>
    )
  }
  return children
}
