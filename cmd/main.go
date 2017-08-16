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

//e.GET("/show", show)
func show(c echo.Context) error {
	// Get team and points from the query string
	playerId := c.QueryParam("playerId")
	points := c.QueryParam("points")
	return c.String(http.StatusOK, "playerId:"+playerId+", points:"+points)
}

func main() {
	e := echo.New()
	_ = func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}
	e.GET("/users/:id", getUser)
	e.GET("/take", show)
	//?playerId=P1&points=300
	e.Logger.Fatal(e.Start(":1323"))
}
