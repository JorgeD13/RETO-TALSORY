/**
 * Manejador centralizado de errores de Express.
 * Debe registrarse como último middleware (4 parámetros).
 */
// eslint-disable-next-line no-unused-vars
function errorHandler(err, req, res, next) {
  console.error(`[error] ${err.message}`);
  res.status(err.status || 500).json({
    error: 'internal_server_error',
    message: err.message || 'An unexpected error occurred',
  });
}

module.exports = { errorHandler };
