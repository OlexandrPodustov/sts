package main

import (
	"github.com/labstack/echo"
	"sts"
)

func main() {
	e := echo.New()

	e.GET("/users/:id", sts.GetUser)
	e.GET("/take", sts.Take)
	e.GET("/fund", sts.Fund)
	e.GET("/announceTournament", sts.AnnounceTournament)
	e.GET("/joinTournament", sts.JoinTournament)
	e.POST("/resultTournament", sts.ResultTournament)
	e.GET("/balance", sts.Balance)
	e.GET("/reset", sts.ResetDB)

	e.Logger.Fatal(e.Start(":8080"))
}
