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

export function getQualityString(value: number): string {
	return qualityMap[value] || 'unknown';
}

export function getTorrentStatusString(value: number): string {
	return torrentStatusMap[value] || 'unknown';
}
