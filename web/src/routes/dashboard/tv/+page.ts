import { env } from '$env/dynamic/public';
import type { PageLoad } from './$types';

const apiUrl = env.PUBLIC_API_URL;

export const load: PageLoad = async ({ fetch }) => {
	const response = fetch(apiUrl + '/tv/shows', {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		},
		credentials: 'include'
	});

	return { tvShows: response };
};
