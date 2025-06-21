package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"api-gateway/config"
	"wildfire-risk-platform/api/proto/generated"
)

type RiskAssessmentHandler struct {
	config             *config.Config
	orchestratorConn   *grpc.ClientConn
	orchestratorClient generated.OrchestratorServiceClient
	connMutex          sync.Mutex
	connInitialized    bool
}

// Request/Response structures for REST API
type CreateJobRequest struct {
	AOI AOIGeometry `json:"aoi" binding:"required"`
}

type AOIGeometry struct {
	Type        string          `json:"type" binding:"required"`
	Coordinates json.RawMessage `json:"coordinates" binding:"required"`
}

type CreateJobResponse struct {
	JobID   string `json:"job_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type JobStatusResponse struct {
	JobID     string `json:"job_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewRiskAssessmentHandler(cfg *config.Config) *RiskAssessmentHandler {
	// Lazy start up to handle orchestrator down error
	return &RiskAssessmentHandler{
		config: cfg,
	}
}

func (h *RiskAssessmentHandler) ensureGRPCConnection() error {
	h.connMutex.Lock()
	defer h.connMutex.Unlock()

	if h.connInitialized {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Attempting to connect to orchestrator at %s", h.config.OrchestratorURL)

	conn, err := grpc.DialContext(
		ctx,
		h.config.OrchestratorURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}

	h.orchestratorConn = conn
	h.orchestratorClient = generated.NewOrchestratorServiceClient(conn)
	h.connInitialized = true

	log.Printf("Successfully connected to orchestrator at %s", h.config.OrchestratorURL)
	return nil
}

// POST /api/v1/wildfire-risk-jobs
func (h *RiskAssessmentHandler) CreateJob(c *gin.Context) {
	var req CreateJobRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	aoiBytes, err := json.Marshal(req.AOI)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid AOI geometry",
			Message: "Failed to serialize AOI geometry",
			Code:    http.StatusBadRequest,
		})
		return
	}

	aoiGeoJSON := string(aoiBytes)
	log.Printf("Received job request with AOI: %s", aoiGeoJSON)

	if err := h.ensureGRPCConnection(); err != nil {
		log.Printf("Failed to connect to orchestrator: %v", err)

		jobID := uuid.New().String()
		log.Printf("Generated mock job ID: %s", jobID)

		response := CreateJobResponse{
			JobID:   jobID,
			Status:  "PENDING",
			Message: "Risk assessment job accepted (MOCK MODE - orchestrator unavailable)",
		}

		c.JSON(http.StatusAccepted, response)
		return
	}

	grpcReq := &generated.CreateJobRequest{
		AoiGeojson: aoiGeoJSON,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	grpcResp, err := h.orchestratorClient.CreateRiskAssessmentJob(ctx, grpcReq)
	if err != nil {
		log.Printf("gRPC call failed: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Message: "Failed to create risk assessment job",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	response := CreateJobResponse{
		JobID:   grpcResp.JobId,
		Status:  grpcResp.Status.String(),
		Message: grpcResp.Message,
	}

	c.JSON(http.StatusAccepted, response)
}

// GET /api/v1/wildfire-risk-jobs/:job_id
func (h *RiskAssessmentHandler) GetJobStatus(c *gin.Context) {
	jobID := c.Param("job_id")

	if jobID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Missing job ID",
			Message: "Job ID is required in the URL path",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Try to connect to orchestrator
	if err := h.ensureGRPCConnection(); err != nil {
		log.Printf("Failed to connect to orchestrator: %v", err)

		// Return mock response when orchestrator is not available
		response := JobStatusResponse{
			JobID:     jobID,
			Status:    "PENDING",
			Message:   "Job is being processed (MOCK MODE - orchestrator unavailable)",
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		c.JSON(http.StatusOK, response)
		return
	}

	// Create gRPC request
	grpcReq := &generated.GetJobStatusRequest{
		JobId: jobID,
	}

	// Call orchestrator service
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcResp, err := h.orchestratorClient.GetJobStatus(ctx, grpcReq)
	if err != nil {
		log.Printf("gRPC call failed: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Message: "Failed to get job status",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Convert gRPC response to REST response
	response := JobStatusResponse{
		JobID:     grpcResp.JobId,
		Status:    grpcResp.Status.String(),
		Message:   grpcResp.Message,
		CreatedAt: grpcResp.CreatedAt,
		UpdatedAt: grpcResp.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// GET /api/v1/wildfire-risk-jobs/:job_id/results
func (h *RiskAssessmentHandler) GetJobResults(c *gin.Context) {
	jobID := c.Param("job_id")

	if jobID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Missing job ID",
			Message: "Job ID is required in the URL path",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// For now, always return mock results
	c.JSON(http.StatusOK, gin.H{
		"job_id":  jobID,
		"status":  "COMPLETE",
		"message": "Mock results (orchestrator not implemented yet)",
		"results": gin.H{
			"total_assets":       42,
			"high_risk_assets":   8,
			"medium_risk_assets": 15,
			"low_risk_assets":    19,
			"download_url":       nil,
		},
	})
}

// Close gRPC connection
func (h *RiskAssessmentHandler) Close() error {
	if h.orchestratorConn != nil {
		return h.orchestratorConn.Close()
	}
	return nil
}
