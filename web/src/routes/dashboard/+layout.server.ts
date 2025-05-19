import {env} from '$env/dynamic/public';
import type {LayoutServerLoad} from './$types';
import {redirect} from '@sveltejs/kit';
import {base} from "$app/paths";

const apiUrl = env.PUBLIC_API_URL;

export const load: LayoutServerLoad = async ({fetch}) => {
	const response = await fetch(apiUrl + '/users/me', {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		},
		credentials: 'include'
	});
	if (!response.ok) {
		console.log('unauthorized, redirecting to login');
		throw redirect(303, base + '/login');
	}
	return {user: await response.json()};
};
