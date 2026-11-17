package api

//nolint:depguard
import "github.com/gin-gonic/gin"

func NewRouter(h *Handler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/convert", h.Convert)
	}

	return r
}
