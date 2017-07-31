package main

import (
	"net/http"

	"github.com/alanfran/gameprofile/profile"
	"github.com/gin-gonic/gin"
)

func (a *App) GetPunishments(c *gin.Context) {
	steamid := c.Param("steamid")

	ps, err := a.profiles.GetPunishments(steamid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No punishments found for that player.",
		})
		return
	}

	c.JSON(http.StatusOK, ps)
}

func (a *App) PostPunishments(c *gin.Context) {
	steamid := c.Param("steamid")

	var p profile.Punishment

	err := c.Bind(&p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request. Make sure your JSON is correct.",
		})
		return
	}

	if steamid != p.PlayerID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "PlayerID in body does not match the ID in the URL.",
		})
		return
	}

	err = a.profiles.PutPunishment(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (a *App) PutPunishments(c *gin.Context) {
	steamid := c.Param("steamid")

	var ps map[string]profile.Punishment
	err := c.Bind(&ps)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error processing your request. Please make sure your JSON is well-formatted.",
		})
		return
	}

	for _, v := range ps {
		if v.PlayerID != steamid {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "PlayerID does not match ID in the URL.",
			})
			return
		}
		err = a.profiles.PutPunishment(v)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error storing punishments. Please try again later.",
			})
			return
		}
	}

	c.String(http.StatusNoContent, "")
}
