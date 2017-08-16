package main

import (
	"net/http"

	"github.com/labstack/echo"
)

// e.GET("/users/:id", getUser)
func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

//e.GET("/take", take)
func take(c echo.Context) error {
	// Get team and points from the query string
	playerId := c.QueryParam("playerId")
	points := c.QueryParam("points")
	return c.String(http.StatusOK, "playerId:"+playerId+", will be charged from his balance, points:"+points)
}

//e.GET("/fund", fund)
func fund(c echo.Context) error {
	// Get team and points from the query string
	playerId := c.QueryParam("playerId")
	points := c.QueryParam("points")
	return c.String(http.StatusOK, "tournamentId:"+playerId+", will receive to his balance, points:"+points)
}

func announceTournament(c echo.Context) error {
	// Get tournamentId and deposit from the query string
	tournamentId := c.QueryParam("tournamentId")
	deposit := c.QueryParam("deposit")
	return c.String(http.StatusOK, "tournamentId:"+tournamentId+", deposit:"+deposit)
}

func joinTournament(c echo.Context) error {
	// Get tournamentId and players from the query string
	tournamentId := c.QueryParam("tournamentId")
	playerId := c.QueryParam("playerId")

	allQueryValues := c.QueryParams()

	var sliceOfBakers []string
	var allBakers string
	for key, val := range allQueryValues {
		if key == "backerId" {
			sliceOfBakers = append(sliceOfBakers, val...)
		}
	}

	for _, v := range sliceOfBakers {
		allBakers += v + ", "
	}
	return c.String(http.StatusOK, "tournamentId:"+tournamentId+", playerId:"+playerId+", allQueryValues:"+allBakers)
}

// e.POST("/resultTournament", resultTournament)
func resultTournament(c echo.Context) error {
	// Get tournamentId and players from the query string
	name := c.FormValue("name")
	email := c.FormValue("email")
	return c.String(http.StatusOK, "name:"+name+", email:"+email)
}

func balance(c echo.Context) error {
	// Get tournamentId and players from the query string
	playerId := c.QueryParam("playerId")
	return c.String(http.StatusOK, "balance"+playerId)
}

func reset(c echo.Context) error {
	// Get tournamentId and players from the query string

	return c.String(http.StatusOK, "DB was cleared")
}

func main() {
	e := echo.New()

	e.GET("/users/:id", getUser)
	e.GET("/take", take)
	e.GET("/fund", fund)
	e.GET("/announceTournament", announceTournament)
	e.GET("/joinTournament", joinTournament)
	e.POST("/resultTournament", resultTournament)
	e.GET("/balance", balance)
	e.GET("/reset", reset)

	e.Logger.Fatal(e.Start(":1323"))
}
