import {env} from '$env/dynamic/public';
import type {PageLoad} from './$types';

const apiUrl = env.PUBLIC_API_URL;

export const load: PageLoad = async ({fetch}) => {

};
