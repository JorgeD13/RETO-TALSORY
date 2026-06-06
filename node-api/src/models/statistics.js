/**
 * @typedef {Object} MatrixPayload
 * @property {number[][]} q - Matriz ortogonal Q
 * @property {number[][]} r - Matriz triangular superior R
 */

/**
 * @typedef {Object} StatisticsResult
 * @property {number}  max          - Valor máximo entre todos los elementos de Q y R
 * @property {number}  min          - Valor mínimo entre todos los elementos de Q y R
 * @property {number}  average      - Promedio de todos los elementos
 * @property {number}  sum          - Suma de todos los elementos
 * @property {boolean} isDiagonalQ  - True si Q es matriz diagonal
 * @property {boolean} isDiagonalR  - True si R es matriz diagonal
 */
