export interface User {
	id: string;
	email: string;
	is_active?: boolean;
	is_superuser?: boolean;
	is_verified?: boolean;
}

export interface BearerResponse {
	access_token: string;
	token_type: string;
}

export interface BodyAuthCookieLoginAuthCookieLoginPost {
	grant_type?: string | null; // anyOf string, null, pattern ^password$ implied but not strongly typed in TS interface
	username: string;
	password: string;
	scope?: string; // default: ""
	client_id?: string | null; // anyOf string, null
	client_secret?: string | null; // anyOf string, null
}

export interface BodyAuthJwtLoginAuthJwtLoginPost {
	grant_type?: string | null; // anyOf string, null, pattern ^password$ implied but not strongly typed in TS interface
	username: string;
	password: string;
	scope?: string; // default: ""
	client_id?: string | null; // anyOf string, null
	client_secret?: string | null; // anyOf string, null
}

export interface BodyResetForgotPasswordAuthForgotPasswordPost {
	email: string; // format: email
}

export interface BodyResetResetPasswordAuthResetPasswordPost {
	token: string;
	password: string;
}

export interface BodyVerifyRequestTokenAuthRequestVerifyTokenPost {
	email: string; // format: email
}

export interface BodyVerifyVerifyAuthVerifyPost {
	token: string;
}

export interface Episode {
	number: number; // type: integer
	external_id: number; // type: integer
	title: string;
	id?: string; // type: string, format: uuid
}

export type Quality = 1 | 2 | 3 | 4 | 5;
export type TorrentStatus = 1 | 2 | 3 | 4;

// You likely want to export these maps and potentially helper functions too

export interface PublicIndexerQueryResult {
	title: string;
	quality: Quality; // $ref: #/components/schemas/Quality
	id: string; // type: string, format: uuid
	seeders: number; // type: integer
	flags: string[]; // items: { type: string }, type: array
	season: number[]; // items: { type: integer }, type: array
	size: number;
}

export interface Season {
	number: number; // type: integer
	name: string;
	overview: string;
	external_id: number; // type: integer
	episodes: Episode[]; // items: { $ref: #/components/schemas/Episode }, type: array
	id?: string; // type: string, format: uuid
}

export interface PublicSeasonFile {
	season_id: string; // type: string, format: uuid
	quality: Quality;
	torrent_id?: string; // type: string, format: uuid
	file_path_suffix?: string;
	downloaded: boolean;
}

export interface PublicSeason {
	number: number; // type: integer
	name: string;
	downloaded: boolean;
	overview: string;
	external_id: number; // type: integer
	episodes: Episode[]; // items: { $ref: #/components/schemas/Episode }, type: array
	id?: string; // type: string, format: uuid
}

export interface Show {
	name: string;
	overview: string;
	year: number; // type: integer
	external_id: number; // type: integer
	metadata_provider: string;
	seasons: Season[]; // items: { $ref: #/components/schemas/Season }, type: array
	id: string; // type: string, format: uuid
	continuous_download: boolean;
	ended: boolean;
}

export interface PublicShow {
	name: string;
	overview: string;
	year: number; // type: integer
	external_id: number; // type: integer
	metadata_provider: string;
	seasons: PublicSeason[]; // items: { $ref: #/components/schemas/Season }, type: array
	id: string; // type: string, format: uuid
	continuous_download: boolean;
	ended: boolean;
}

export interface Torrent {
	status: TorrentStatus; // $ref: #/components/schemas/TorrentStatus
	title: string;
	quality: Quality; // $ref: #/components/schemas/Quality
	imported: boolean;
	hash: string;
	id?: string; // type: string, format: uuid
}

export interface UserCreate {
	email: string; // format: email
	password: string;
	is_active?: boolean | null; // anyOf boolean, null, default: true
	is_superuser?: boolean | null; // anyOf boolean, null, default: false
	is_verified?: boolean | null; // anyOf boolean, null, default: false
}

export interface UserUpdate {
	password?: string | null; // anyOf string, null
	email?: string | null; // anyOf string, null, format: email
	is_active?: boolean | null; // anyOf boolean, null
	is_superuser?: boolean | null; // anyOf boolean, null
	is_verified?: boolean | null; // anyOf boolean, null
}

export interface MetaDataProviderShowSearchResult {
	poster_path: string | null;
	overview: string | null;
	name: string;
	external_id: number;
	year: number | null;
	metadata_provider: string;
	added: boolean;
	vote_average: number;
}
export interface RichSeasonTorrent {
	torrent_id: string;
	torrent_title: string;
	status: TorrentStatus;
	quality: Quality;
	imported: boolean;

	file_path_suffix: string;
	seasons: number[];
}

export interface RichShowTorrent {
	show_id: string;
	name: string;
	year: number | null;
	metadata_provider: string;
	torrents: RichSeasonTorrent[];
}

interface SeasonRequestBase {
	min_quality: Quality;
	wanted_quality: Quality;
}

export interface CreateSeasonRequest extends SeasonRequestBase {
	season_id: string;
}

export interface UpdateSeasonRequest extends SeasonRequestBase {
	id: string;
}

export interface SeasonRequest extends SeasonRequestBase {
	id: string;
	season: Season;
	requested_by?: User;
	authorized: boolean;
	authorized_by?: User;
	show: Show;
}
