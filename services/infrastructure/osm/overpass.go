package osm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	geojson "github.com/paulmach/go.geojson"
)

type OverpassClient struct {
	apiURL     string
	httpClient *http.Client
	maxRetries int
	retryDelay time.Duration
}

func NewOverpassClient(apiURL string, timeout time.Duration, maxRetries int, retryDelay time.Duration) *OverpassClient {
	return &OverpassClient{
		apiURL: apiURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

type Asset struct {
	Type       string                 `json:"type"`
	Geometry   *geojson.Geometry      `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

func (c *OverpassClient) QueryAssets(ctx context.Context, areaGeoJSON string) ([]Asset, error) {
	bounds, err := extractBounds(areaGeoJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to extract bounds: %w", err)
	}

	query := c.buildQuery(bounds)

	var response *OverpassResponse
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("Retrying Overpass API query (attempt %d/%d)", attempt, c.maxRetries)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.retryDelay):
			}
		}

		response, err = c.executeQuery(ctx, query)
		if err == nil {
			break
		}

		log.Printf("Overpass API error: %v", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query Overpass API after %d attempts: %w", c.maxRetries+1, err)
	}

	return c.convertToAssets(response), nil
}

type Bounds struct {
	MinLat, MinLon, MaxLat, MaxLon float64
}

func extractBounds(geoJSONStr string) (*Bounds, error) {
	var geoData map[string]interface{}
	if err := json.Unmarshal([]byte(geoJSONStr), &geoData); err != nil {
		return nil, fmt.Errorf("invalid GeoJSON: %w", err)
	}

	coords, ok := geoData["coordinates"].([]interface{})
	if !ok || len(coords) == 0 {
		return nil, fmt.Errorf("no coordinates found in GeoJSON")
	}

	ring, ok := coords[0].([]interface{})
	if !ok || len(ring) == 0 {
		return nil, fmt.Errorf("invalid polygon structure")
	}

	firstCoord, ok := ring[0].([]interface{})
	if !ok || len(firstCoord) < 2 {
		return nil, fmt.Errorf("invalid coordinate format")
	}

	minLon, _ := firstCoord[0].(float64)
	minLat, _ := firstCoord[1].(float64)
	maxLon, maxLat := minLon, minLat

	for _, coord := range ring {
		c, ok := coord.([]interface{})
		if !ok || len(c) < 2 {
			continue
		}
		lon, _ := c[0].(float64)
		lat, _ := c[1].(float64)

		if lon < minLon {
			minLon = lon
		}
		if lon > maxLon {
			maxLon = lon
		}
		if lat < minLat {
			minLat = lat
		}
		if lat > maxLat {
			maxLat = lat
		}
	}

	return &Bounds{
		MinLat: minLat,
		MinLon: minLon,
		MaxLat: maxLat,
		MaxLon: maxLon,
	}, nil
}

func (c *OverpassClient) buildQuery(bounds *Bounds) string {
	// Build a query that fetches:
	// - Buildings (way["building"])
	// - Roads (way["highway"])
	// - Power lines (way["power"="line"])
	// - Railways (way["railway"])

	bbox := fmt.Sprintf("%f,%f,%f,%f", bounds.MinLat, bounds.MinLon, bounds.MaxLat, bounds.MaxLon)

	query := fmt.Sprintf(`
		[out:json][timeout:60];
		(
			// Buildings
			way["building"](%s);
			relation["building"](%s);
			
			// Roads
			way["highway"](%s);
			
			// Power infrastructure
			way["power"="line"](%s);
			node["power"="tower"](%s);
			node["power"="pole"](%s);
			
			// Railways
			way["railway"](%s);
			
			// Critical facilities
			node["amenity"="hospital"](%s);
			node["amenity"="fire_station"](%s);
			way["amenity"="hospital"](%s);
			way["amenity"="fire_station"](%s);
		);
		out geom;
	`, bbox, bbox, bbox, bbox, bbox, bbox, bbox, bbox, bbox, bbox, bbox)

	return query
}

type OverpassResponse struct {
	Elements []OSMElement `json:"elements"`
}

type OSMElement struct {
	Type     string            `json:"type"`
	ID       int64             `json:"id"`
	Lat      float64           `json:"lat,omitempty"`
	Lon      float64           `json:"lon,omitempty"`
	Tags     map[string]string `json:"tags"`
	Geometry []OSMNode         `json:"geometry,omitempty"`
}

// OSMNode represents a node in OSM geometry
type OSMNode struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// executeQuery executes the Overpass query
func (c *OverpassClient) executeQuery(ctx context.Context, query string) (*OverpassResponse, error) {

	req, err := http.NewRequestWithContext(ctx, "POST", c.apiURL, strings.NewReader(query))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Overpass API returned status %d: %s", resp.StatusCode, string(body))
	}

	var overpassResp OverpassResponse
	if err := json.NewDecoder(resp.Body).Decode(&overpassResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &overpassResp, nil
}

func (c *OverpassClient) convertToAssets(response *OverpassResponse) []Asset {
	assets := make([]Asset, 0, len(response.Elements))

	for _, element := range response.Elements {
		asset := c.convertElement(element)
		if asset != nil {
			assets = append(assets, *asset)
		}
	}

	return assets
}

// convertElement converts a single OSM element to an Asset
func (c *OverpassClient) convertElement(element OSMElement) *Asset {

	assetType := c.determineAssetType(element.Tags)
	if assetType == "" {
		return nil
	}

	var geometry *geojson.Geometry

	switch element.Type {
	case "node":

		geometry = geojson.NewPointGeometry([]float64{element.Lon, element.Lat})

	case "way":

		if len(element.Geometry) < 2 {
			return nil // Invalid geometry
		}

		coords := make([][]float64, len(element.Geometry))
		for i, node := range element.Geometry {
			coords[i] = []float64{node.Lon, node.Lat}
		}

		if element.Tags["building"] != "" && len(coords) > 3 &&
			coords[0][0] == coords[len(coords)-1][0] &&
			coords[0][1] == coords[len(coords)-1][1] {
			// It's a polygon (building)
			geometry = geojson.NewPolygonGeometry([][][]float64{coords})
		} else {
			// It's a line (road, power line, etc.)
			geometry = geojson.NewLineStringGeometry(coords)
		}

	case "relation":
		return nil

	default:
		return nil
	}

	properties := make(map[string]interface{})
	for k, v := range element.Tags {
		properties[k] = v
	}
	properties["osm_id"] = element.ID
	properties["osm_type"] = element.Type

	return &Asset{
		Type:       assetType,
		Geometry:   geometry,
		Properties: properties,
	}
}

// determineAssetType determines the asset type from OSM tags
func (c *OverpassClient) determineAssetType(tags map[string]string) string {

	if tags["building"] != "" {
		return "building"
	}
	if tags["highway"] != "" {
		return "road"
	}
	if tags["power"] == "line" {
		return "power_line"
	}
	if tags["power"] == "tower" || tags["power"] == "pole" {
		return "power_infrastructure"
	}
	if tags["railway"] != "" {
		return "railway"
	}
	if tags["amenity"] == "hospital" {
		return "hospital"
	}
	if tags["amenity"] == "fire_station" {
		return "fire_station"
	}

	return ""
}
