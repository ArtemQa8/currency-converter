package api

//nolint:depguard
import "github.com/gin-gonic/gin"

// NewRouter создаёт новый роутер.
func NewRouter(h *Handler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/convert", h.Convert)
	}

	return r
}
