package sources

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) HandleAdd(c *gin.Context) {
	var req SourcesAdd
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// request ke BE AI
	var parseReq struct {
		Type string `json:"type"`
		Url  string `json:"url"`
	}
	parseReq.Type = req.Type
	parseReq.Url = req.Url

	jsonReq, err := json.Marshal(parseReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	resp, err := http.Post(os.Getenv("BE_AI_URL")+"/probe", "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Probe service unreachable"})
		return
	}
	defer resp.Body.Close()

	var probeResult struct {
		Exists bool   `json:"exists"`
		Detail string `json:"detail"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&probeResult); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse probe response"})
		return
	}

	if !probeResult.Exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "Source verification failed",
			"detail": probeResult.Detail,
		})
		return
	}

	sources, err := h.svc.Add(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sources)
}

func (h *Handler) HandleRequest(c *gin.Context) {
	sources, err := h.svc.Request(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, SourcesResponse{Sources: sources})
}

func (h *Handler) HandleRequestSourceID(c *gin.Context) {
	sources, err := h.svc.RequestSourcesId(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, SourcesResponse{Sources: sources})
}

func (h *Handler) HandleRequestById(c *gin.Context) {
	sourceId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": c.Param(":id")})
		return
	}

	source, err := h.svc.RequestByID(c.Request.Context(), sourceId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, source)
}

func (h *Handler) HandleUpdateById(c *gin.Context) {
	sourceId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req SourcesUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	source, err := h.svc.UpdateById(c.Request.Context(), req, sourceId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, source)
}

func (h *Handler) HandleDeleteById(c *gin.Context) {
	sourceId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	source, err := h.svc.DeleteById(c.Request.Context(), sourceId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, source)
}
