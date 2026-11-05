import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { vi } from 'vitest'
import { ErrorApi } from '../api.js'
import { validarLogin } from '../validaciones.js'
import Login from './Login.jsx'

// Se simula el contexto de sesión para probar la página sin la API real.
const iniciar = vi.fn()
vi.mock('../sesion/useSesion.js', async (original) => ({
  ...(await original()),
  useSesion: () => ({ iniciar }),
}))

function renderizar() {
  render(
    <MemoryRouter initialEntries={['/ingresar']}>
      <Routes>
        <Route path="/ingresar" element={<Login />} />
        <Route path="/mis-productos" element={<h1>Mis productos</h1>} />
        <Route path="/admin/pendientes" element={<h1>Por revisar</h1>} />
      </Routes>
    </MemoryRouter>,
  )
}

describe('validarLogin', () => {
  it('pide correo y clave', () => {
    expect(validarLogin({ correo: '', clave: '' })).toEqual({
      correo: 'Escriba su correo.',
      clave: 'Escriba su clave.',
    })
  })

  it('rechaza un correo mal escrito', () => {
    expect(validarLogin({ correo: 'sin-arroba', clave: 'x' }).correo).toMatch(/formato/)
  })

  it('acepta datos completos', () => {
    expect(validarLogin({ correo: 'ana@asociacion.test', clave: 'x' })).toEqual({})
  })
})

describe('Página de ingreso', () => {
  beforeEach(() => {
    iniciar.mockReset()
  })

  it('muestra errores y no llama a la API si faltan datos', async () => {
    renderizar()
    await userEvent.click(screen.getByRole('button', { name: 'Ingresar' }))

    expect(screen.getByLabelText(/Correo/)).toHaveAttribute('aria-invalid', 'true')
    expect(screen.getByText('Escriba su clave.')).toBeInTheDocument()
    expect(iniciar).not.toHaveBeenCalled()
  })

  it('lleva el foco al primer campo con error', async () => {
    renderizar()
    await userEvent.click(screen.getByRole('button', { name: 'Ingresar' }))
    expect(screen.getByLabelText(/Correo/)).toHaveFocus()
  })

  it('con credenciales válidas redirige según el rol', async () => {
    iniciar.mockResolvedValue({ rol: 'emprendedora' })
    renderizar()

    await userEvent.type(screen.getByLabelText(/Correo/), 'emprendedora1@asociacion.test')
    await userEvent.type(screen.getByLabelText(/Clave/), 'una-clave')
    await userEvent.click(screen.getByRole('button', { name: 'Ingresar' }))

    expect(iniciar).toHaveBeenCalledWith('emprendedora1@asociacion.test', 'una-clave')
    expect(await screen.findByRole('heading', { name: 'Mis productos' })).toBeInTheDocument()
  })

  it('muestra el error de la API en una alerta', async () => {
    iniciar.mockImplementation(async () => {
      throw new ErrorApi('Correo o clave incorrectos.', 401)
    })
    renderizar()

    await userEvent.type(screen.getByLabelText(/Correo/), 'admin@asociacion.test')
    await userEvent.type(screen.getByLabelText(/Clave/), 'mala')
    await userEvent.click(screen.getByRole('button', { name: 'Ingresar' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('Correo o clave incorrectos.')
  })
})
