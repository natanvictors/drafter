// Package ddragon fetches static game data from Riot's Data Dragon CDN.
package ddragon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/natanvictors/drafter/internal/parser"
)

const baseURL = "https://ddragon.leagueoflegends.com"

// ChampionNames returns the display name of every champion in the latest
// Data Dragon version.
func ChampionNames() ([]string, error) {
	version, err := latestVersion()
	if err != nil {
		return nil, err
	}

	data, err := get(fmt.Sprintf("%s/cdn/%s/data/en_US/champion.json", baseURL, version))
	if err != nil {
		return nil, fmt.Errorf("failed to get champion API: %w", err)
	}

	body, err := parser.ParseChampion(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data: %w", err)
	}

	names := make([]string, 0, len(body.Data))
	for key, raw := range body.Data {
		var champ parser.Champion
		if err := json.Unmarshal(raw, &champ); err != nil {
			return nil, fmt.Errorf("failed to decode champion %q: %w", key, err)
		}
		names = append(names, champ.Name)
	}
	return names, nil
}

func latestVersion() (string, error) {
	data, err := get(baseURL + "/api/versions.json")
	if err != nil {
		return "", fmt.Errorf("failed to get version: %w", err)
	}

	var versions []string
	if err := json.Unmarshal(data, &versions); err != nil {
		return "", fmt.Errorf("failed to decode versions API: %w", err)
	}
	if len(versions) == 0 {
		return "", fmt.Errorf("versions API returned no versions")
	}
	return versions[0], nil
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
