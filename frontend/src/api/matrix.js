const BASE_URL = import.meta.env.VITE_API_URL ?? '';

/**
 * Envía una matriz al endpoint QR de la Go API y retorna
 * { q, r, statistics } o lanza un Error con el mensaje del servidor.
 *
 * @param {number[][]} matrix
 * @returns {Promise<{ q: number[][], r: number[][], statistics: object }>}
 */
export async function computeQR(matrix) {
  const res = await fetch(`${BASE_URL}/api/v1/matrix/qr`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ matrix }),
  });

  const data = await res.json();

  if (!res.ok) {
    throw new Error(data.message ?? `Error ${res.status}`);
  }

  return data;
}
