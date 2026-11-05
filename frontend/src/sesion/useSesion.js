import { createContext, useContext } from 'react'

export const SesionContexto = createContext(null)

export function useSesion() {
  const ctx = useContext(SesionContexto)
  if (!ctx) throw new Error('useSesion debe usarse dentro de SesionProvider')
  return ctx
}

// Página de inicio según el rol.
export function inicioPorRol(rol) {
  if (rol === 'administradora') return '/admin/pendientes'
  if (rol === 'emprendedora') return '/mis-productos'
  return '/'
}
