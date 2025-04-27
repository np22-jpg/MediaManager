import {env} from "$env/dynamic/public";
import type {LayoutLoad} from "./$types";

let apiUrl = env.PUBLIC_API_URL

export const load: LayoutLoad = async ({fetch}) => {
    const response = await fetch(apiUrl + "/users/me", {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
        credentials: 'include',
    });
    return {user: await response.json()};
}

