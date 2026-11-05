// Cliente HTTP mínimo para la API.
// Todas las peticiones envían la cookie de sesión (credentials: 'same-origin').

export class ErrorApi extends Error {
  constructor(mensaje, estado, errores = null) {
    super(mensaje)
    this.estado = estado
    this.errores = errores // errores por campo, cuando la API los devuelve
  }
}

async function pedir(ruta, { metodo = 'GET', cuerpo, formulario, signal } = {}) {
  const opciones = { method: metodo, credentials: 'same-origin', headers: {}, signal }

  if (formulario) {
    opciones.body = formulario // FormData: el navegador pone el Content-Type
  } else if (cuerpo !== undefined) {
    opciones.headers['Content-Type'] = 'application/json'
    opciones.body = JSON.stringify(cuerpo)
  }

  let resp
  try {
    resp = await fetch(ruta, opciones)
  } catch (err) {
    if (err.name === 'AbortError') throw err // petición cancelada a propósito
    throw new ErrorApi('No se pudo conectar con el servidor.', 0)
  }

  if (resp.status === 204) return null

  const datos = await resp.json().catch(() => null)
  if (!resp.ok) {
    throw new ErrorApi(datos?.error ?? 'Ocurrió un error.', resp.status, datos?.errores ?? null)
  }
  return datos
}

export const api = {
  iniciarSesion: (correo, clave) => pedir('/api/sesion', { metodo: 'POST', cuerpo: { correo, clave } }),
  sesionActual: () => pedir('/api/sesion'),
  cerrarSesion: () => pedir('/api/sesion', { metodo: 'DELETE' }),
  categorias: () => pedir('/api/categorias'),

  // HU-01
  crearProducto: (formulario) => pedir('/api/productos', { metodo: 'POST', formulario }),
  misProductos: () => pedir('/api/mis-productos'),

  // HU-02
  pendientes: () => pedir('/api/admin/pendientes'),
  aprobar: (id) => pedir(`/api/admin/productos/${id}/aprobar`, { metodo: 'POST', cuerpo: {} }),
  rechazar: (id, motivo) => pedir(`/api/admin/productos/${id}/rechazar`, { metodo: 'POST', cuerpo: { motivo } }),

  // HU-03
  catalogo: (categoriaId, signal) =>
    pedir(categoriaId ? `/api/catalogo?categoria=${encodeURIComponent(categoriaId)}` : '/api/catalogo', { signal }),
}
