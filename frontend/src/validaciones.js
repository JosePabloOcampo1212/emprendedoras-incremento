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

export const LIMITES = {
  nombreMin: 3,
  nombreMax: 80,
  descripcionMin: 10,
  descripcionMax: 500,
  precioMax: 10_000_000,
  altMax: 150,
  imagenMaxBytes: 2 * 1024 * 1024,
}

const TIPOS_IMAGEN = ['image/jpeg', 'image/png']

// Cuenta caracteres reales (una letra con tilde cuenta como 1).
const largo = (texto) => [...texto.trim()].length

export function validarProducto({ nombre, descripcion, precio, categoriaId, imagen, imagenAlt }) {
  const e = {}

  const n = largo(nombre)
  if (n === 0) e.nombre = 'El nombre es obligatorio.'
  else if (n < LIMITES.nombreMin) e.nombre = 'El nombre debe tener al menos 3 caracteres.'
  else if (n > LIMITES.nombreMax) e.nombre = 'El nombre no puede pasar de 80 caracteres.'

  const d = largo(descripcion)
  if (d === 0) e.descripcion = 'La descripción es obligatoria.'
  else if (d < LIMITES.descripcionMin) e.descripcion = 'La descripción debe tener al menos 10 caracteres.'
  else if (d > LIMITES.descripcionMax) e.descripcion = 'La descripción no puede pasar de 500 caracteres.'

  const p = String(precio).trim()
  if (!/^\d+$/.test(p) || Number(p) <= 0) e.precio = 'El precio debe ser un número mayor que cero.'
  else if (Number(p) > LIMITES.precioMax) e.precio = 'El precio no puede ser mayor a ₡10 000 000.'

  if (!categoriaId) e.categoriaId = 'Seleccione una categoría.'

  if (imagen) {
    if (!TIPOS_IMAGEN.includes(imagen.type)) e.imagen = 'La imagen debe ser JPG o PNG.'
    else if (imagen.size > LIMITES.imagenMaxBytes) e.imagen = 'La imagen no puede pesar más de 2 MB.'

    const a = largo(imagenAlt)
    if (a === 0) e.imagenAlt = 'Describa la imagen para personas que usan lector de pantalla.'
    else if (a > LIMITES.altMax) e.imagenAlt = 'La descripción de la imagen no puede pasar de 150 caracteres.'
  }

  return e
}
