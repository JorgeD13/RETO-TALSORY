const { computeStatistics } = require('../services/statistics.service');
const { validateMatrix2D } = require('../utils/matrix.utils');

/**
 * POST /api/v1/statistics
 * Recibe { q, r } y devuelve las estadísticas calculadas.
 *
 * @param {import('express').Request}  req
 * @param {import('express').Response} res
 * @param {import('express').NextFunction} next
 */
function getStatistics(req, res, next) {
  try {
    const { q, r } = req.body;

    const qValidation = validateMatrix2D(q);
    if (!qValidation.valid) {
      return res.status(400).json({ error: 'bad_request', message: `q: ${qValidation.message}` });
    }

    const rValidation = validateMatrix2D(r);
    if (!rValidation.valid) {
      return res.status(400).json({ error: 'bad_request', message: `r: ${rValidation.message}` });
    }

    const result = computeStatistics(q, r);
    console.log(`[statistics] max=${result.max.toFixed(4)} min=${result.min.toFixed(4)} avg=${result.average.toFixed(4)}`);

    return res.status(200).json(result);
  } catch (err) {
    next(err);
  }
}

module.exports = { getStatistics };
