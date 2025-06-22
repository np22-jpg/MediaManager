import {env} from '$env/dynamic/public';
import type {PageLoad} from './$types';

const apiUrl = env.PUBLIC_API_URL;

export const load: PageLoad = async ({fetch, params}) => {
	const url = `${apiUrl}/tv/seasons/${params.SeasonId}/files`;
	const url2 = `${apiUrl}/tv/seasons/${params.SeasonId}`;

	try {
		console.log(`Fetching data from: ${url} and ${url2}`);
		const response = await fetch(url, {
			method: 'GET',
			credentials: 'include'
		});
		const response2 = await fetch(url2, {
			method: 'GET',
			credentials: 'include'
		});

		if (!response.ok) {
			const errorText = await response.text();
			console.error(`API request failed with status ${response.status}: ${errorText}`);
		}

		if (!response2.ok) {
			const errorText = await response.text();
			console.error(`API request failed with status ${response.status}: ${errorText}`);
		}

		const filesData = await response.json();
		const seasonData = await response2.json();
		console.log('received season_files data: ', filesData);
		console.log('received season data: ', seasonData);
		return {
			files: filesData,
			season: seasonData
		};
	} catch (error) {
		console.error('An error occurred while fetching TV show files:', error);
		return {
			error: `An unexpected error occurred: ${error.message || 'Unknown error'}`,
			files: [],
			season: null
		};
	}
};
