/**
 * Tests de integración del controller de estadísticas.
 * Usa supertest para hacer requests HTTP reales contra la app Express
 * sin necesidad de levantar un puerto de red.
 *
 * Cubre: parsing del body, validaciones de entrada, respuesta correcta
 * y formato del error response ante payloads inválidos.
 */

const request = require('supertest');
const app = require('../app');

// Payload de referencia: matrices Q y R del ejemplo del reto.
const validPayload = {
  q: [
    [-0.8571,  0.3943,  0.3314],
    [-0.4286, -0.9029, -0.0292],
    [ 0.2857, -0.1714,  0.9429],
  ],
  r: [
    [-14,    -21,  14],
    [  0,  -175,  70],
    [  0,     0, -35],
  ],
};

describe('POST /api/v1/statistics', () => {

  // ── Camino feliz ────────────────────────────────────────────────────────────

  test('200 — payload válido devuelve todas las métricas', async () => {
    // Caso: el cliente envía q y r correctos → debe recibir las 6 métricas.
    const res = await request(app).post('/api/v1/statistics').send(validPayload);

    expect(res.status).toBe(200);

    // Verificar que existen todos los campos del contrato.
    expect(typeof res.body.max).toBe('number');
    expect(typeof res.body.min).toBe('number');
    expect(typeof res.body.average).toBe('number');
    expect(typeof res.body.sum).toBe('number');
    expect(typeof res.body.isDiagonalQ).toBe('boolean');
    expect(typeof res.body.isDiagonalR).toBe('boolean');
  });

  test('200 — max y min son correctos para el payload de referencia', async () => {
    // Caso: verifica que los valores extremos se calculen sobre Q∪R juntos,
    // no solo sobre una de las matrices por separado.
    const res = await request(app).post('/api/v1/statistics').send(validPayload);

    expect(res.status).toBe(200);
    // El valor mayor de R es 70, el mayor de Q es ~0.9429.
    expect(res.body.max).toBeCloseTo(70, 4);
    // El valor menor de R es -175, el menor de Q es ~-0.9029.
    expect(res.body.min).toBeCloseTo(-175, 4);
  });

  test('200 — matrices diagonales → isDiagonalQ y isDiagonalR son true', async () => {
    // Caso: cuando ambas matrices son diagonales, los flags deben reportarlo.
    const payload = {
      q: [[2, 0, 0], [0, 3, 0], [0, 0, 5]],
      r: [[1, 0], [0, 4]],
    };
    const res = await request(app).post('/api/v1/statistics').send(payload);

    expect(res.status).toBe(200);
    expect(res.body.isDiagonalQ).toBe(true);
    expect(res.body.isDiagonalR).toBe(true);
  });

  test('200 — matrices 1×1 — caso mínimo válido', async () => {
    // Caso límite: la matriz más pequeña posible aún debe procesarse.
    const res = await request(app)
      .post('/api/v1/statistics')
      .send({ q: [[7]], r: [[-3]] });

    expect(res.status).toBe(200);
    expect(res.body.max).toBe(7);
    expect(res.body.min).toBe(-3);
    expect(res.body.sum).toBeCloseTo(4, 10);
    expect(res.body.average).toBeCloseTo(2, 10);
  });

  // ── Validación de campo 'q' ─────────────────────────────────────────────────

  test('400 — campo q ausente en el body', async () => {
    // Caso: el cliente envía solo r sin q.
    // validateMatrix2D detecta que q no es un array.
    const res = await request(app)
      .post('/api/v1/statistics')
      .send({ r: validPayload.r });

    expect(res.status).toBe(400);
    expect(res.body.error).toBe('bad_request');
    // El mensaje debe mencionar 'q' para que el cliente identifique el campo.
    expect(res.body.message).toMatch(/q/);
  });

  test('400 — q es un string en lugar de array', async () => {
    // Caso: tipo incorrecto → el campo existe pero no es una matriz.
    const res = await request(app)
      .post('/api/v1/statistics')
      .send({ q: 'not-an-array', r: validPayload.r });

    expect(res.status).toBe(400);
    expect(res.body.error).toBe('bad_request');
  });

  test('400 — q es un array vacío', async () => {
    // Caso: array vacío no es una matriz válida (sin filas, sin factorización).
    const res = await request(app)
      .post('/api/v1/statistics')
      .send({ q: [], r: validPayload.r });

    expect(res.status).toBe(400);
    expect(res.body.error).toBe('bad_request');
  });

  test('400 — q tiene una fila vacía', async () => {
    // Caso: [[], [1,2]] — una fila sin columnas no es una matriz rectangular.
    const res = await request(app)
      .post('/api/v1/statistics')
      .send({ q: [[]], r: validPayload.r });

    expect(res.status).toBe(400);
    expect(res.body.error).toBe('bad_request');
  });

  test('400 — q es una matriz jagged (filas con distinta longitud)', async () => {
    // Caso: [[1,2,3],[4,5]] — no rectangular, los cálculos serían incorrectos.
    const res = await request(app)
      .post('/api/v1/statistics')
      .send({ q: [[1, 2, 3], [4, 5]], r: validPayload.r });

    expect(res.status).toBe(400);
    expect(res.body.error).toBe('bad_request');
  });

  test('400 — q contiene valores no numéricos', async () => {
    // Caso: [["a", "b"]] — los strings no pueden participar en operaciones
    // matemáticas y produciría NaN silencioso sin validación.
    const res = await request(app)
      .post('/api/v1/statistics')
      .send({ q: [['a', 'b']], r: validPayload.r });

    expect(res.status).toBe(400);
    expect(res.body.error).toBe('bad_request');
  });

  // ── Validación de campo 'r' ─────────────────────────────────────────────────

  test('400 — campo r ausente en el body', async () => {
    // Caso: el cliente envía solo q sin r.
    const res = await request(app)
      .post('/api/v1/statistics')
      .send({ q: validPayload.q });

    expect(res.status).toBe(400);
    expect(res.body.error).toBe('bad_request');
    expect(res.body.message).toMatch(/r/);
  });

  test('400 — r contiene NaN', async () => {
    // Caso: NaN es de tipo "number" en JS pero no es un número finito válido.
    // Sin validación explícita, NaN se propaga silenciosamente a max/min/sum.
    const res = await request(app)
      .post('/api/v1/statistics')
      .send({ q: validPayload.q, r: [[NaN, 0], [0, NaN]] });

    expect(res.status).toBe(400);
    expect(res.body.error).toBe('bad_request');
  });

  // ── Formato de la respuesta de error ───────────────────────────────────────

  test('400 — respuesta de error sigue el esquema { error, message }', async () => {
    // Contrato: todos los errores deben tener los campos 'error' y 'message'
    // para que el cliente pueda manejarlos uniformemente.
    const res = await request(app)
      .post('/api/v1/statistics')
      .send({ q: null, r: null });

    expect(res.status).toBe(400);
    expect(res.body).toHaveProperty('error');
    expect(res.body).toHaveProperty('message');
    expect(typeof res.body.error).toBe('string');
    expect(typeof res.body.message).toBe('string');
  });

});
