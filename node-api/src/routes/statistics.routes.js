const { Router } = require('express');
const { getStatistics } = require('../controllers/statistics.controller');

const router = Router();

/**
 * @openapi
 * /api/v1/statistics:
 *   post:
 *     summary: Calcula estadísticas sobre las matrices Q y R
 *     description: >
 *       Recibe las matrices Q y R resultantes de la factorización QR y
 *       devuelve estadísticas (max, min, promedio, suma) y si cada matriz es diagonal.
 *     tags:
 *       - statistics
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/StatisticsRequest'
 *           example:
 *             q: [[-0.857, 0.394, 0.331], [-0.429, -0.903, -0.029], [0.286, -0.171, 0.943]]
 *             r: [[-14, -21, 14], [0, -175, 70], [0, 0, -35]]
 *     responses:
 *       200:
 *         description: Estadísticas calculadas correctamente
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/StatisticsResponse'
 *       400:
 *         description: Payload inválido
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/ErrorResponse'
 *       500:
 *         description: Error interno del servidor
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/ErrorResponse'
 */
router.post('/statistics', getStatistics);

module.exports = router;
