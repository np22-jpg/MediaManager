package anilist

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"relay/app/common"
	"relay/app/metrics"
)

// GetURLParam extracts a URL parameter from the request path
func GetURLParam(r *http.Request, key string) string {
	path := strings.TrimPrefix(r.URL.Path, "/anilist/")
	parts := strings.Split(path, "/")

	switch key {
	case "id":
		if len(parts) >= 2 {
			return parts[1]
		}
	}
	return ""
}

// GetMediaByIDHandler handles AniList media lookup by ID endpoint
// GET /anilist/anime/{id} or /anilist/manga/{id}
func GetMediaByIDHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling AniList get media by ID route")

	mediaIDStr := GetURLParam(r, "id")
	mediaID, err := strconv.Atoi(mediaIDStr)
	if err != nil {
		common.WriteErrorResponse(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	// Record metrics with timer
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		status := "200"
		if err != nil {
			status = "500"
		}
		metrics.RecordAPIRequest("anilist", "media_by_id", status, duration)
	}()

	media, err := GetMediaByID(r.Context(), mediaID)
	if err != nil {
		slog.Error("failed to get AniList media", "error", err, "id", mediaID)
		common.WriteErrorResponse(w, "Failed to fetch media", http.StatusInternalServerError)
		return
	}

	if media == nil {
		common.WriteErrorResponse(w, "Media not found", http.StatusNotFound)
		return
	}

	// Transform to REST-friendly format
	result := transformMediaToREST(media)
	common.WriteJSONResponse(w, result)
}

// SearchMediaHandler handles AniList media search endpoint
// GET /anilist/search?q=query&type=anime|manga&page=1&per_page=20
func SearchMediaHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling AniList search route")

	query := r.URL.Query().Get("q")
	if query == "" {
		common.WriteErrorResponse(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	mediaType := r.URL.Query().Get("type")
	if mediaType != "" {
		mediaType = strings.ToUpper(mediaType)
		if mediaType != "ANIME" && mediaType != "MANGA" {
			common.WriteErrorResponse(w, "Type must be 'anime' or 'manga'", http.StatusBadRequest)
			return
		}
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	perPage := 20
	if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 50 {
			perPage = pp
		}
	}

	// Record metrics with timer
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		status := "200"
		metrics.RecordAPIRequest("anilist", "search", status, duration)
	}()

	pageResult, err := SearchMedia(r.Context(), query, mediaType, page, perPage)
	if err != nil {
		slog.Error("failed to search AniList media", "error", err, "query", query)
		common.WriteErrorResponse(w, "Failed to search media", http.StatusInternalServerError)
		return
	}

	// Transform to REST-friendly format
	result := transformPageToREST(pageResult)
	common.WriteJSONResponse(w, result)
}

// GetTrendingAnimeHandler handles trending anime endpoint
// GET /anilist/trending/anime?page=1&per_page=20
func GetTrendingAnimeHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling AniList trending anime route")

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	perPage := 20
	if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 50 {
			perPage = pp
		}
	}

	// Record metrics with timer
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		status := "200"
		metrics.RecordAPIRequest("anilist", "trending", status, duration)
	}()

	pageResult, err := GetTrendingAnime(r.Context(), page, perPage)
	if err != nil {
		slog.Error("failed to get trending anime", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch trending anime", http.StatusInternalServerError)
		return
	}

	// Transform to REST-friendly format
	result := transformPageToREST(pageResult)
	common.WriteJSONResponse(w, result)
}

// GetSeasonalAnimeHandler handles seasonal anime endpoint
// GET /anilist/seasonal?year=2024&season=winter&page=1&per_page=20
func GetSeasonalAnimeHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handling AniList seasonal anime route")

	year := 0
	if yearStr := r.URL.Query().Get("year"); yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil && y > 1950 && y <= 2030 {
			year = y
		}
	}

	season := r.URL.Query().Get("season")
	if season != "" {
		season = strings.ToUpper(season)
		if season != "WINTER" && season != "SPRING" && season != "SUMMER" && season != "FALL" {
			common.WriteErrorResponse(w, "Season must be 'winter', 'spring', 'summer', or 'fall'", http.StatusBadRequest)
			return
		}
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	perPage := 20
	if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 50 {
			perPage = pp
		}
	}

	// Record metrics with timer
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		status := "200"
		metrics.RecordAPIRequest("anilist", "seasonal", status, duration)
	}()

	pageResult, err := GetSeasonalAnime(r.Context(), year, season, page, perPage)
	if err != nil {
		slog.Error("failed to get seasonal anime", "error", err)
		common.WriteErrorResponse(w, "Failed to fetch seasonal anime", http.StatusInternalServerError)
		return
	}

	// Transform to REST-friendly format
	result := transformPageToREST(pageResult)
	common.WriteJSONResponse(w, result)
}

// Helper functions to transform GraphQL responses to REST-friendly format

func transformMediaToREST(media *Media) map[string]interface{} {
	result := map[string]interface{}{
		"id":           media.ID,
		"title":        transformTitle(media.Title),
		"type":         strings.ToLower(media.Type),
		"format":       strings.ToLower(media.Format),
		"status":       strings.ToLower(media.Status),
		"description":  media.Description,
		"genres":       media.Genres,
		"score":        media.AverageScore,
		"popularity":   media.Popularity,
		"favourites":   media.Favourites,
		"is_adult":     media.IsAdult,
		"cover_image":  transformCoverImage(media.CoverImage),
		"banner_image": media.BannerImage,
	}

	// Add type-specific fields
	switch media.Type {
	case "ANIME":
		result["episodes"] = media.Episodes
		result["duration"] = media.Duration
		result["season"] = strings.ToLower(media.Season)
		result["season_year"] = media.SeasonYear
		if media.NextAiringEpisode != nil {
			result["next_airing_episode"] = transformAiringSchedule(media.NextAiringEpisode)
		}
	case "MANGA":
		result["chapters"] = media.Chapters
		result["volumes"] = media.Volumes
	}

	// Add dates
	if media.StartDate != nil {
		result["start_date"] = transformFuzzyDate(media.StartDate)
	}
	if media.EndDate != nil {
		result["end_date"] = transformFuzzyDate(media.EndDate)
	}

	// Add studios
	if len(media.Studios.Nodes) > 0 {
		result["studios"] = transformStudios(media.Studios.Nodes)
	}

	// Add trailer
	if media.Trailer != nil {
		result["trailer"] = transformTrailer(media.Trailer)
	}

	// Add external links
	if len(media.ExternalLinks) > 0 {
		result["external_links"] = transformExternalLinks(media.ExternalLinks)
	}

	// Add tags
	if len(media.Tags) > 0 {
		result["tags"] = transformTags(media.Tags)
	}

	return result
}

func transformPageToREST(page *Page) map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(page.Media))
	for _, media := range page.Media {
		results = append(results, transformMediaToREST(&media))
	}

	return map[string]interface{}{
		"results": results,
		"pagination": map[string]interface{}{
			"total":         page.PageInfo.Total,
			"current_page":  page.PageInfo.CurrentPage,
			"last_page":     page.PageInfo.LastPage,
			"has_next_page": page.PageInfo.HasNextPage,
			"per_page":      page.PageInfo.PerPage,
		},
	}
}

func transformTitle(title MediaTitle) map[string]interface{} {
	return map[string]interface{}{
		"romaji":  title.Romaji,
		"english": title.English,
		"native":  title.Native,
	}
}

func transformCoverImage(cover CoverImage) map[string]interface{} {
	return map[string]interface{}{
		"extra_large": cover.ExtraLarge,
		"large":       cover.Large,
		"medium":      cover.Medium,
		"color":       cover.Color,
	}
}

func transformFuzzyDate(date *FuzzyDate) map[string]interface{} {
	return map[string]interface{}{
		"year":  date.Year,
		"month": date.Month,
		"day":   date.Day,
	}
}

func transformStudios(studios []Studio) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(studios))
	for _, studio := range studios {
		result = append(result, map[string]interface{}{
			"id":                  studio.ID,
			"name":                studio.Name,
			"is_animation_studio": studio.IsAnimationStudio,
		})
	}
	return result
}

func transformTrailer(trailer *MediaTrailer) map[string]interface{} {
	return map[string]interface{}{
		"id":        trailer.ID,
		"site":      trailer.Site,
		"thumbnail": trailer.Thumbnail,
	}
}

func transformExternalLinks(links []ExternalLink) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(links))
	for _, link := range links {
		result = append(result, map[string]interface{}{
			"id":       link.ID,
			"url":      link.URL,
			"site":     link.Site,
			"type":     strings.ToLower(link.Type),
			"language": link.Language,
			"color":    link.Color,
			"icon":     link.Icon,
		})
	}
	return result
}

func transformTags(tags []MediaTag) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(tags))
	for _, tag := range tags {
		result = append(result, map[string]interface{}{
			"id":          tag.ID,
			"name":        tag.Name,
			"description": tag.Description,
			"category":    tag.Category,
			"rank":        tag.Rank,
			"is_adult":    tag.IsAdult,
		})
	}
	return result
}

func transformAiringSchedule(schedule *AiringSchedule) map[string]interface{} {
	return map[string]interface{}{
		"id":                schedule.ID,
		"airing_at":         schedule.AiringAt,
		"time_until_airing": schedule.TimeUntilAiring,
		"episode":           schedule.Episode,
		"media_id":          schedule.MediaID,
	}
}
