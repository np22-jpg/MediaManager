import {env} from '$env/dynamic/public';
import type {LayoutLoad} from './$types';
import {redirect} from '@sveltejs/kit';
import {base} from '$app/paths';
import {browser} from '$app/environment';
import {goto} from '$app/navigation';

const apiUrl = env.PUBLIC_API_URL;

export const load: LayoutLoad = async ({fetch}) => {
	const response = await fetch(apiUrl + '/users/me', {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		},
		credentials: 'include'
	});
	if (!response.ok) {
		console.log('unauthorized, redirecting to login');
		if (browser) {
			await goto(base + '/login');
		} else {
			throw redirect(303, base + '/login');
		}
	}
	return {user: await response.json()};
};
