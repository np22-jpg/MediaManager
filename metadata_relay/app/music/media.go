package music

import (
	"relay/app/music/lrclib"
	"relay/app/music/spotify"
)

// Type aliases for backward compatibility
type SpotifyClient = spotify.Client
type LRCLibClient = lrclib.Client

// Constructor functions for backward compatibility
func NewSpotify(clientID, clientSecret string) *SpotifyClient {
	return spotify.NewClient(clientID, clientSecret)
}

func NewLRCLib(baseURL string) *LRCLibClient {
	return lrclib.NewClient(baseURL)
}
