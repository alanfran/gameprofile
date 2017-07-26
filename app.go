package main

import (
	"github.com/alanfran/gameprofile/profile"
	"github.com/gin-gonic/gin"
)

type App struct {
	profiles profile.Storer
	Config
	engine *gin.Engine
}

type Config struct {
	dbAddress  string
	dbUser     string
	dbPassword string
	dbDatabase string
}

// NewApp initializes a new App with a profile.Storer, registers application routes, then returns a reference to the App.
func NewApp(store profile.Storer) *App {
	a := &App{profiles: store}

	a.initRoutes()

	return a
}

// Run runs the application on the given interface/port.
// Example: app.Run(":80")
func (a *App) Run(port string) {
	a.engine.Run(port)
}
