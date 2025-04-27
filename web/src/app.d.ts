// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
		interface User {
			id: string;
			email: string;
			is_active?: boolean;
			is_superuser?: boolean;
			is_verified?: boolean;
		}
	}
}

export {};
