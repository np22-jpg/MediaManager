import {env} from '$env/dynamic/public';

export const ssr = (env.PUBLIC_WEB_SSR.toLowerCase() == 'true');
