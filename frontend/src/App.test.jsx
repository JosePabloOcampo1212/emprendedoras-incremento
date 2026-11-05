import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { vi } from 'vitest'
import { api } from './api.js'
import App from './App.jsx'
import SesionProvider from './sesion/SesionProvider.jsx'

// Regresión DEF-03: al cerrar sesión desde una página protegida
// la aplicación debe ir al catálogo, no a la página de ingreso.
describe('Cerrar sesión', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
    vi.spyOn(api, 'sesionActual').mockResolvedValue({ id: 2, rol: 'emprendedora', negocio: 'Dulces de la Abuela' })
    vi.spyOn(api, 'cerrarSesion').mockResolvedValue(null)
    vi.spyOn(api, 'categorias').mockResolvedValue([])
    vi.spyOn(api, 'misProductos').mockResolvedValue([])
    vi.spyOn(api, 'catalogo').mockResolvedValue([])
  })

  it('desde Mis productos lleva al catálogo', async () => {
    render(
      <MemoryRouter initialEntries={['/mis-productos']}>
        <SesionProvider>
          <App />
        </SesionProvider>
      </MemoryRouter>,
    )
    await screen.findByRole('heading', { name: 'Mis productos' })
    await userEvent.click(screen.getByRole('button', { name: 'Salir' }))

    expect(await screen.findByRole('heading', { name: 'Catálogo de productos' })).toBeInTheDocument()
    expect(screen.queryByRole('heading', { name: 'Ingresar' })).not.toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Ingresar' })).toBeInTheDocument()
  })
})
