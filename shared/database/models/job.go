package models

import (
	"time"

	"github.com/google/uuid"
)

type JobStatus string

const (
	JobStatusPending       JobStatus = "PENDING"
	JobStatusGatheringData JobStatus = "GATHERING_DATA"
	JobStatusProcessing    JobStatus = "PROCESSING"
	JobStatusComplete      JobStatus = "COMPLETE"
	JobStatusFailed        JobStatus = "FAILED"
)

type RiskAssessmentJob struct {
	JobID      uuid.UUID `json:"job_id" db:"job_id"`
	AOIPolygon string    `json:"aoi_polygon" db:"aoi_polygon"`
	JobStatus  JobStatus `json:"job_status" db:"job_status"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type InfrastructureAsset struct {
	AssetID       int64                  `json:"asset_id" db:"asset_id"`
	JobID         uuid.UUID              `json:"job_id" db:"job_id"`
	AssetType     string                 `json:"asset_type" db:"asset_type"`
	AssetGeometry string                 `json:"asset_geometry" db:"asset_geometry"`
	Properties    map[string]interface{} `json:"properties" db:"properties"`
}

type AssetRiskAnalysis struct {
	AnalysisID         int64     `json:"analysis_id" db:"analysis_id"`
	AssetID            int64     `json:"asset_id" db:"asset_id"`
	RiskFromVegetation float64   `json:"risk_from_vegetation" db:"risk_from_vegetation"`
	RiskFromSlope      float64   `json:"risk_from_slope" db:"risk_from_slope"`
	RiskFromWind       float64   `json:"risk_from_wind" db:"risk_from_wind"`
	OverallRiskScore   float64   `json:"overall_risk_score" db:"overall_risk_score"`
	AnalysisTimestamp  time.Time `json:"analysis_timestamp" db:"analysis_timestamp"`
}
