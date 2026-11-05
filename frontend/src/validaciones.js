// Reglas de validación del lado del cliente.
// Deben coincidir con las del backend (internal/producto/producto.go).
// La validación del cliente mejora la experiencia; la del servidor es la que protege los datos.

const formatoCorreo = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

export function validarLogin({ correo, clave }) {
  const errores = {}
  if (!correo.trim()) errores.correo = 'Escriba su correo.'
  else if (!formatoCorreo.test(correo.trim())) errores.correo = 'El correo no tiene un formato válido.'
  if (!clave) errores.clave = 'Escriba su clave.'
  return errores
}

// TODO HU-01: validarProducto({ nombre, descripcion, precio, categoriaId, imagen, imagenAlt })
