package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/satriofjar/dashboard-api/handler"
)

func main(){
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})	

	r.GET("/summary", handler.DashboardSummary)

	r.Run(":8000")

}
