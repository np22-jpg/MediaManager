import {env} from '$env/dynamic/public';
import type {LayoutLoad} from './$types';

export const load: LayoutLoad = async ({params, fetch}) => {
	const showId = params.showId;

	if (!showId) {
		return {
			showData: null,
			torrentsData: null,
			error: 'Show ID is missing'
		};
	}

	try {
		const response = await fetch(`${env.PUBLIC_API_URL}/tv/shows/${showId}`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include'
		});

		const torrents = await fetch(`${env.PUBLIC_API_URL}/tv/shows/${showId}/torrents`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include'
		});

		if (!response.ok || !torrents.ok) {
			console.error(`Failed to fetch show ${showId}: ${response.statusText}`);
			return {
				showData: null,
				torrentsData: null,
				error: `Failed to load show or/and its torrents: ${response.statusText}`
			};
		}

		const showData = await response.json();
		const torrentsData = await torrents.json();
		console.log('Fetched show data:', showData);
		console.log('Fetched torrents data:', torrentsData);

		return {
			showData: showData,
			torrentsData: torrentsData,
		};
	} catch (error) {
		console.error('Error fetching show:', error);
		return {
			showData: null,
			torrentsData: null,
			error: 'An error occurred while fetching show data.'
		};
	}
};
