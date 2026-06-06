const { flatten, isDiagonal } = require('../utils/matrix.utils');

/**
 * Calcula estadísticas sobre los elementos de las matrices Q y R.
 *
 * @param {number[][]} q - Matriz ortogonal Q
 * @param {number[][]} r - Matriz triangular superior R
 * @returns {import('../models/statistics').StatisticsResult}
 */
function computeStatistics(q, r) {
  const allValues = [...flatten(q), ...flatten(r)];

  let sum = 0;
  let max = -Infinity;
  let min = Infinity;

  for (const v of allValues) {
    sum += v;
    if (v > max) max = v;
    if (v < min) min = v;
  }

  const average = allValues.length > 0 ? sum / allValues.length : 0;

  return {
    max,
    min,
    average,
    sum,
    isDiagonalQ: isDiagonal(q),
    isDiagonalR: isDiagonal(r),
  };
}

module.exports = { computeStatistics };
