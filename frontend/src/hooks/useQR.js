import { useState, useCallback } from 'react';
import { computeQR } from '../api/matrix';

/**
 * Encapsula el estado y la llamada a la API QR.
 * Retorna { result, loading, error, submit }.
 */
export function useQR() {
  const [result, setResult]   = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError]     = useState(null);

  const submit = useCallback(async (matrix) => {
    setLoading(true);
    setError(null);
    setResult(null);
    try {
      const data = await computeQR(matrix);
      setResult(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  return { result, loading, error, submit };
}
