import {env} from '$env/dynamic/public';

let ssrMode = false
if (env.PUBLIC_WEB_SSR == undefined) {
    ssrMode = false;
} else {
    ssrMode = env.PUBLIC_WEB_SSR.toLowerCase() == 'true';
}
console.log('SSR Mode:', ssrMode);
export const ssr = ssrMode;