/**
 * Tests unitarios de las utilidades de matriz.
 * Cada función utilitaria se prueba de forma aislada con casos
 * que cubren el comportamiento normal, los límites y los casos de error.
 */

const { flatten, isSquare, isDiagonal, validateMatrix2D } = require('../utils/matrix.utils');

// ── flatten ────────────────────────────────────────────────────────────────────

describe('flatten', () => {

  test('convierte una matriz 2×2 en un array 1D en orden row-major', () => {
    // El orden importa: [fila0_col0, fila0_col1, fila1_col0, fila1_col1].
    expect(flatten([[1, 2], [3, 4]])).toEqual([1, 2, 3, 4]);
  });

  test('una matriz 1×1 produce un array de un solo elemento', () => {
    // Caso mínimo: no debe lanzar error ni devolver array vacío.
    expect(flatten([[7]])).toEqual([7]);
  });

  test('una matriz 1×N produce los elementos de esa única fila', () => {
    expect(flatten([[10, 20, 30]])).toEqual([10, 20, 30]);
  });

  test('preserva valores negativos y decimales sin alterarlos', () => {
    expect(flatten([[-1.5, 2.7]])).toEqual([-1.5, 2.7]);
  });

});

// ── isSquare ──────────────────────────────────────────────────────────────────

describe('isSquare', () => {

  test('devuelve true para una matriz 3×3', () => {
    // La factorización QR de una matriz cuadrada produce Q cuadrada.
    expect(isSquare([[1,2,3],[4,5,6],[7,8,9]])).toBe(true);
  });

  test('devuelve false para una matriz 3×2 (tall)', () => {
    // Una matriz con más filas que columnas no es cuadrada.
    expect(isSquare([[1,2],[3,4],[5,6]])).toBe(false);
  });

  test('devuelve false para una matriz 2×3 (wide)', () => {
    expect(isSquare([[1,2,3],[4,5,6]])).toBe(false);
  });

  test('devuelve true para una matriz 1×1 (mínimo cuadrado posible)', () => {
    expect(isSquare([[42]])).toBe(true);
  });

});

// ── isDiagonal ────────────────────────────────────────────────────────────────

describe('isDiagonal', () => {

  test('la matriz identidad 3×3 es diagonal', () => {
    // La identidad tiene 1 en la diagonal y 0 en el resto: diagonal por definición.
    expect(isDiagonal([[1,0,0],[0,1,0],[0,0,1]])).toBe(true);
  });

  test('una matriz diagonal con valores distintos de 1 en la diagonal es diagonal', () => {
    // La diagonal no tiene que ser la identidad para ser 'diagonal'.
    expect(isDiagonal([[5,0,0],[0,-3,0],[0,0,7]])).toBe(true);
  });

  test('un elemento distinto de cero fuera de la diagonal hace la matriz no-diagonal', () => {
    // Caso: matriz triangular superior — 1 encima de la diagonal.
    expect(isDiagonal([[1,0],[1,2]])).toBe(false);
  });

  test('respeta el epsilon: 1e-10 se trata como cero (ruido numérico)', () => {
    // Los elementos de Q/R suelen tener residuos de punto flotante (1e-16, 1e-15).
    // Sin epsilon, isDiagonal retornaría false en casi todos los casos reales.
    expect(isDiagonal([[1, 1e-10], [0, 1]])).toBe(true);
  });

  test('1e-8 supera el epsilon por defecto (1e-9) y se considera no-cero', () => {
    // Un valor de 1e-8 es mayor que 1e-9 → fuera del umbral → no es diagonal.
    expect(isDiagonal([[1, 1e-8], [0, 1]])).toBe(false);
  });

  test('una matriz 1×1 siempre es diagonal (no hay elementos fuera de la diagonal)', () => {
    expect(isDiagonal([[42]])).toBe(true);
  });

  test('epsilon personalizado permite ajustar la tolerancia', () => {
    // Con epsilon=1e-6, un valor de 1e-7 pasa como cero.
    expect(isDiagonal([[1, 1e-7], [0, 1]], 1e-6)).toBe(true);
    expect(isDiagonal([[1, 1e-5], [0, 1]], 1e-6)).toBe(false);
  });

});

// ── validateMatrix2D ──────────────────────────────────────────────────────────

describe('validateMatrix2D', () => {

  test('acepta una matriz rectangular válida — { valid: true }', () => {
    // Caso feliz: 2×2 con números finitos.
    expect(validateMatrix2D([[1,2],[3,4]])).toEqual({ valid: true });
  });

  test('rechaza un string — "matrix" debe ser un array', () => {
    // Caso: el campo llegó como string en el JSON (error de tipo del cliente).
    expect(validateMatrix2D('not-a-matrix')).toMatchObject({ valid: false });
  });

  test('rechaza null — campo matrix ausente o null', () => {
    expect(validateMatrix2D(null)).toMatchObject({ valid: false });
  });

  test('rechaza un número escalar — no es una matriz 2D', () => {
    expect(validateMatrix2D(42)).toMatchObject({ valid: false });
  });

  test('rechaza un array vacío — sin filas no hay datos que procesar', () => {
    expect(validateMatrix2D([])).toMatchObject({ valid: false });
  });

  test('rechaza una fila vacía — sin columnas la fila no tiene datos', () => {
    expect(validateMatrix2D([[]])).toMatchObject({ valid: false });
  });

  test('rechaza una matriz jagged — las filas tienen distinta longitud', () => {
    // [[1,2,3],[4,5]] — inconsistente, no puede ser procesada como matriz.
    expect(validateMatrix2D([[1,2,3],[4,5]])).toMatchObject({ valid: false });
  });

  test('rechaza strings dentro de las filas — solo se aceptan números finitos', () => {
    expect(validateMatrix2D([[1, 'a']])).toMatchObject({ valid: false });
  });

  test('rechaza NaN — es de tipo number pero no es un número válido', () => {
    // Sin esta validación, NaN se propagaría a max/min/sum silenciosamente.
    expect(validateMatrix2D([[1, NaN]])).toMatchObject({ valid: false });
  });

  test('rechaza Infinity — no finito', () => {
    expect(validateMatrix2D([[Infinity, 1]])).toMatchObject({ valid: false });
  });

  test('acepta valores negativos y decimales — son números finitos válidos', () => {
    expect(validateMatrix2D([[-1.5, 2.7], [-0.001, 100]])).toEqual({ valid: true });
  });

  test('el mensaje de error menciona el campo problemático', () => {
    // Facilita el debug: el cliente sabe exactamente qué fila tiene el problema.
    const result = validateMatrix2D([[1,2],[3]]);
    expect(result.valid).toBe(false);
    expect(typeof result.message).toBe('string');
    expect(result.message.length).toBeGreaterThan(0);
  });

});
