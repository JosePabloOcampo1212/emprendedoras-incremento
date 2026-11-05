import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi } from 'vitest'
import { ErrorApi, api } from '../api.js'
import MisProductos from './MisProductos.jsx'

vi.mock('../sesion/useSesion.js', async (original) => ({
  ...(await original()),
  useSesion: () => ({ usuario: { rol: 'emprendedora', negocio: 'Dulces de la Abuela' } }),
}))

const productoDevuelto = {
  id: 7,
  nombre: 'Cajeta de coco',
  categoria: 'Alimentos',
  precio: 3000,
  estado: 'borrador',
  motivoRechazo: 'Falta indicar el tamaño de la caja.',
}

beforeEach(() => {
  vi.restoreAllMocks()
  vi.spyOn(api, 'categorias').mockResolvedValue([
    { id: 1, nombre: 'Alimentos' },
    { id: 2, nombre: 'Textiles' },
  ])
  vi.spyOn(api, 'misProductos').mockResolvedValue([productoDevuelto])
})

async function llenarValido() {
  await userEvent.type(screen.getByLabelText(/Nombre del producto/), 'Mermelada de mora')
  await userEvent.type(screen.getByLabelText(/^Descripción/), 'Frasco de 250 g sin preservantes.')
  await userEvent.type(screen.getByLabelText(/Precio/), '2800')
  await userEvent.selectOptions(screen.getByLabelText(/Categoría/), '1')
}

describe('Mis productos (HU-01)', () => {
  it('muestra los productos y el motivo de devolución', async () => {
    render(<MisProductos />)
    const item = await screen.findByText('Cajeta de coco')
    const tarjeta = item.closest('li')
    expect(within(tarjeta).getByText('Devuelto para corregir')).toBeInTheDocument()
    expect(within(tarjeta).getByText(/Falta indicar el tamaño/)).toBeInTheDocument()
  })

  it('CA1: con campos vacíos muestra errores junto a cada campo y no envía', async () => {
    const crear = vi.spyOn(api, 'crearProducto')
    render(<MisProductos />)
    await userEvent.click(screen.getByRole('button', { name: 'Enviar a revisión' }))

    expect(screen.getByLabelText(/Nombre del producto/)).toHaveAttribute('aria-invalid', 'true')
    expect(screen.getByLabelText(/Precio/)).toHaveAccessibleDescription(/El precio debe ser un número mayor que cero/)
    expect(screen.getByLabelText(/Nombre del producto/)).toHaveFocus()
    expect(crear).not.toHaveBeenCalled()
  })

  it('CA2: con datos válidos envía el formulario y confirma', async () => {
    const crear = vi.spyOn(api, 'crearProducto').mockResolvedValue({ id: 10, estado: 'pendiente' })
    render(<MisProductos />)
    await screen.findByRole('option', { name: 'Alimentos' })
    await llenarValido()
    await userEvent.click(screen.getByRole('button', { name: 'Enviar a revisión' }))

    expect(crear).toHaveBeenCalledTimes(1)
    const fd = crear.mock.calls[0][0]
    expect(fd.get('nombre')).toBe('Mermelada de mora')
    expect(fd.get('precio')).toBe('2800')
    expect(fd.get('categoriaId')).toBe('1')
    expect(await screen.findByRole('status')).toHaveTextContent('Producto enviado a revisión')
  })

  it('CA3: al elegir imagen pide el texto alternativo', async () => {
    render(<MisProductos />)
    expect(screen.queryByLabelText(/Descripción de la imagen/)).not.toBeInTheDocument()

    const png = new File(['x'], 'foto.png', { type: 'image/png' })
    await userEvent.upload(screen.getByLabelText(/Imagen/), png)
    expect(screen.getByLabelText(/Descripción de la imagen/)).toBeRequired()
  })

  it('muestra los errores por campo que devuelve la API', async () => {
    vi.spyOn(api, 'crearProducto').mockImplementation(async () => {
      throw new ErrorApi('Revise los campos marcados.', 422, { categoriaId: 'La categoría seleccionada no existe.' })
    })
    render(<MisProductos />)
    await screen.findByRole('option', { name: 'Alimentos' })
    await llenarValido()
    await userEvent.click(screen.getByRole('button', { name: 'Enviar a revisión' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('Revise los campos marcados.')
    expect(screen.getByLabelText(/Categoría/)).toHaveAccessibleDescription('La categoría seleccionada no existe.')
  })
})
