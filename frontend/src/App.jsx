import { NavLink, Route, Routes, useNavigate } from 'react-router-dom'
import RutaProtegida from './components/RutaProtegida.jsx'
import Catalogo from './pages/Catalogo.jsx'
import Login from './pages/Login.jsx'
import MisProductos from './pages/MisProductos.jsx'
import Pendientes from './pages/Pendientes.jsx'
import { useSesion } from './sesion/useSesion.js'

export default function App() {
  const { usuario, cerrar } = useSesion()
  const navegar = useNavigate()

  // Primero se sale de la página protegida y luego se cierra la sesión.
  // Al revés, RutaProtegida detecta que ya no hay usuario y manda a /ingresar (DEF-03).
  async function salir() {
    navegar('/')
    await cerrar()
  }

  return (
    <>
      <a href="#contenido" className="saltar">
        Saltar al contenido
      </a>

      <header className="encabezado">
        <div className="contenedor fila">
          <NavLink to="/" className="marca">
            Asociación de Emprendedoras
          </NavLink>
          <nav aria-label="Principal">
            <ul className="menu">
              <li>
                <NavLink to="/" end>
                  Catálogo
                </NavLink>
              </li>
              {usuario?.rol === 'emprendedora' && (
                <li>
                  <NavLink to="/mis-productos">Mis productos</NavLink>
                </li>
              )}
              {usuario?.rol === 'administradora' && (
                <li>
                  <NavLink to="/admin/pendientes">Por revisar</NavLink>
                </li>
              )}
              <li>
                {usuario ? (
                  <button type="button" className="enlace" onClick={salir}>
                    Salir
                  </button>
                ) : (
                  <NavLink to="/ingresar">Ingresar</NavLink>
                )}
              </li>
            </ul>
          </nav>
        </div>
      </header>

      <main id="contenido" className="contenedor" tabIndex={-1}>
        <Routes>
          <Route path="/" element={<Catalogo />} />
          <Route path="/ingresar" element={<Login />} />
          <Route
            path="/mis-productos"
            element={
              <RutaProtegida rol="emprendedora">
                <MisProductos />
              </RutaProtegida>
            }
          />
          <Route
            path="/admin/pendientes"
            element={
              <RutaProtegida rol="administradora">
                <Pendientes />
              </RutaProtegida>
            }
          />
          <Route
            path="*"
            element={
              <section>
                <h1>Página no encontrada</h1>
              </section>
            }
          />
        </Routes>
      </main>

      <footer className="pie">
        <div className="contenedor">Prototipo académico. Todos los datos son ficticios.</div>
      </footer>
    </>
  )
}
