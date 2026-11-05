import { useId } from 'react'

// Campo de formulario accesible:
// - la etiqueta está asociada al control (htmlFor/id)
// - el error se anuncia con aria-invalid y aria-describedby
// - la ayuda opcional también queda asociada
// Para usarlo con <select> o <textarea>, pasar como="select" o como="textarea".
export default function Campo({ etiqueta, error, ayuda, como = 'input', obligatorio = false, children, ...props }) {
  const id = useId()
  const idError = `${id}-error`
  const idAyuda = `${id}-ayuda`
  const describe = [ayuda && idAyuda, error && idError].filter(Boolean).join(' ') || undefined
  const Control = como

  return (
    <div className="campo">
      <label htmlFor={id}>
        {etiqueta}
        {obligatorio && <span aria-hidden="true"> *</span>}
      </label>
      {ayuda && (
        <p id={idAyuda} className="campo-ayuda">
          {ayuda}
        </p>
      )}
      <Control
        id={id}
        aria-invalid={error ? 'true' : undefined}
        aria-describedby={describe}
        required={obligatorio}
        {...props}
      >
        {children}
      </Control>
      {error && (
        <p id={idError} className="campo-error">
          {error}
        </p>
      )}
    </div>
  )
}
