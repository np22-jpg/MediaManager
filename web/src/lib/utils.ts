import {type ClassValue, clsx} from 'clsx';
import {twMerge} from 'tailwind-merge';

export const qualityMap: { [key: number]: string } = {
	1: 'high',
	2: 'medium',
	3: 'low',
	4: 'very_low',
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

export function convertTorrentSeasonRangeToIntegerRange(torrent: any): string {
	if (torrent.seasons.length === 1) return torrent.seasons[0]?.toString();
	if (torrent.seasons.length >= 2) return torrent.seasons[0]?.toString() + "-" + torrent.seasons.at(-1).toString();
	else {
		console.log("Error parsing season range: " + torrent.seasons);
		return "Error parsing season range: " + torrent.seasons;
	}

}