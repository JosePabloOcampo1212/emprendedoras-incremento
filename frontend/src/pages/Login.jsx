import { useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import Campo from '../components/Campo.jsx'
import { inicioPorRol, useSesion } from '../sesion/useSesion.js'
import { validarLogin } from '../validaciones.js'

export default function Login() {
  const { iniciar } = useSesion()
  const navegar = useNavigate()
  const ubicacion = useLocation()

  const [datos, setDatos] = useState({ correo: '', clave: '' })
  const [errores, setErrores] = useState({})
  const [errorGeneral, setErrorGeneral] = useState('')
  const [enviando, setEnviando] = useState(false)

  const cambiar = (e) => setDatos({ ...datos, [e.target.name]: e.target.value })

  async function enviar(e) {
    e.preventDefault()
    setErrorGeneral('')
    const errs = validarLogin(datos)
    setErrores(errs)
    if (Object.keys(errs).length > 0) {
      // Lleva el foco al primer campo con error.
      document.getElementsByName(Object.keys(errs)[0])[0]?.focus()
      return
    }

    setEnviando(true)
    try {
      const u = await iniciar(datos.correo.trim(), datos.clave)
      navegar(ubicacion.state?.desde ?? inicioPorRol(u.rol), { replace: true })
    } catch (err) {
      setErrorGeneral(err.message)
    } finally {
      setEnviando(false)
    }
  }

  return (
    <section className="tarjeta angosta">
      <h1>Ingresar</h1>
      <p className="texto-suave">Acceso para emprendedoras y administración de la asociación.</p>

      {errorGeneral && (
        <p className="alerta" role="alert">
          {errorGeneral}
        </p>
      )}

      <form onSubmit={enviar} noValidate>
        <Campo
          etiqueta="Correo"
          name="correo"
          type="email"
          autoComplete="username"
          obligatorio
          value={datos.correo}
          onChange={cambiar}
          error={errores.correo}
        />
        <Campo
          etiqueta="Clave"
          name="clave"
          type="password"
          autoComplete="current-password"
          obligatorio
          value={datos.clave}
          onChange={cambiar}
          error={errores.clave}
        />
        <button type="submit" className="boton" disabled={enviando}>
          {enviando ? 'Ingresando…' : 'Ingresar'}
        </button>
      </form>
    </section>
  )
}
