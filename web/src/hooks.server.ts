import {createImageOptimizer} from 'sveltekit-image-optimize';
import type {Handle} from '@sveltejs/kit';
import {createFileSystemCache} from 'sveltekit-image-optimize/cache-adapters';

const cache = createFileSystemCache('./cache');
const imageHandler = createImageOptimizer({
    cache: cache,
    fallbackFormat: 'avif',
    quality: 20
});

export const handle: Handle = imageHandler;
