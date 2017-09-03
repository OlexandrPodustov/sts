package sts

import (
	"fmt"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
)

var (
	mongo_address = os.Getenv("MONGO_ADDRESS")
	MgoAddr       = mongo_address //+ ":27017"
)

type Player struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	Name   string
	Points int
}

func init() {
	session, err := mgo.Dial("172.17.0.2")
	if err != nil {
		log.Println("2")
		panic(err)
		log.Println("3", MgoAddr)
	}
	defer session.Close()

	c := session.DB("stsDB").C("PlayersCollection")
	err = c.Insert(&Player{Name: "Ale", Points: 11112},
		&Player{Name: "Cla", Points: 222})
	if err != nil {
		log.Fatal(err)
	}

	result := Player{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ale's Points:", result.Points)
}

// e.GET("/users/:id", getUser)
func GetUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

//e.GET("/take", take)
func Take(c echo.Context) error {
	// Get team and points from the query string
	playerId := c.QueryParam("playerId")
	points := c.QueryParam("points")
	return c.String(http.StatusOK, "playerId:"+playerId+", will be charged from his balance, points:"+points)
}

//e.GET("/fund", fund)
func Fund(c echo.Context) error {
	// Get team and points from the query string
	playerId := c.QueryParam("playerId")
	points := c.QueryParam("points")

	log.Println("Fund", playerId, points)

	return c.String(http.StatusOK, "playerId:"+playerId+", will receive to his balance, points:"+points)
}

func AnnounceTournament(c echo.Context) error {
	// Get tournamentId and deposit from the query string
	tournamentId := c.QueryParam("tournamentId")
	deposit := c.QueryParam("deposit")
	return c.String(http.StatusOK, "tournamentId:"+tournamentId+", deposit:"+deposit)
}

func JoinTournament(c echo.Context) error {
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
func ResultTournament(c echo.Context) error {
	// Get tournamentId and players from the query string
	name := c.FormValue("name")
	email := c.FormValue("email")
	return c.String(http.StatusOK, "name:"+name+", email:"+email)
}

func Balance(c echo.Context) error {
	// Get tournamentId and players from the query string
	playerId := c.QueryParam("playerId")
	return c.String(http.StatusOK, "balance"+playerId)
}

func ResetDB(c echo.Context) error {
	m := "DB was cleared. "
	m2 := "mongo_address:" + MgoAddr
	completeMessage := m + m2
	log.Println(completeMessage)

	return c.String(http.StatusOK, completeMessage)
}
