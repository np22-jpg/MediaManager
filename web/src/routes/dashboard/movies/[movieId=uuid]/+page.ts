import type {LayoutLoad} from './$types';
import {PUBLIC_API_URL} from '$env/static/public';
import {error} from '@sveltejs/kit';

export const load: LayoutLoad = async ({params, fetch}) => {
    const res = await fetch(`${PUBLIC_API_URL}/movies/${params.movieId}`, {
        credentials: 'include'
    });
    if (!res.ok) throw error(res.status, `Failed to load movie`);
    const movieData = await res.json();
    return {movie: movieData, torrents: []};
};
