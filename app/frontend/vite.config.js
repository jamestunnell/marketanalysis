import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({mode}) => {
    const env = loadEnv(mode, process.cwd(), '');

    console.log("loaded vite env: VITE_API_BASE_URL=%s", env.VITE_API_BASE_URL);

    return {
        server: {
            proxy: {
                '/api': {
                    target: 'http://localhost:4002',
                    changeOrigin: true,
                    rewrite: (path) => {
                        const modified = path.replace(/^\/api/, '')

                        console.log('rewriting path %s -> %s', path, modified)

                        return modified
                    },
                },
            },
        },
    };
});