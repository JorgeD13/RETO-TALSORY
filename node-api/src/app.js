require('dotenv').config();

const express = require('express');
const swaggerUi = require('swagger-ui-express');
const swaggerSpec = require('./utils/swagger');
const statisticsRoutes = require('./routes/statistics.routes');
const { errorHandler } = require('./middleware/errorHandler');

const app = express();

app.use(express.json());

// Swagger UI — disponible en /api-docs
app.use('/api-docs', swaggerUi.serve, swaggerUi.setup(swaggerSpec));

// Health-check
app.get('/health', (_req, res) => res.json({ status: 'ok', env: process.env.NODE_ENV }));

// Rutas de la API
app.use('/api/v1', statisticsRoutes);

// Manejador de errores global (debe ir al final)
app.use(errorHandler);

// La app se exporta SIN llamar a listen() para que los tests puedan
// importarla sin ocupar un puerto de red.
module.exports = app;
