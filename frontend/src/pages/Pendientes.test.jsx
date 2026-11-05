import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi } from 'vitest'
import { ErrorApi, api } from '../api.js'
import Pendientes from './Pendientes.jsx'

const lista = [
  { id: 4, nombre: 'Portavasos de fibra', negocio: 'Tejidos Montaña Verde', categoria: 'Artesanía', precio: 6000, descripcion: 'Juego de 6.', estado: 'pendiente' },
  { id: 5, nombre: 'Jabón de avena', negocio: 'Dulces de la Abuela', categoria: 'Cosmética natural', precio: 2500, descripcion: 'Barra de 100 g.', estado: 'pendiente' },
]

beforeEach(() => {
  vi.restoreAllMocks()
  vi.spyOn(api, 'pendientes').mockResolvedValue(lista)
})

const tarjeta = async (nombre) => (await screen.findByRole('heading', { name: nombre })).closest('li')

describe('Productos por revisar (HU-02)', () => {
  it('lista los pendientes', async () => {
    render(<Pendientes />)
    expect(await screen.findByRole('heading', { name: 'Portavasos de fibra' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: 'Jabón de avena' })).toBeInTheDocument()
  })

  it('CA2: aprobar quita el producto de la lista y lo anuncia', async () => {
    const aprobar = vi.spyOn(api, 'aprobar').mockResolvedValue({ id: 4, estado: 'publicado' })
    render(<Pendientes />)
    const t = await tarjeta('Portavasos de fibra')
    await userEvent.click(within(t).getByRole('button', { name: 'Aprobar' }))

    expect(aprobar).toHaveBeenCalledWith(4)
    expect(screen.queryByRole('heading', { name: 'Portavasos de fibra' })).not.toBeInTheDocument()
    expect(screen.getByRole('status')).toHaveTextContent('fue aprobado')
  })

  it('CA3: rechazar sin motivo muestra error y no llama a la API', async () => {
    const rechazar = vi.spyOn(api, 'rechazar')
    render(<Pendientes />)
    const t = await tarjeta('Jabón de avena')
    await userEvent.click(within(t).getByRole('button', { name: 'Rechazar' }))
    await userEvent.click(within(t).getByRole('button', { name: 'Confirmar rechazo' }))

    expect(within(t).getByLabelText(/Motivo del rechazo/)).toHaveAttribute('aria-invalid', 'true')
    expect(rechazar).not.toHaveBeenCalled()
  })

  it('CA3: rechazar con motivo devuelve el producto', async () => {
    const rechazar = vi.spyOn(api, 'rechazar').mockResolvedValue({ id: 5, estado: 'borrador' })
    render(<Pendientes />)
    const t = await tarjeta('Jabón de avena')
    await userEvent.click(within(t).getByRole('button', { name: 'Rechazar' }))
    expect(within(t).getByLabelText(/Motivo del rechazo/)).toHaveFocus()
    await userEvent.type(within(t).getByLabelText(/Motivo del rechazo/), 'Falta la lista de ingredientes.')
    await userEvent.click(within(t).getByRole('button', { name: 'Confirmar rechazo' }))

    expect(rechazar).toHaveBeenCalledWith(5, 'Falta la lista de ingredientes.')
    expect(screen.getByRole('status')).toHaveTextContent('fue devuelto')
  })

  it('si otra persona ya lo revisó (409), avisa y recarga la lista', async () => {
    vi.spyOn(api, 'aprobar').mockImplementation(async () => {
      throw new ErrorApi('Este producto ya fue revisado. Actualice la lista.', 409)
    })
    render(<Pendientes />)
    const t = await tarjeta('Portavasos de fibra')
    await userEvent.click(within(t).getByRole('button', { name: 'Aprobar' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('ya fue revisado')
    expect(api.pendientes).toHaveBeenCalledTimes(2)
  })

  it('sin pendientes muestra un mensaje', async () => {
    api.pendientes.mockResolvedValue([])
    render(<Pendientes />)
    expect(await screen.findByText('No hay productos pendientes de revisión.')).toBeInTheDocument()
  })
})
