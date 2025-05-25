import {env} from '$env/dynamic/public';
import type {PageLoad} from './$types';

export const load: PageLoad = async ({fetch}) => {
	try {
		const users = await fetch(env.PUBLIC_API_URL + '/users/all', {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include'
		});

		if (!users.ok) {
			console.error(`Failed to fetch users: ${users.statusText}`);
			return {
				users: null
			};
		}

		const usersData = await users.json();
		console.log('Fetched users:', usersData);

		return {
			users: usersData
		};
	} catch (error) {
		console.error('Error fetching users:', error);
		return {
			users: null
		};
	}
};
