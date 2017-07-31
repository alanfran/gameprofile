package main

import (
	"net/http"

	"github.com/alanfran/gameprofile/profile"
	"github.com/gin-gonic/gin"
)

func (a *App) GetProfile(c *gin.Context) {
	steamid := c.Param("steamid")

	if steamid == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Please supply a SteamID",
		})
		return
	}

	p, err := a.profiles.GetProfile(steamid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	profile := NewProfileWithHash(p)

	c.JSON(http.StatusOK, profile)
}

func (a *App) PostProfile(c *gin.Context) {
	var p profile.Profile
	err := c.Bind(&p)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if p.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please supply a profile with an ID.",
		})
		return
	}

	p2, err := a.profiles.GetProfile(p.ID)
	if err == nil || p2.ID != "" {
		c.JSON(http.StatusConflict, NewProfileWithHash(p2))
		return
	}

	err = a.profiles.PutProfile(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "An error occurred while storing your profile. Please try again later.",
		})
		return
	}

	c.JSON(http.StatusCreated, NewProfileWithHash(p))
}

func (a *App) PutProfile(c *gin.Context) {
	steamid := c.Param("steamid")

	if steamid == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Please supply a SteamID in the URL.",
		})
		return
	}

	var pwh ProfileWithHash
	err := c.Bind(&pwh)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "There was an error processing your request. Please make sure your JSON is well-formed.",
		})
		return
	}

	if steamid != pwh.ID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The ID in the request body does not match the one in the URL.",
		})
		return
	}

	p, err := a.profiles.GetProfile(steamid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Could not find a profile with that SteamID.",
		})
		return
	}

	currentPwh := NewProfileWithHash(p)

	if pwh.Hash != currentPwh.Hash {
		c.JSON(http.StatusConflict, currentPwh)
		return
	}

	err = a.profiles.PutProfile(pwh.Profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "An error occurred while storing the profile. Please try again later.",
		})
		return
	}

	c.JSON(http.StatusOK, NewProfileWithHash(pwh.Profile))
}
