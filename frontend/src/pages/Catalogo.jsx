import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { api } from '../api.js'
import Campo from '../components/Campo.jsx'
import { formatoPrecio } from '../formato.js'

// HU-03 Ver catálogo público.
// El filtro cambia la lista sin recargar la página y queda en la URL (?categoria=ID),
// así un enlace compartido abre el catálogo ya filtrado.

export default function Catalogo() {
  const [parametros, setParametros] = useSearchParams()
  const categoriaId = parametros.get('categoria') ?? ''

  const [categorias, setCategorias] = useState([])
  // Resultado de la última consulta y el filtro al que corresponde.
  const [resultado, setResultado] = useState({ clave: null, lista: [], error: '' })

  useEffect(() => {
    api.categorias().then(setCategorias).catch(() => {})
  }, [])

  useEffect(() => {
    // Si el filtro cambia antes de que llegue la respuesta, se cancela la petición vieja.
    const control = new AbortController()
    api
      .catalogo(categoriaId, control.signal)
      .then((lista) => setResultado({ clave: categoriaId, lista, error: '' }))
      .catch((e) => {
        if (!control.signal.aborted) setResultado({ clave: categoriaId, lista: [], error: e.message })
      })
    return () => control.abort()
  }, [categoriaId])

  const productos = resultado.clave === categoriaId && !resultado.error ? resultado.lista : null
  const error = resultado.clave === categoriaId ? resultado.error : ''

  const cambiarCategoria = (e) => {
    const valor = e.target.value
    setParametros(valor ? { categoria: valor } : {})
  }

  const nombreCategoria = categorias.find((c) => String(c.id) === categoriaId)?.nombre

  return (
    <>
      <h1>Catálogo de productos</h1>
      <p className="texto-suave">Productos hechos por las emprendedoras de la asociación.</p>

      <div className="filtros">
        <Campo como="select" etiqueta="Filtrar por categoría" name="categoria" value={categoriaId} onChange={cambiarCategoria}>
          <option value="">Todas las categorías</option>
          {categorias.map((c) => (
            <option key={c.id} value={c.id}>
              {c.nombre}
            </option>
          ))}
        </Campo>
      </div>

      {/* Anuncia el resultado del filtro a lectores de pantalla */}
      <p role="status" className="texto-suave" aria-live="polite">
        {productos === null
          ? 'Cargando productos…'
          : `${productos.length} ${productos.length === 1 ? 'producto' : 'productos'}${nombreCategoria ? ` en ${nombreCategoria}` : ''}`}
      </p>

      {error && (
        <p className="alerta" role="alert">
          {error}
        </p>
      )}

      {productos?.length === 0 && <p>No hay productos publicados en esta categoría por ahora.</p>}

      {productos?.length > 0 && (
        <ul className="catalogo">
          {productos.map((p) => (
            <li key={p.id} className="tarjeta producto">
              {p.imagenUrl ? (
                <img src={p.imagenUrl} alt={p.imagenAlt} loading="lazy" width="320" height="240" />
              ) : (
                <div className="sin-imagen" aria-hidden="true">
                  {p.categoria}
                </div>
              )}
              <h2>{p.nombre}</h2>
              <p className="texto-suave negocio">{p.negocio}</p>
              <p>{p.descripcion}</p>
              <p className="precio">{formatoPrecio(p.precio)}</p>
            </li>
          ))}
        </ul>
      )}
    </>
  )
}
