import {env} from '$env/dynamic/public';
import type {PageLoad} from './$types';

const apiUrl = env.PUBLIC_API_URL;

export const load: PageLoad = async ({fetch, params}) => {
	const url = `${apiUrl}/tv/shows/${params.showId}/${params.SeasonNumber}/files`;

	try {
		console.log(`Fetching data from: ${url}`);
		const response = await fetch(url, {
			method: 'GET',
			credentials: 'include'
		});

		if (!response.ok) {
			const errorText = await response.text();
			console.error(`API request failed with status ${response.status}: ${errorText}`);
			return {
				error: `Failed to load TV show files. Status: ${response.status}`,
				files: []
			};
		}

		const filesData = await response.json();
		console.log('received season_files data: ', filesData);
		return {
			files: filesData
		};
	} catch (error) {
		console.error('An error occurred while fetching TV show files:', error);
		return {
			error: `An unexpected error occurred: ${error.message || 'Unknown error'}`,
			files: []
		};
	}
};
