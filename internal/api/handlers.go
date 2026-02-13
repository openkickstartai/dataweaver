package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/dataweaver/internal/schema"
	"github.com/dataweaver/internal/workflow"
)

type Handler struct {
	db       *sql.DB
	detector *schema.Detector
	workflow *workflow.Engine
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db:       db,
		detector: schema.NewDetector(),
		workflow: workflow.NewEngine(),
	}
}

func (h *Handler) Health(c *gin.Context) {
	status := "healthy"
	httpStatus := http.StatusOK
	dbStatus := "connected"

	if err := h.db.Ping(); err != nil {
		status = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
		dbStatus = "disconnected"
	}

	c.JSON(httpStatus, gin.H{
		"status":  status,
		"service": "dataweaver",
		"db":      dbStatus,
	})
}

type DetectSchemaRequest struct {
	Data   string `json:"data" binding:"required"`
	Format string `json:"format" binding:"required"`
}

func (h *Handler) DetectSchema(c *gin.Context) {
	var req DetectSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Format != "json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only JSON format supported"})
		return
	}

	schema, err := h.detector.DetectFromJSON([]byte(req.Data))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schema)
}

type CreateWorkflowRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Description string                 `json:"description"`
	Steps       []workflow.Step        `json:"steps" binding:"required"`
	Config      map[string]interface{} `json:"config"`
}

func (h *Handler) CreateWorkflow(c *gin.Context) {
	var req CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wf := &workflow.Workflow{
		Name:        req.Name,
		Description: req.Description,
		Steps:       req.Steps,
		Config:      req.Config,
	}

	id, err := h.workflow.Create(wf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "workflow": wf})
}

func (h *Handler) GetWorkflow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workflow ID"})
		return
	}

	wf, err := h.workflow.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workflow not found"})
		return
	}

	c.JSON(http.StatusOK, wf)
}

type TransformRequest struct {
	WorkflowID int                    `json:"workflow_id" binding:"required"`
	Data       map[string]interface{} `json:"data" binding:"required"`
}

func (h *Handler) TransformData(c *gin.Context) {
	var req TransformRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.workflow.Execute(req.WorkflowID, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}