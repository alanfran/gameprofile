package main

import (
	"github.com/gin-gonic/gin"
)

func (a *App) initRoutes() {
	r := gin.Default()
	a.engine = r

	r.GET("/:steamid", a.GetProfile)
	r.POST("/", a.PostProfile)
	r.PUT("/:steamid", a.PutProfile)

	r.GET("/:steamid/punishments", a.GetPunishments)
	r.POST("/:steamid/punishments", a.PostPunishments)
	r.PUT("/:steamid/punishments", a.PutPunishments)
}
