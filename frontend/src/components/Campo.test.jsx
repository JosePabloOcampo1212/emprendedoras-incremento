import { render, screen } from '@testing-library/react'
import Campo from './Campo.jsx'

describe('Campo', () => {
  it('asocia la etiqueta con el control', () => {
    render(<Campo etiqueta="Nombre" name="nombre" />)
    expect(screen.getByLabelText('Nombre')).toBeInTheDocument()
  })

  it('sin error no marca el campo como inválido', () => {
    render(<Campo etiqueta="Nombre" name="nombre" />)
    expect(screen.getByLabelText('Nombre')).not.toHaveAttribute('aria-invalid')
  })

  it('con error lo anuncia y lo vincula al control', () => {
    render(<Campo etiqueta="Precio" name="precio" error="El precio debe ser mayor que cero." />)
    const control = screen.getByLabelText('Precio')
    expect(control).toHaveAttribute('aria-invalid', 'true')
    expect(control).toHaveAccessibleDescription('El precio debe ser mayor que cero.')
  })

  it('vincula la ayuda y el error juntos', () => {
    render(<Campo etiqueta="Imagen" name="imagen" ayuda="JPG o PNG, máximo 2 MB." error="Formato no permitido." />)
    expect(screen.getByLabelText('Imagen')).toHaveAccessibleDescription('JPG o PNG, máximo 2 MB. Formato no permitido.')
  })

  it('puede renderizar un select', () => {
    render(
      <Campo etiqueta="Categoría" name="categoria" como="select">
        <option value="1">Alimentos</option>
      </Campo>,
    )
    expect(screen.getByLabelText('Categoría').tagName).toBe('SELECT')
  })
})
