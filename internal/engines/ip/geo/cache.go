package geo

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"osint/internal/models"
)

const cacheTTL = 24 * time.Hour

type cacheEntry struct {
	Result    GeoResult `json:"result"`
	ExpiresAt time.Time `json:"expires_at"`
}

var cacheMu sync.Mutex

func getCached(
	ip string,
) (GeoResult, bool) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	cache, err := readCache()
	if err != nil {
		return GeoResult{}, false
	}

	entry, exists := cache[ip]
	if !exists {
		return GeoResult{}, false
	}

	if time.Now().After(
		entry.ExpiresAt,
	) {
		delete(
			cache,
			ip,
		)

		_ = writeCache(cache)

		return GeoResult{}, false
	}

	return entry.Result, true
}

func setCached(
	ip string,
	result GeoResult,
) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	cache, err := readCache()
	if err != nil {
		cache = make(
			map[string]cacheEntry,
		)
	}

	cache[ip] = cacheEntry{
		Result: result,

		ExpiresAt: time.Now().Add(
			cacheTTL,
		),
	}

	_ = writeCache(cache)
}

func readCache() (
	map[string]cacheEntry,
	error,
) {
	path, err := models.GeoCachePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(
				map[string]cacheEntry,
			), nil
		}

		return nil, err
	}

	if len(data) == 0 {
		return make(
			map[string]cacheEntry,
		), nil
	}

	var cache map[string]cacheEntry

	if err := json.Unmarshal(
		data,
		&cache,
	); err != nil {
		return nil, err
	}

	if cache == nil {
		cache = make(
			map[string]cacheEntry,
		)
	}

	return cache, nil
}

func writeCache(
	cache map[string]cacheEntry,
) error {
	path, err := models.GeoCachePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(
		cache,
		"",
		"  ",
	)
	if err != nil {
		return err
	}

	tmpPath := path + ".tmp"

	if err := os.WriteFile(
		tmpPath,
		data,
		0644,
	); err != nil {
		return err
	}

	if err := os.Rename(
		tmpPath,
		path,
	); err != nil {
		_ = os.Remove(tmpPath)

		return err
	}

	return nil
}
