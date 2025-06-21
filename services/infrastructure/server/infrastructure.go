package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"infrastructure/config"
	"infrastructure/osm"
	pb "wildfire-risk-platform/api/proto/generated"
)

// InfrastructureServer implements the InfrastructureService gRPC server
type InfrastructureServer struct {
	pb.UnimplementedInfrastructureServiceServer

	config         *config.Config
	overpassClient *osm.OverpassClient
}

// NewInfrastructureServer creates a new infrastructure server
func NewInfrastructureServer(cfg *config.Config) *InfrastructureServer {

	overpassClient := osm.NewOverpassClient(
		cfg.OverpassAPIURL,
		cfg.HTTPTimeout,
		cfg.MaxRetries,
		cfg.RetryDelay,
	)

	return &InfrastructureServer{
		config:         cfg,
		overpassClient: overpassClient,
	}
}

// GetAssetsInAOI retrieves infrastructure assets within the area of interest
func (s *InfrastructureServer) GetAssetsInAOI(ctx context.Context, req *pb.GetAssetsRequest) (*pb.GetAssetsResponse, error) {

	if req.AoiGeojson == "" {
		return nil, status.Error(codes.InvalidArgument, "aoi_geojson is required")
	}

	// Validate GeoJSON structure
	var geojson map[string]interface{}
	if err := json.Unmarshal([]byte(req.AoiGeojson), &geojson); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid GeoJSON format")
	}

	log.Printf("Fetching assets for AOI with bounds")

	assets, err := s.overpassClient.QueryAssets(ctx, req.AoiGeojson)
	if err != nil {
		log.Printf("Failed to query Overpass API: %v", err)
		return nil, status.Error(codes.Internal, "failed to fetch assets from OpenStreetMap")
	}

	log.Printf("Retrieved %d assets from Overpass API", len(assets))

	pbAssets := make([]*pb.Asset, 0, len(assets))

	for _, asset := range assets {

		geometryJSON, err := json.Marshal(asset.Geometry)
		if err != nil {
			log.Printf("Failed to marshal geometry for asset: %v", err)
			continue // Skip this asset
		}

		properties := make(map[string]string)
		for k, v := range asset.Properties {

			properties[k] = fmt.Sprintf("%v", v)
		}

		pbAsset := &pb.Asset{
			AssetType:            asset.Type,
			AssetGeometryGeojson: string(geometryJSON),
			Properties:           properties,
		}

		pbAssets = append(pbAssets, pbAsset)
	}

	typeCount := make(map[string]int)
	for _, asset := range pbAssets {
		typeCount[asset.AssetType]++
	}
	log.Printf("Asset distribution: %v", typeCount)

	return &pb.GetAssetsResponse{
		Assets:     pbAssets,
		TotalCount: int32(len(pbAssets)),
	}, nil
}
