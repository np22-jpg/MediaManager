import {env} from '$env/dynamic/public';
import type {PageLoad} from './$types';
import {browser} from "$app/environment";

const apiUrl = browser ? env.PUBLIC_API_URL : env.PUBLIC_SSR_API_URL;

export const load: PageLoad = async ({fetch}) => {
	const response = await fetch(apiUrl + '/tv/shows/torrents', {
		method: 'GET',
		credentials: 'include'
	});
	return {shows: response.json()};
};
