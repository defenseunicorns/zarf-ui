import { sveltekit } from '@sveltejs/kit/vite';
import type { UserConfig } from 'vite';

const backendAPI = {
	target: 'http://127.0.0.1:3333',
	changeOrigin: true,
	secure: false,
	ws: true,
};

const config: UserConfig = {
	plugins: [sveltekit()],
	server: {
		proxy: {
			'/api': backendAPI,
			'/sbom-viewer/': backendAPI,
		},
		fs: {
			strict: false,
		},
	},
	optimizeDeps: {
		include: [
			'prismjs',
			'ansi-to-html',
			'@floating-ui/dom',
			'prismjs/components/prism-yaml',
			'yaml',
		],
	},
	resolve: {
		alias: {
			'@images': __dirname + '/images',
			'@ui': __dirname + '/node_modules/@defense-unicorns/unicorn-ui',
		},
	},
};

export default config;
