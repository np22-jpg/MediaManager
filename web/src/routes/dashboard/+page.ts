import {env} from '$env/dynamic/public';
import type {PageLoad} from './$types';

const apiUrl = env.PUBLIC_API_URL;
import {toast} from 'svelte-sonner';


export const load: PageLoad = async ({fetch}) => {
    const response = await fetch(apiUrl + '/tv/recommended', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include'
    });

    if (!response.ok) {
        console.error(`Failed to fetch TV recommendations: ${response.statusText}`);
        toast.error(`Failed to fetch TV recommendations: ${response.statusText}`);
        return {tvRecommendations: []};
    } else {
        console.log('Fetched TV recommendations successfully');
        toast.success('Fetched TV recommendations successfully');
    }

    const recommendations = await response.json();
    console.log('Tv Recommendations:', recommendations);
    return {tvRecommendations: recommendations};
};