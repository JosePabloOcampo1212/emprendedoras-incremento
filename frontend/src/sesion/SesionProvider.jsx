import { useCallback, useEffect, useState } from 'react'
import { api } from '../api.js'
import { SesionContexto } from './useSesion.js'

export default function SesionProvider({ children }) {
  const [usuario, setUsuario] = useState(null)
  const [cargando, setCargando] = useState(true)

  // Al abrir la app, revisa si ya hay una sesión vigente.
  useEffect(() => {
    api
      .sesionActual()
      .then(setUsuario)
      .catch(() => setUsuario(null))
      .finally(() => setCargando(false))
  }, [])

  const iniciar = useCallback(async (correo, clave) => {
    const u = await api.iniciarSesion(correo, clave)
    setUsuario(u)
    return u
  }, [])

  const cerrar = useCallback(async () => {
    await api.cerrarSesion().catch(() => {})
    setUsuario(null)
  }, [])

  return <SesionContexto.Provider value={{ usuario, cargando, iniciar, cerrar }}>{children}</SesionContexto.Provider>
}
