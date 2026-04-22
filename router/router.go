package router

import (
	"github.com/Mayank-Bihani/eta-service/handlers"
	"github.com/gin-gonic/gin"
)

// func Setup(etaHandler *handlers.ETAHandler) *gin.Engine {
// 	r := gin.Default()

// 	r.GET("/health", func(c *gin.Context) {
// 		c.JSON(200, gin.H{"status": "ok"})
// 	})

// 	api := r.Group("/api")
// 	{
// 		api.POST("/order/eta", etaHandler.GetETA)
// 	}

// 	return r
// }

// package router

// import (
// 	"github.com/Mayank-Bihani/eta-service/handlers"
// 	"github.com/gin-gonic/gin"
// )

func Setup(etaHandler *handlers.ETAHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		api.POST("/order/eta", etaHandler.GetETA)
		api.GET("/order/:id", etaHandler.GetOrder)
		api.POST("/restaurant/queue", etaHandler.SetQueueDepth)
	}

	return r
}