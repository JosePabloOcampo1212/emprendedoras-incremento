import { useCallback, useEffect, useState } from 'react'
import { api } from '../api.js'
import Campo from '../components/Campo.jsx'
import { formatoPrecio } from '../formato.js'

// HU-02 Aprobar o rechazar productos (rol administradora).

export default function Pendientes() {
  const [productos, setProductos] = useState(null) // null = cargando
  const [mensaje, setMensaje] = useState('')
  const [errorGeneral, setErrorGeneral] = useState('')

  const cargar = useCallback(() => {
    return api
      .pendientes()
      .then(setProductos)
      .catch((e) => setErrorGeneral(e.message))
  }, [])

  useEffect(() => {
    cargar()
  }, [cargar])

  // Después de revisar, quita el producto de la lista y anuncia el resultado.
  const alRevisar = (producto, texto) => {
    setProductos((lista) => lista.filter((p) => p.id !== producto.id))
    setMensaje(texto)
  }

  const alFallar = async (err) => {
    setErrorGeneral(err.message)
    if (err.estado === 409 || err.estado === 404) await cargar()
  }

  return (
    <>
      <h1>Productos por revisar</h1>
      <p className="texto-suave">Se muestran en orden de llegada. Al aprobar, el producto se publica en el catálogo.</p>

      <p role="status" className="exito" aria-live="polite">
        {mensaje}
      </p>
      {errorGeneral && (
        <p className="alerta" role="alert">
          {errorGeneral}
        </p>
      )}

      {productos === null ? (
        <p>Cargando…</p>
      ) : productos.length === 0 ? (
        <p className="texto-suave">No hay productos pendientes de revisión.</p>
      ) : (
        <ul className="lista-productos">
          {productos.map((p) => (
            <TarjetaRevision
              key={p.id}
              producto={p}
              alRevisar={alRevisar}
              alFallar={alFallar}
              limpiar={() => {
                setMensaje('')
                setErrorGeneral('')
              }}
            />
          ))}
        </ul>
      )}
    </>
  )
}

function TarjetaRevision({ producto: p, alRevisar, alFallar, limpiar }) {
  const [rechazando, setRechazando] = useState(false)
  const [motivo, setMotivo] = useState('')
  const [errorMotivo, setErrorMotivo] = useState('')
  const [ocupado, setOcupado] = useState(false)
  const idTitulo = `producto-${p.id}`

  async function aprobar() {
    limpiar()
    setOcupado(true)
    try {
      await api.aprobar(p.id)
      alRevisar(p, `"${p.nombre}" fue aprobado y ya está en el catálogo.`)
    } catch (err) {
      alFallar(err)
      setOcupado(false)
    }
  }

  async function confirmarRechazo(e) {
    e.preventDefault()
    limpiar()
    const texto = motivo.trim()
    if (texto.length < 10) {
      setErrorMotivo(
        texto ? 'El motivo debe tener al menos 10 caracteres.' : 'Escriba el motivo para que la emprendedora sepa qué corregir.',
      )
      return
    }
    setOcupado(true)
    try {
      await api.rechazar(p.id, texto)
      alRevisar(p, `"${p.nombre}" fue devuelto a la emprendedora con el motivo indicado.`)
    } catch (err) {
      if (err.errores?.motivo) setErrorMotivo(err.errores.motivo)
      else alFallar(err)
      setOcupado(false)
    }
  }

  return (
    <li className="tarjeta revision" aria-labelledby={idTitulo}>
      {p.imagenUrl && <img src={p.imagenUrl} alt={p.imagenAlt} className="miniatura" width="120" height="120" />}
      <div>
        <h2 id={idTitulo}>{p.nombre}</h2>
        <p className="texto-suave">
          {p.negocio} · {p.categoria} · {formatoPrecio(p.precio)}
        </p>
        <p>{p.descripcion}</p>

        {!rechazando ? (
          <div className="acciones">
            <button type="button" className="boton" onClick={aprobar} disabled={ocupado}>
              Aprobar
            </button>
            <button type="button" className="boton secundario" onClick={() => setRechazando(true)} disabled={ocupado}>
              Rechazar
            </button>
          </div>
        ) : (
          <form onSubmit={confirmarRechazo} noValidate>
            <Campo
              como="textarea"
              etiqueta="Motivo del rechazo"
              name="motivo"
              rows={3}
              obligatorio
              maxLength={300}
              autoFocus
              ayuda="La emprendedora verá este mensaje. Indique qué debe corregir."
              value={motivo}
              onChange={(e) => setMotivo(e.target.value)}
              error={errorMotivo}
            />
            <div className="acciones">
              <button type="submit" className="boton peligro" disabled={ocupado}>
                Confirmar rechazo
              </button>
              <button
                type="button"
                className="boton secundario"
                onClick={() => {
                  setRechazando(false)
                  setErrorMotivo('')
                }}
                disabled={ocupado}
              >
                Cancelar
              </button>
            </div>
          </form>
        )}
      </div>
    </li>
  )
}
