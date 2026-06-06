/**
 * StatisticsDisplay — muestra las métricas calculadas por la API Node.js
 * en tarjetas visuales agrupadas por tipo.
 */
export function StatisticsDisplay({ stats }) {
  if (!stats) return null;

  const numeric = [
    { label: 'Máximo',   value: fmt(stats.max),     icon: '↑' },
    { label: 'Mínimo',   value: fmt(stats.min),     icon: '↓' },
    { label: 'Promedio', value: fmt(stats.average), icon: '∅' },
    { label: 'Suma',     value: fmt(stats.sum),     icon: 'Σ' },
  ];

  return (
    <div className="space-y-4">
      <h3 className="font-semibold text-sm uppercase tracking-widest text-emerald-400">
        Estadísticas
      </h3>

      {/* Métricas numéricas */}
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
        {numeric.map(({ label, value, icon }) => (
          <div
            key={label}
            className="rounded-xl border border-slate-700 bg-slate-800/50 p-4 text-center"
          >
            <div className="text-2xl text-slate-500 mb-1">{icon}</div>
            <div className="text-lg font-mono font-semibold text-slate-100">{value}</div>
            <div className="text-xs text-slate-400 mt-1">{label}</div>
          </div>
        ))}
      </div>

      {/* Flags de diagonal */}
      <div className="grid grid-cols-2 gap-3">
        <DiagonalBadge label="Q es diagonal" value={stats.isDiagonalQ} />
        <DiagonalBadge label="R es diagonal" value={stats.isDiagonalR} />
      </div>
    </div>
  );
}

function DiagonalBadge({ label, value }) {
  return (
    <div
      className={`flex items-center gap-3 rounded-xl border px-4 py-3
        ${value
          ? 'border-emerald-700 bg-emerald-950/40 text-emerald-300'
          : 'border-slate-700 bg-slate-800/50 text-slate-400'
        }`}
    >
      <span className="text-xl">{value ? '✓' : '✗'}</span>
      <span className="text-sm font-medium">{label}</span>
    </div>
  );
}

function fmt(n) {
  if (n === undefined || n === null) return '—';
  if (Math.abs(n) >= 1e6 || (Math.abs(n) < 0.001 && n !== 0)) {
    return n.toExponential(3);
  }
  return n.toFixed(4);
}
