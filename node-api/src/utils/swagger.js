const swaggerJsdoc = require('swagger-jsdoc');

const options = {
  definition: {
    openapi: '3.0.0',
    info: {
      title: 'Statistics API',
      version: '1.0.0',
      description:
        'API Node.js que recibe las matrices Q y R calculadas por la API Go y devuelve estadísticas sobre sus valores.',
    },
    servers: [{ url: 'http://localhost:3000', description: 'Local development' }],
    components: {
      schemas: {
        StatisticsRequest: {
          type: 'object',
          required: ['q', 'r'],
          properties: {
            q: {
              type: 'array',
              items: { type: 'array', items: { type: 'number' } },
              description: 'Matriz ortogonal Q',
            },
            r: {
              type: 'array',
              items: { type: 'array', items: { type: 'number' } },
              description: 'Matriz triangular superior R',
            },
          },
        },
        StatisticsResponse: {
          type: 'object',
          properties: {
            max:         { type: 'number', example: 175 },
            min:         { type: 'number', example: -68 },
            average:     { type: 'number', example: 12.5 },
            sum:         { type: 'number', example: 100 },
            isDiagonalQ: { type: 'boolean', example: false },
            isDiagonalR: { type: 'boolean', example: false },
          },
        },
        ErrorResponse: {
          type: 'object',
          properties: {
            error:   { type: 'string', example: 'bad_request' },
            message: { type: 'string', example: '"q" and "r" must be arrays' },
          },
        },
      },
    },
  },
  apis: ['./src/routes/*.js'],
};

const swaggerSpec = swaggerJsdoc(options);

module.exports = swaggerSpec;
