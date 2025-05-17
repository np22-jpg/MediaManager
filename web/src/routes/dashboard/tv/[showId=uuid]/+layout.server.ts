import {env} from '$env/dynamic/public';
import type {LayoutServerLoad} from './$types';

export const load: LayoutServerLoad = async ({params, fetch}) => {
	const showId = params.showId;

	if (!showId) {
		return {
			showData: null,
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

		if (!response.ok) {
			console.error(`Failed to fetch show ${showId}: ${response.statusText}`);
			return {
				showData: null,
				error: `Failed to load show: ${response.statusText}`
			};
		}

		const showData = await response.json();
		console.log('Fetched show data:', showData);

		return {
			showData: showData
		};
	} catch (error) {
		console.error('Error fetching show:', error);
		return {
			showData: null,
			error: 'An error occurred while fetching show data.'
		};
	}
};
