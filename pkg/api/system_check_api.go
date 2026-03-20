package api

import (
	"net/http"

	"novel-video-workflow/pkg/providers"

	"github.com/gin-gonic/gin"
)

// SystemCheckAPI handles system health check endpoints.
type SystemCheckAPI struct {
	providers providers.ProviderBundle
}

// NewSystemCheckAPI creates a new system check API.
func NewSystemCheckAPI(bundle providers.ProviderBundle) *SystemCheckAPI {
	return &SystemCheckAPI{providers: bundle}
}

// RegisterRoutes registers system check routes.
func (api *SystemCheckAPI) RegisterRoutes(router *gin.Engine) {
	router.GET("/api/system/check", api.CheckSystem)
}

// CheckSystem runs health checks on all providers.
func (api *SystemCheckAPI) CheckSystem(c *gin.Context) {
	results := []SystemCheckResultResponse{
		convertHealthCheckResult(api.providers.TTS.HealthCheck()),
		convertHealthCheckResult(api.providers.Subtitle.HealthCheck()),
		convertHealthCheckResult(api.providers.Image.HealthCheck()),
		convertHealthCheckResult(api.providers.Project.HealthCheck()),
	}

	canStart := true
	for _, result := range results {
		if result.Severity == string(providers.SeverityBlocking) {
			canStart = false
			break
		}
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Status: "success",
		Data: SystemCheckUpdatedPayload{
			Results:  results,
			CanStart: canStart,
		},
	})
}

func convertHealthCheckResult(result providers.HealthCheckResult) SystemCheckResultResponse {
	return SystemCheckResultResponse{
		Provider: result.Provider,
		Severity: string(result.Severity),
		Message:  result.Message,
	}
}