package main

import (
	"sts"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.GET("/fund", sts.Fund)
	e.GET("/announceTournament", sts.AnnounceTournament)
	e.GET("/joinTournament", sts.JoinTournament)
	e.POST("/resultTournament", sts.ResultTournament)
	e.GET("/balance", sts.Balance)
	e.GET("/reset", sts.ResetDB)

	e.Logger.Fatal(e.Start(":8080"))
}
