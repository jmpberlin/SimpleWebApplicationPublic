const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  // Determine backend URL based on environment
  // In Docker: use service name 'backend'
  // Locally: use localhost
  const target = process.env.REACT_APP_BACKEND_URL || 'http://backend:8081';

  app.use(
    '/api',
    createProxyMiddleware({
      target: target,
      changeOrigin: true,
      pathRewrite: {
        '^/api': '', // remove /api prefix when forwarding to backend
      },
      onProxyReq: (proxyReq, req, res) => {
        console.log(`[Proxy] ${req.method} ${req.path} -> ${target}${req.path.replace('/api', '')}`);
      },
      onError: (err, req, res) => {
        console.error('[Proxy Error]', err.message);
      }
    })
  );
};
