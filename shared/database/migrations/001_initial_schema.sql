-- Enable PostGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;

-- Main table to track user-initiated risk assessment jobs
CREATE TABLE risk_assessment_jobs (
    job_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aoi_polygon GEOMETRY(Polygon, 4326) NOT NULL,
    job_status VARCHAR(50) NOT NULL DEFAULT 'PENDING', -- PENDING, GATHERING_DATA, PROCESSING, COMPLETE, FAILED
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Stores the assets (buildings, etc.) identified within the AOI for a specific job
CREATE TABLE infrastructure_assets (
    asset_id BIGSERIAL PRIMARY KEY,
    job_id UUID NOT NULL REFERENCES risk_assessment_jobs(job_id) ON DELETE CASCADE,
    asset_type VARCHAR(50) NOT NULL, -- 'building', 'road', 'power_line'
    -- Using a more generic GEOMETRY type to accommodate points, lines, or polygons from OSM
    asset_geometry GEOMETRY(Geometry, 4326) NOT NULL,
    properties JSONB -- Store additional OSM properties
);

-- The final output table, storing time-series risk scores for each asset
CREATE TABLE asset_risk_analysis (
    analysis_id BIGSERIAL PRIMARY KEY,
    asset_id BIGINT NOT NULL REFERENCES infrastructure_assets(asset_id) ON DELETE CASCADE,
    -- Individual risk components (normalized to 0.0 - 1.0) for detailed analysis
    risk_from_vegetation NUMERIC(4, 3) NOT NULL,
    risk_from_slope NUMERIC(4, 3) NOT NULL,
    risk_from_wind NUMERIC(4, 3) NOT NULL,
    -- The final, weighted risk score, for quick filtering and visualization
    overall_risk_score NUMERIC(4, 3) NOT NULL,
    analysis_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes are crucial for performance when joining tables or querying by location
CREATE INDEX idx_assets_job_id ON infrastructure_assets(job_id);
CREATE INDEX idx_risk_asset_id ON asset_risk_analysis(asset_id);
CREATE INDEX idx_jobs_status ON risk_assessment_jobs(job_status);
CREATE INDEX idx_jobs_created ON risk_assessment_jobs(created_at);

-- A spatial index is critical for any location-based queries
CREATE INDEX idx_assets_spatial ON infrastructure_assets USING GIST (asset_geometry);
CREATE INDEX idx_jobs_spatial ON risk_assessment_jobs USING GIST (aoi_polygon);

-- Function to automatically update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to automatically update updated_at on risk_assessment_jobs
CREATE TRIGGER update_risk_assessment_jobs_updated_at 
    BEFORE UPDATE ON risk_assessment_jobs 
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();