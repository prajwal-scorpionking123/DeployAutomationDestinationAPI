package main

import (
	"github.com/team_six/DESTINATION_API/controllers"
	"github.com/team_six/DESTINATION_API/controllers/deploycontroller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("api/backup", controllers.BackupFiles)
	r.POST("api/deployFiles", controllers.DeployeFiles)
	r.POST("/api/deployMultiple", deploycontroller.DeployFiles)

	r.Static("../PRODUCTION", "../PRODUCTION")
	r.Static("../BACKUP", "../BACKUP")
	r.Static("../SOURCE", "../SOURCE")

	r.Run(":3002")
}
