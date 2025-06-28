import { env } from '$env/dynamic/public';
import type { LayoutLoad } from './$types';

const apiUrl = env.PUBLIC_API_URL;
export const load: LayoutLoad = async ({ params, fetch }) => {
	const showId = params.showId;

	if (!showId) {
		return {
			showData: null,
			torrentsData: null
		};
	}

	try {
		const show = await fetch(`${apiUrl}/tv/shows/${showId}`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include'
		});

		const torrents = await fetch(`${apiUrl}/tv/shows/${showId}/torrents`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include'
		});

		if (!show.ok || !torrents.ok) {
			console.error(`Failed to fetch show ${showId}: ${show.statusText}`);
			return {
				showData: null,
				torrentsData: null
			};
		}

		const showData = await show.json();
		const torrentsData = await torrents.json();
		console.log('Fetched show data:', showData);
		console.log('Fetched torrents data:', torrentsData);

		return {
			showData: showData,
			torrentsData: torrentsData
		};
	} catch (error) {
		console.error('Error fetching show:', error);
		return {
			showData: null,
			torrentsData: null
		};
	}
};
