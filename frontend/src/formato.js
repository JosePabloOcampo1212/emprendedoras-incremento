const colones = new Intl.NumberFormat('es-CR', {
  style: 'currency',
  currency: 'CRC',
  maximumFractionDigits: 0,
})

export const formatoPrecio = (monto) => colones.format(monto)

export const ETIQUETA_ESTADO = {
  borrador: 'Devuelto para corregir',
  pendiente: 'En revisión',
  publicado: 'Publicado',
  archivado: 'Archivado',
}
