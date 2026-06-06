/**
 * Tests unitarios del servicio de estadísticas.
 * Prueba la lógica de cálculo pura (sin HTTP) sobre los elementos
 * de las matrices Q y R: max, min, average, sum, isDiagonalQ, isDiagonalR.
 */

const { computeStatistics } = require('../services/statistics.service');

// Matrices de referencia del reto: factorización QR de [[12,-51,4],[6,167,-68],[-4,24,-41]]
const Q = [
  [-0.8571,  0.3943,  0.3314],
  [-0.4286, -0.9029, -0.0292],
  [ 0.2857, -0.1714,  0.9429],
];
const R = [
  [-14,    -21,  14],
  [  0,  -175,  70],
  [  0,     0, -35],
];

describe('computeStatistics — cálculo sobre Q∪R combinadas', () => {

  test('max — el mayor valor de Q∪R es 70 (elemento de R)', () => {
    // Verifica que se busca el máximo sobre todos los valores de ambas matrices,
    // no solo sobre una. El mayor elemento de Q es ~0.9429, el de R es 70.
    const result = computeStatistics(Q, R);
    expect(result.max).toBeCloseTo(70, 4);
  });

  test('min — el menor valor de Q∪R es -175 (elemento de R)', () => {
    // Análogamente al max: el menor de Q es ~-0.9029, el de R es -175.
    const result = computeStatistics(Q, R);
    expect(result.min).toBeCloseTo(-175, 4);
  });

  test('sum — es la suma aritmética de todos los elementos de Q y R', () => {
    // La suma se calcula sobre el conjunto completo Q∪R (18 elementos).
    const expected = [...Q.flat(), ...R.flat()].reduce((a, b) => a + b, 0);
    const result = computeStatistics(Q, R);
    expect(result.sum).toBeCloseTo(expected, 5);
  });

  test('average — es sum / (|Q| + |R|)', () => {
    // El promedio divide la suma entre el total de elementos (no solo filas).
    const allValues = [...Q.flat(), ...R.flat()];
    const expected = allValues.reduce((a, b) => a + b, 0) / allValues.length;
    const result = computeStatistics(Q, R);
    expect(result.average).toBeCloseTo(expected, 5);
  });

});

describe('computeStatistics — isDiagonalQ / isDiagonalR', () => {

  test('Q ortogonal no es diagonal — tiene elementos fuera de la diagonal principal', () => {
    // La matriz Q del ejemplo tiene valores ~0.39, ~0.33, etc. fuera de la diagonal.
    const result = computeStatistics(Q, R);
    expect(result.isDiagonalQ).toBe(false);
  });

  test('R triangular superior no es diagonal — tiene elementos encima de la diagonal', () => {
    // R tiene -21, 14, 70 fuera de la diagonal.
    const result = computeStatistics(Q, R);
    expect(result.isDiagonalR).toBe(false);
  });

  test('detecta matrices verdaderamente diagonales en ambas entradas', () => {
    // Una matriz diagonal tiene ceros en todas las posiciones (i,j) donde i≠j.
    const diag3x3 = [[1, 0, 0], [0, 2, 0], [0, 0, 3]];
    const identity = [[1, 0], [0, 1]];
    const result = computeStatistics(diag3x3, identity);
    expect(result.isDiagonalQ).toBe(true);
    expect(result.isDiagonalR).toBe(true);
  });

  test('detecta cuando solo una de las dos matrices es diagonal', () => {
    // Caso mixto: Q no diagonal, R diagonal.
    // Importante para que el cliente sepa exactamente cuál cumple la propiedad.
    const nonDiag = [[1, 2], [0, 3]];
    const diag    = [[1, 0], [0, 1]];
    const result = computeStatistics(nonDiag, diag);
    expect(result.isDiagonalQ).toBe(false);
    expect(result.isDiagonalR).toBe(true);
  });

  test('valores menores que epsilon (1e-9) se tratan como cero para isDiagonal', () => {
    // Los elementos de Q y R tienen errores de punto flotante (e.g. 1.2e-16).
    // Sin epsilon, isDiagonal devolvería false incorrectamente.
    const nearlyDiag = [[5, 0], [1e-15, 3]]; // 1e-15 < 1e-9 → se trata como 0
    const result = computeStatistics(nearlyDiag, [[1]]);
    expect(result.isDiagonalQ).toBe(true);
  });

});

describe('computeStatistics — casos límite', () => {

  test('matrices 1×1 — el mínimo caso posible con un solo elemento cada una', () => {
    // Con una sola celda por matriz, max=5, min=-3, sum=2, avg=1.
    const result = computeStatistics([[5]], [[-3]]);
    expect(result.max).toBe(5);
    expect(result.min).toBe(-3);
    expect(result.sum).toBeCloseTo(2, 10);
    expect(result.average).toBeCloseTo(1, 10);
    // Una matriz 1×1 siempre es diagonal por definición.
    expect(result.isDiagonalQ).toBe(true);
    expect(result.isDiagonalR).toBe(true);
  });

  test('todos los valores negativos — max y min son negativos', () => {
    // Verifica que no se asuma que el máximo siempre es positivo.
    const result = computeStatistics([[-1, -2]], [[-3, -4]]);
    expect(result.max).toBe(-1);
    expect(result.min).toBe(-4);
  });

  test('todos los valores iguales — max = min = average = valor', () => {
    // Caso degenerado: todos los elementos son el mismo número.
    const result = computeStatistics([[3, 3]], [[3, 3]]);
    expect(result.max).toBe(3);
    expect(result.min).toBe(3);
    expect(result.average).toBe(3);
  });

  test('valores muy grandes — no hay overflow en la suma (JS usa float64)', () => {
    // Number.MAX_SAFE_INTEGER es ~9e15. La suma de dos valores así
    // puede perder precisión, pero no debe crashear.
    const big = [[1e15, 0], [0, 1e15]];
    const result = computeStatistics(big, [[0]]);
    expect(result.sum).toBeCloseTo(2e15, -10);
    expect(result.max).toBeCloseTo(1e15, -10);
  });

});
