// Punto de entrada que arranca el servidor HTTP.
// Separado de app.js para que los tests puedan importar la app sin
// abrir un puerto de red.
const app = require('./app');

const PORT = process.env.PORT || 3000;

app.listen(PORT, () => {
  console.log(`[node-api] running on port ${PORT} (${process.env.NODE_ENV})`);
  console.log(`[node-api] swagger docs → http://localhost:${PORT}/api-docs`);
});
