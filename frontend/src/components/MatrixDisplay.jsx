/**
 * MatrixDisplay — muestra una matriz 2D con formato numérico fijo.
 */
export function MatrixDisplay({ label, matrix, accentClass = 'text-violet-400' }) {
  if (!matrix?.length) return null;

  return (
    <div className="space-y-2">
      <h3 className={`font-semibold text-sm uppercase tracking-widest ${accentClass}`}>
        {label}
      </h3>
      <div className="overflow-x-auto rounded-xl border border-slate-700 bg-slate-800/50 p-4">
        <table className="border-collapse font-mono text-sm">
          <tbody>
            {matrix.map((row, r) => (
              <tr key={r}>
                {/* Corchete izquierdo solo en la primera columna */}
                <td className="pr-1 text-slate-500 select-none">
                  {r === 0 ? '⎡' : r === matrix.length - 1 ? '⎣' : '⎢'}
                </td>

                {row.map((val, c) => (
                  <td
                    key={c}
                    className="px-2 py-0.5 text-right tabular-nums text-slate-200"
                  >
                    {formatNum(val)}
                  </td>
                ))}

                {/* Corchete derecho */}
                <td className="pl-1 text-slate-500 select-none">
                  {r === 0 ? '⎤' : r === matrix.length - 1 ? '⎦' : '⎥'}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

function formatNum(n) {
  if (Math.abs(n) < 1e-9) return '0.0000';
  return n.toFixed(4);
}
