package utils

import (
	"encoding/json"
	"fmt"
)

type GeoJSON struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}

func ValidateGeoJSONPolygon(geojsonStr string) error {
	var geojson GeoJSON
	if err := json.Unmarshal([]byte(geojsonStr), &geojson); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if geojson.Type != "Polygon" {
		return fmt.Errorf("expected Polygon type, got %s", geojson.Type)
	}

	return nil
}

func ConvertGeoJSONToWKT(geojsonStr string) (string, error) {

	var geojson GeoJSON
	if err := json.Unmarshal([]byte(geojsonStr), &geojson); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	return geojsonStr, nil
}
