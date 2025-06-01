import {type ClassValue, clsx} from 'clsx';
import {twMerge} from 'tailwind-merge';
import {env} from '$env/dynamic/public';
import {goto} from '$app/navigation';
import {base} from '$app/paths';
import {toast} from 'svelte-sonner';
import {browser} from "$app/environment";

const apiUrl = browser ? env.PUBLIC_API_URL : env.PUBLIC_SSR_API_URL;

export const qualityMap: { [key: number]: string } = {
	1: '4K/UHD',
	2: '1080p/FullHD',
	3: '720p/HD',
	4: '480p/SD',
	5: 'unknown'
};
export const torrentStatusMap: { [key: number]: string } = {
	1: 'finished',
	2: 'downloading',
	3: 'error',
	4: 'unknown'
};

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function getTorrentQualityString(value: number): string {
	return qualityMap[value] || 'unknown';
}

export function getTorrentStatusString(value: number): string {
	return torrentStatusMap[value] || 'unknown';
}
export function getFullyQualifiedShowName(show: { name: string; year: number }): string {
	let name = show.name;
	if (show.year != null) {
		name += ' (' + show.year + ')';
	}
	return name;
}

export function convertTorrentSeasonRangeToIntegerRange(torrent: {
	season?: number[];
	seasons?: number[];
}): string {
	if (torrent?.season?.length === 1) return torrent.season[0]?.toString();
	if (torrent?.season?.length >= 2)
		return torrent.season[0]?.toString() + '-' + torrent.season.at(-1).toString();
	if (torrent?.seasons?.length === 1) return torrent.seasons[0]?.toString();
	if (torrent?.seasons?.length >= 2)
		return torrent.seasons[0]?.toString() + '-' + torrent.seasons.at(-1).toString();
	else {
		console.log('Error parsing season range: ' + torrent?.seasons + torrent?.season);
		return 'Error parsing season range: ' + torrent?.seasons + torrent?.season;
	}
}

export async function handleLogout() {
	const response = await fetch(apiUrl + '/auth/cookie/logout', {
		method: 'POST',
		credentials: 'include'
	});
	if (response.ok) {
		console.log('Logout successful!');
		toast.success('Logout successful!');
		await goto(base + '/login');
	} else {
		console.error('Logout failed:', response.status);
		toast.error('Logout failed: ' + response.status);
	}
}
