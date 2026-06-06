/**
 * Utilidades para operaciones sobre matrices 2D.
 */

/**
 * Aplana una matriz 2D en un array 1D.
 * @param {number[][]} matrix
 * @returns {number[]}
 */
function flatten(matrix) {
  return matrix.flat();
}

/**
 * Devuelve true si la matriz es cuadrada.
 * @param {number[][]} matrix
 * @returns {boolean}
 */
function isSquare(matrix) {
  return matrix.every((row) => row.length === matrix.length);
}

/**
 * Devuelve true si todos los elementos fuera de la diagonal principal son
 * menores que EPSILON en valor absoluto.
 * @param {number[][]} matrix
 * @param {number} [epsilon=1e-9]
 * @returns {boolean}
 */
function isDiagonal(matrix, epsilon = 1e-9) {
  for (let i = 0; i < matrix.length; i++) {
    for (let j = 0; j < matrix[i].length; j++) {
      if (i !== j && Math.abs(matrix[i][j]) > epsilon) return false;
    }
  }
  return true;
}

/**
 * Valida que un valor sea una matriz 2D de números no vacía y rectangular.
 * @param {unknown} value
 * @returns {{ valid: boolean, message?: string }}
 */
function validateMatrix2D(value) {
  if (!Array.isArray(value)) return { valid: false, message: 'must be an array' };
  if (value.length === 0) return { valid: false, message: 'must not be empty' };
  const cols = value[0].length;
  if (cols === 0) return { valid: false, message: 'rows must not be empty' };
  for (let i = 0; i < value.length; i++) {
    if (!Array.isArray(value[i])) return { valid: false, message: `row ${i} is not an array` };
    if (value[i].length !== cols) {
      return { valid: false, message: `row ${i} has ${value[i].length} columns, expected ${cols}` };
    }
    if (!value[i].every((v) => typeof v === 'number' && isFinite(v))) {
      return { valid: false, message: `row ${i} contains non-numeric values` };
    }
  }
  return { valid: true };
}

module.exports = { flatten, isSquare, isDiagonal, validateMatrix2D };
