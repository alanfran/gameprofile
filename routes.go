package main

import (
	"github.com/gin-gonic/gin"
)

func (a *App) initRoutes() {
	r := gin.Default()

	a.engine = r
}
