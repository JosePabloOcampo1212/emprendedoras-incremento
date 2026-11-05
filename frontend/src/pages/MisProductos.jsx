import { useCallback, useEffect, useRef, useState } from 'react'
import { api } from '../api.js'
import Campo from '../components/Campo.jsx'
import { ETIQUETA_ESTADO, formatoPrecio } from '../formato.js'
import { useSesion } from '../sesion/useSesion.js'
import { LIMITES, validarProducto } from '../validaciones.js'

// HU-01 Registrar producto (rol emprendedora).

const VACIO = { nombre: '', descripcion: '', precio: '', categoriaId: '', imagenAlt: '' }
const ORDEN_CAMPOS = ['nombre', 'descripcion', 'precio', 'categoriaId', 'imagen', 'imagenAlt']

export default function MisProductos() {
  const { usuario } = useSesion()
  const [categorias, setCategorias] = useState([])
  const [productos, setProductos] = useState([])
  const [datos, setDatos] = useState(VACIO)
  const [imagen, setImagen] = useState(null)
  const [errores, setErrores] = useState({})
  const [mensaje, setMensaje] = useState('')
  const [errorGeneral, setErrorGeneral] = useState('')
  const [enviando, setEnviando] = useState(false)
  const formRef = useRef(null)

  const cargarProductos = useCallback(() => {
    api.misProductos().then(setProductos).catch((e) => setErrorGeneral(e.message))
  }, [])

  useEffect(() => {
    api.categorias().then(setCategorias).catch((e) => setErrorGeneral(e.message))
    cargarProductos()
  }, [cargarProductos])

  const cambiar = (e) => setDatos({ ...datos, [e.target.name]: e.target.value })
  const elegirImagen = (e) => setImagen(e.target.files?.[0] ?? null)

  function enfocarPrimerError(errs) {
    const primero = ORDEN_CAMPOS.find((c) => errs[c])
    if (primero) formRef.current?.elements.namedItem(primero)?.focus()
  }

  async function enviar(e) {
    e.preventDefault()
    setMensaje('')
    setErrorGeneral('')

    const errs = validarProducto({ ...datos, imagen })
    setErrores(errs)
    if (Object.keys(errs).length > 0) {
      enfocarPrimerError(errs)
      return
    }

    const fd = new FormData()
    Object.entries(datos).forEach(([k, v]) => fd.append(k, v.trim()))
    if (imagen) fd.append('imagen', imagen)

    setEnviando(true)
    try {
      await api.crearProducto(fd)
      setDatos(VACIO)
      setImagen(null)
      formRef.current?.reset()
      setMensaje('Producto enviado a revisión. Aparecerá en el catálogo cuando la administración lo apruebe.')
      cargarProductos()
    } catch (err) {
      if (err.errores) {
        setErrores(err.errores)
        enfocarPrimerError(err.errores)
      }
      setErrorGeneral(err.message)
    } finally {
      setEnviando(false)
    }
  }

  return (
    <>
      <h1>Mis productos</h1>
      {usuario?.negocio && <p className="texto-suave">{usuario.negocio}</p>}

      <div className="dos-columnas">
        <section className="tarjeta" aria-labelledby="titulo-registro">
          <h2 id="titulo-registro">Registrar un producto</h2>
          <p className="texto-suave">Los campos con * son obligatorios.</p>

          {errorGeneral && (
            <p className="alerta" role="alert">
              {errorGeneral}
            </p>
          )}

          <form ref={formRef} onSubmit={enviar} noValidate>
            <Campo
              etiqueta="Nombre del producto"
              name="nombre"
              obligatorio
              maxLength={LIMITES.nombreMax}
              value={datos.nombre}
              onChange={cambiar}
              error={errores.nombre}
            />
            <Campo
              como="textarea"
              etiqueta="Descripción"
              name="descripcion"
              rows={4}
              obligatorio
              maxLength={LIMITES.descripcionMax}
              ayuda="Entre 10 y 500 caracteres. Cuente de qué está hecho, tamaño o presentación."
              value={datos.descripcion}
              onChange={cambiar}
              error={errores.descripcion}
            />
            <Campo
              etiqueta="Precio en colones"
              name="precio"
              type="text"
              inputMode="numeric"
              obligatorio
              ayuda="Solo números, sin puntos ni símbolos. Ejemplo: 4500"
              value={datos.precio}
              onChange={cambiar}
              error={errores.precio}
            />
            <Campo
              como="select"
              etiqueta="Categoría"
              name="categoriaId"
              obligatorio
              value={datos.categoriaId}
              onChange={cambiar}
              error={errores.categoriaId}
            >
              <option value="">Seleccione…</option>
              {categorias.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.nombre}
                </option>
              ))}
            </Campo>
            <Campo
              etiqueta="Imagen (opcional)"
              name="imagen"
              type="file"
              accept="image/jpeg,image/png"
              ayuda="JPG o PNG, máximo 2 MB."
              onChange={elegirImagen}
              error={errores.imagen}
            />
            {imagen && (
              <Campo
                etiqueta="Descripción de la imagen"
                name="imagenAlt"
                obligatorio
                maxLength={LIMITES.altMax}
                ayuda="Para personas que usan lector de pantalla. Ejemplo: Frasco de mermelada de mora con tapa de tela."
                value={datos.imagenAlt}
                onChange={cambiar}
                error={errores.imagenAlt}
              />
            )}
            <button type="submit" className="boton" disabled={enviando}>
              {enviando ? 'Enviando…' : 'Enviar a revisión'}
            </button>
          </form>

          <p role="status" className="exito" aria-live="polite">
            {mensaje}
          </p>
        </section>

        <section aria-labelledby="titulo-lista">
          <h2 id="titulo-lista">Registrados</h2>
          {productos.length === 0 ? (
            <p className="texto-suave">Todavía no hay productos registrados.</p>
          ) : (
            <ul className="lista-productos">
              {productos.map((p) => (
                <li key={p.id} className="tarjeta compacta">
                  <div className="fila-producto">
                    <strong>{p.nombre}</strong>
                    <span className={`estado estado-${p.estado}`}>{ETIQUETA_ESTADO[p.estado]}</span>
                  </div>
                  <span className="texto-suave">
                    {p.categoria} · {formatoPrecio(p.precio)}
                  </span>
                  {p.estado === 'borrador' && p.motivoRechazo && (
                    <p className="motivo">
                      <strong>Motivo de la devolución:</strong> {p.motivoRechazo}
                    </p>
                  )}
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>
    </>
  )
}
