import {env} from '$env/dynamic/public';
import type {LayoutLoad} from './$types';

const apiUrl = env.PUBLIC_API_URL;
export const load: LayoutLoad = async ({fetch}) => {
    try {
        const requests = await fetch(`${apiUrl}/movies/requests`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include'
        });

        if (!requests.ok) {
            console.error(`Failed to fetch season requests ${requests.statusText}`);
            return {
                requestsData: null
            };
        }

        const requestsData = await requests.json();
        console.log('Fetched season requests:', requestsData);

        return {
            requestsData: requestsData
        };
    } catch (error) {
        console.error('Error fetching season requests:', error);
        return {
            requestsData: null
        };
    }
};
