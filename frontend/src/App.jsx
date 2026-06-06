import { useState } from 'react';
import { useQR } from './hooks/useQR';
import { MatrixInput } from './components/MatrixInput';
import { MatrixDisplay } from './components/MatrixDisplay';
import { StatisticsDisplay } from './components/StatisticsDisplay';

// Matriz de ejemplo del reto precargada para facilitar la demo.
const INITIAL_MATRIX = [
  [12, -51, 4],
  [6, 167, -68],
  [-4, 24, -41],
];

export default function App() {
  const [matrix, setMatrix] = useState(INITIAL_MATRIX);
  const { result, loading, error, submit } = useQR();

  function handleSubmit(e) {
    e.preventDefault();

    // Convertir cualquier string vacío a 0 antes de enviar.
    const clean = matrix.map((row) =>
      row.map((v) => (v === '' || isNaN(Number(v)) ? 0 : Number(v))),
    );
    submit(clean);
  }

  return (
    <div className="min-h-screen bg-slate-950 px-4 py-10">
      <div className="mx-auto max-w-5xl space-y-10">

        {/* Header */}
        <header className="text-center space-y-2">
          <h1 className="text-3xl font-bold tracking-tight text-slate-100">
            Factorización <span className="text-violet-400">QR</span>
          </h1>
          <p className="text-slate-400 text-sm">
            Ingresa una matriz numérica y calcula sus factores Q y R con estadísticas
          </p>
        </header>

        {/* Formulario */}
        <form
          onSubmit={handleSubmit}
          className="rounded-2xl border border-slate-800 bg-slate-900 p-6 space-y-6"
        >
          <div className="space-y-1">
            <label className="text-xs font-semibold uppercase tracking-widest text-slate-400">
              Matriz de entrada
            </label>
            <MatrixInput matrix={matrix} onChange={setMatrix} />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-xl bg-violet-600 py-3 text-sm font-semibold text-white
                       hover:bg-violet-500 active:bg-violet-700
                       disabled:cursor-not-allowed disabled:opacity-50
                       transition-colors focus:outline-none focus:ring-2 focus:ring-violet-500"
          >
            {loading ? (
              <span className="flex items-center justify-center gap-2">
                <Spinner /> Calculando…
              </span>
            ) : (
              'Calcular factorización QR'
            )}
          </button>

          {/* Error */}
          {error && (
            <div className="rounded-xl border border-red-800 bg-red-950/40 px-4 py-3 text-sm text-red-400">
              <span className="font-semibold">Error:</span> {error}
            </div>
          )}
        </form>

        {/* Resultados */}
        {result && (
          <div className="space-y-8 rounded-2xl border border-slate-800 bg-slate-900 p-6">
            <div className="grid gap-8 md:grid-cols-2">
              <MatrixDisplay
                label="Matriz Q (ortogonal)"
                matrix={result.q}
                accentClass="text-violet-400"
              />
              <MatrixDisplay
                label="Matriz R (triangular superior)"
                matrix={result.r}
                accentClass="text-sky-400"
              />
            </div>

            <hr className="border-slate-800" />

            <StatisticsDisplay stats={result.statistics} />
          </div>
        )}

        {/* Footer */}
        <footer className="text-center text-xs text-slate-600">
          Go API · Node API · Gonum QR
        </footer>
      </div>
    </div>
  );
}

function Spinner() {
  return (
    <svg
      className="h-4 w-4 animate-spin"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83" />
    </svg>
  );
}
