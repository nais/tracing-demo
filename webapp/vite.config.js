export default {
  build: {
    manifest: true,
    rollupOptions: {
      external: [
        "./nais.js",
      ],
    },
  },
	server: {
		proxy: {
			'/grafana': {
				target: 'http://127.0.0.1:12347',
				changeOrigin: true,
				rewrite: (path) => path.replace(/^\/grafana/, ''),
			},
			'/api': {
				target: 'http://127.0.0.1:8080',
				changeOrigin: true,
				rewrite: (path) => path.replace(/^\/api/, ''),
			}
		}
	}
}
