export default {
	server: {
		proxy: {
			'/grafana': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				rewrite: (path) => path.replace(/^\/grafana/, ''),
			}
		}
	}
}
