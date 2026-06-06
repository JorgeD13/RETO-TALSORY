/**
 * MatrixInput — grilla de inputs numéricos editable.
 * Permite al usuario ingresar los valores de la matriz,
 * y agregar/eliminar filas y columnas dinámicamente.
 */
export function MatrixInput({ matrix, onChange }) {
  const rows = matrix.length;
  const cols = matrix[0]?.length ?? 0;

  function updateCell(r, c, raw) {
    const value = raw === '' ? '' : Number(raw);
    const next = matrix.map((row, i) =>
      i === r ? row.map((v, j) => (j === c ? value : v)) : row,
    );
    onChange(next);
  }

  function addRow() {
    onChange([...matrix, Array(cols).fill(0)]);
  }

  function removeRow() {
    if (rows > 1) onChange(matrix.slice(0, -1));
  }

  function addCol() {
    onChange(matrix.map((row) => [...row, 0]));
  }

  function removeCol() {
    if (cols > 1) onChange(matrix.map((row) => row.slice(0, -1)));
  }

  return (
    <div className="space-y-4">
      {/* Grilla */}
      <div className="overflow-x-auto">
        <table className="border-collapse">
          <tbody>
            {matrix.map((row, r) => (
              <tr key={r}>
                {row.map((val, c) => (
                  <td key={c} className="p-1">
                    <input
                      type="number"
                      value={val}
                      onChange={(e) => updateCell(r, c, e.target.value)}
                      className="w-20 rounded-lg border border-slate-700 bg-slate-800 px-2 py-1.5
                                 text-center text-sm text-slate-100 font-mono
                                 focus:border-violet-500 focus:outline-none focus:ring-1 focus:ring-violet-500
                                 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none
                                 [&::-webkit-inner-spin-button]:appearance-none"
                    />
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Controles de dimensión */}
      <div className="flex flex-wrap gap-2 text-xs">
        <span className="self-center text-slate-400">
          {rows} × {cols}
        </span>

        <div className="flex gap-1">
          <DimBtn onClick={addRow}    label="+ fila"    />
          <DimBtn onClick={removeRow} label="− fila"    disabled={rows <= 1} />
        </div>
        <div className="flex gap-1">
          <DimBtn onClick={addCol}    label="+ col"     />
          <DimBtn onClick={removeCol} label="− col"     disabled={cols <= 1} />
        </div>
      </div>
    </div>
  );
}

function DimBtn({ onClick, label, disabled }) {
  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled}
      className="rounded border border-slate-600 px-2 py-1 text-slate-300
                 hover:border-violet-500 hover:text-violet-400
                 disabled:cursor-not-allowed disabled:opacity-30
                 transition-colors"
    >
      {label}
    </button>
  );
}
