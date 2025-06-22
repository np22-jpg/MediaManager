import {env} from '$env/dynamic/public';
import type {PageLoad} from './$types';

const apiUrl = env.PUBLIC_API_URL;

export const load: PageLoad = async ({fetch}) => {
	const response = await fetch(apiUrl + '/tv/shows/torrents', {
		method: 'GET',
		credentials: 'include'
	});
    return {shows: response.json()};
};
