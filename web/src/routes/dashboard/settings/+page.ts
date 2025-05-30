import {env} from '$env/dynamic/public';
import type {PageLoad} from './$types';
import {browser} from "$app/environment";

const apiUrl = browser ? env.PUBLIC_API_URL : env.PUBLIC_SSR_API_URL;
export const load: PageLoad = async ({fetch}) => {
	try {
		const users = await fetch(apiUrl + '/users/all', {
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
