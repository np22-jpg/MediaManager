package anidb

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAniDBWithMockServer(t *testing.T) {
	// Activate HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Initialize AniDB with mock server
	mockBaseURL := "http://mock-anidb-api.com/httpapi"
	InitAniDB(mockBaseURL, "testclient", "1")

	// Mock anime by ID endpoint
	animeXML := `<?xml version="1.0"?>
<anime id="1" restricted="false">
  <type>TV Series</type>
  <episodecount>13</episodecount>
  <startdate>1999-01-03</startdate>
  <enddate>1999-03-28</enddate>
  <titles>
    <title xml:lang="x-jat" type="main">Test Anime</title>
    <title xml:lang="en" type="official">Test Anime English</title>
  </titles>
  <description>A test anime for unit testing.</description>
  <ratings>
    <permanent count="100">8.5</permanent>
    <temporary count="120">8.7</temporary>
  </ratings>
  <picture>test_anime.jpg</picture>
  <url>http://test-anime.com</url>
</anime>`

	httpmock.RegisterResponder("GET", mockBaseURL,
		func(req *http.Request) (*http.Response, error) {
			query := req.URL.Query()
			request := query.Get("request")

			switch request {
			case "anime":
				return httpmock.NewStringResponse(200, animeXML), nil
			case "hotanime":
				hotAnimeXML := `<?xml version="1.0"?>
<hotanime>
  <anime id="123" restricted="false">
    <episodecount>12</episodecount>
    <startdate>2024-01-01</startdate>
    <title xml:lang="x-jat" type="main">Hot Test Anime</title>
    <ratings>
      <permanent count="200">9.0</permanent>
      <temporary count="250">9.2</temporary>
    </ratings>
    <picture>hot_anime.jpg</picture>
  </anime>
</hotanime>`
				return httpmock.NewStringResponse(200, hotAnimeXML), nil
			case "randomrecommendation":
				recommendationXML := `<?xml version="1.0"?>
<randomrecommendation>
  <recommendation>
    <anime id="456" restricted="false">
      <type>OVA</type>
      <episodecount>3</episodecount>
      <startdate>2024-06-01</startdate>
      <title xml:lang="x-jat" type="main">Recommended Test Anime</title>
      <picture>recommended_anime.jpg</picture>
      <ratings>
        <permanent count="50">7.8</permanent>
        <recommendations>5</recommendations>
      </ratings>
    </anime>
  </recommendation>
</randomrecommendation>`
				return httpmock.NewStringResponse(200, recommendationXML), nil
			case "randomsimilar":
				similarXML := `<?xml version="1.0"?>
<randomsimilar>
  <similar>
    <source aid="789" restricted="false">
      <title xml:lang="x-jat" type="main">Source Anime</title>
      <picture>source_anime.jpg</picture>
    </source>
    <target aid="987" restricted="false">
      <title xml:lang="x-jat" type="main">Similar Anime</title>
      <picture>similar_anime.jpg</picture>
    </target>
  </similar>
</randomsimilar>`
				return httpmock.NewStringResponse(200, similarXML), nil
			case "main":
				mainXML := `<?xml version="1.0"?>
<main>
  <hotanime>
    <anime id="123" restricted="false">
      <title xml:lang="x-jat" type="main">Hot Test Anime</title>
    </anime>
  </hotanime>
  <randomrecommendation>
    <recommendation>
      <anime id="456" restricted="false">
        <title xml:lang="x-jat" type="main">Recommended Test Anime</title>
      </anime>
    </recommendation>
  </randomrecommendation>
</main>`
				return httpmock.NewStringResponse(200, mainXML), nil
			default:
				return httpmock.NewStringResponse(400, "Bad Request"), nil
			}
		})

	ctx := context.Background()

	t.Run("GetAnimeByID", func(t *testing.T) {
		result, err := GetAnimeByID(ctx, 1)
		require.NoError(t, err)
		assert.NotNil(t, result)

		// Check result structure
		anime, ok := result.(AnimeInfo)
		require.True(t, ok, "Result should be AnimeInfo")

		assert.Equal(t, 1, anime.ID)
		assert.Equal(t, "TV Series", anime.Type)
		assert.Equal(t, 13, anime.EpisodeCount)
		assert.Equal(t, "1999-01-03", anime.StartDate)
		assert.Equal(t, "1999-03-28", anime.EndDate)
		assert.Contains(t, anime.Description, "test anime")
		assert.Equal(t, "test_anime.jpg", anime.Picture)
	})

	t.Run("GetHotAnime", func(t *testing.T) {
		result, err := GetHotAnime(ctx)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("GetRandomRecommendation", func(t *testing.T) {
		result, err := GetRandomRecommendation(ctx)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("GetRandomSimilar", func(t *testing.T) {
		result, err := GetRandomSimilar(ctx)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("GetMainPageData", func(t *testing.T) {
		result, err := GetMainPageData(ctx)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	// Verify expected calls were made
	assert.GreaterOrEqual(t, httpmock.GetTotalCallCount(), 5)
}

func TestAniDBErrorHandling(t *testing.T) {
	// Test with no client configured
	InitAniDB("", "", "")

	ctx := context.Background()

	t.Run("NoClientConfigured", func(t *testing.T) {
		result, err := GetAnimeByID(ctx, 1)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "client")
	})
}

func TestAniDBServerErrors(t *testing.T) {
	// Activate HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Initialize AniDB with mock server
	mockBaseURL := "http://mock-anidb-api.com/httpapi"
	InitAniDB(mockBaseURL, "testclient", "1")

	// Mock server error
	httpmock.RegisterResponder("GET", mockBaseURL,
		httpmock.NewStringResponder(500, "Internal Server Error"))

	ctx := context.Background()

	t.Run("ServerError", func(t *testing.T) {
		result, err := GetAnimeByID(ctx, 1)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestAniDBBannedError(t *testing.T) {
	// Activate HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Initialize AniDB with mock server
	mockBaseURL := "http://mock-anidb-api.com/httpapi"
	InitAniDB(mockBaseURL, "testclient", "1")

	// Mock banned error
	bannedXML := `<?xml version="1.0"?><error>Banned</error>`
	httpmock.RegisterResponder("GET", mockBaseURL,
		httpmock.NewStringResponder(200, bannedXML))

	ctx := context.Background()

	t.Run("BannedError", func(t *testing.T) {
		result, err := GetAnimeByID(ctx, 1)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Banned")
	})
}
