import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { vi } from 'vitest'
import { api } from '../api.js'
import Catalogo from './Catalogo.jsx'

const todos = [
  { id: 1, nombre: 'Cajeta de leche', negocio: 'Dulces de la Abuela', categoria: 'Alimentos', categoriaId: 1, precio: 3500, descripcion: 'Caja de 12.' },
  { id: 2, nombre: 'Bufanda tejida', negocio: 'Tejidos Montaña Verde', categoria: 'Textiles', categoriaId: 4, precio: 12000, descripcion: 'Lana.', imagenUrl: '/uploads/a.png', imagenAlt: 'Bufanda color tierra doblada' },
]

beforeEach(() => {
  vi.restoreAllMocks()
  vi.spyOn(api, 'categorias').mockResolvedValue([
    { id: 1, nombre: 'Alimentos' },
    { id: 4, nombre: 'Textiles' },
  ])
  vi.spyOn(api, 'catalogo').mockImplementation(async (cat) => (cat ? todos.filter((p) => String(p.categoriaId) === cat) : todos))
})

const renderizar = (ruta = '/') =>
  render(
    <MemoryRouter initialEntries={[ruta]}>
      <Catalogo />
    </MemoryRouter>,
  )

describe('Catálogo (HU-03)', () => {
  it('CA1: muestra los productos con precio en colones', async () => {
    renderizar()
    expect(await screen.findByRole('heading', { name: 'Cajeta de leche' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: 'Bufanda tejida' })).toBeInTheDocument()
    expect(screen.getByText(/3\s?500/)).toBeInTheDocument()
    expect(screen.getByRole('status')).toHaveTextContent('2 productos')
  })

  it('las imágenes tienen texto alternativo', async () => {
    renderizar()
    expect(await screen.findByRole('img', { name: 'Bufanda color tierra doblada' })).toBeInTheDocument()
  })

  it('CA2: el filtro cambia la lista sin recargar y lo anuncia', async () => {
    renderizar()
    await screen.findByRole('option', { name: 'Textiles' })
    await userEvent.selectOptions(screen.getByLabelText('Filtrar por categoría'), '4')

    expect(api.catalogo).toHaveBeenLastCalledWith('4', expect.any(AbortSignal))
    expect(await screen.findByRole('status')).toHaveTextContent('1 producto en Textiles')
    expect(screen.queryByRole('heading', { name: 'Cajeta de leche' })).not.toBeInTheDocument()
  })

  it('respeta el filtro que viene en la URL', async () => {
    renderizar('/?categoria=1')
    expect(await screen.findByRole('heading', { name: 'Cajeta de leche' })).toBeInTheDocument()
    expect(api.catalogo).toHaveBeenCalledWith('1', expect.any(AbortSignal))
  })

  it('categoría sin productos muestra un mensaje', async () => {
    api.catalogo.mockResolvedValue([])
    renderizar()
    expect(await screen.findByText(/No hay productos publicados/)).toBeInTheDocument()
  })
})
