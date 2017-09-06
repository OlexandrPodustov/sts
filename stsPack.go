package sts

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	dbName          = "stsDB"
	playersColl     = "PlayersCollection"
	tournamentsColl = "TournamentsCollection"
)

var mongo_address = os.Getenv("MONGO_ADDRESS")

type Player struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	PlayerId  string
	PointsOld int
	Balance   int
	Timestamp time.Time
}

type resultPlayer struct {
	PlayerId string
	Balance  int
}

//e.GET("/take?playerId=P1&points=300", fund)
func Take(c echo.Context) error {
	// Get playerID and points from the query string
	playerId := c.QueryParam("playerId")
	points := c.QueryParam("points")
	pointsToCharge, err := strconv.Atoi(points)
	if err != nil {
		panic(err)
	}

	session, err := mgo.Dial(mongo_address)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	collection := session.DB(dbName).C(playersColl)

	result := Player{}
	err = collection.Find(bson.M{"playerid": playerId}).Sort("-timestamp").One(&result)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, fmt.Sprint(err))
	}
	var response string
	currentPointsString := strconv.Itoa(result.Balance)
	if calculatedPoints := result.Balance - pointsToCharge; calculatedPoints >= 0 {
		err = collection.Insert(&Player{
			PlayerId:  playerId,
			PointsOld: result.Balance,
			Balance:   calculatedPoints,
			Timestamp: time.Now()})

		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, fmt.Sprint(err))
		}
		response = "player " + playerId + " points from db " + currentPointsString + " points wanted to charge " + points
	} else {
		response = "Take can't be processed. Insufficient amount of points: player - " + playerId + " points wanted to charge " + points
	}

	return c.String(http.StatusOK, response)
}

//e.GET("/fund?playerId=P1&points=300", fund)
func Fund(c echo.Context) error {
	var response string
	// Get playerID and points from the query string
	playerId := c.QueryParam("playerId")
	points := c.QueryParam("points")
	pointsToFund, err := strconv.Atoi(points)
	if err != nil {
		panic(err)
	}

	session, err := mgo.Dial(mongo_address)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	collection := session.DB(dbName).C(playersColl)

	result := Player{}
	err = collection.Find(bson.M{"playerid": playerId}).Sort("-timestamp").One(&result)
	if err == nil {
		log.Println("player found, updating...")
		err = collection.Insert(&Player{
			PlayerId:  playerId,
			PointsOld: result.Balance,
			Balance:   result.Balance + pointsToFund,
			Timestamp: time.Now()})

		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, fmt.Sprint(err))
		}
	} else if p := err.Error(); p == "not found" {
		log.Println("player wasn't found, first insertion")
		err = collection.Insert(&Player{
			PlayerId:  playerId,
			PointsOld: 0,
			Balance:   pointsToFund,
			Timestamp: time.Now()})

		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, fmt.Sprint(err))
		}
	} else {
		log.Println(err)
		return c.String(http.StatusInternalServerError, fmt.Sprint(err))
	}

	return c.String(http.StatusOK, response)
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

	session, err := mgo.Dial(mongo_address)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	collection := session.DB(dbName).C(playersColl)

	result := Player{}
	err = collection.Find(bson.M{"playerid": playerId}).Sort("-timestamp").One(&result)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, fmt.Sprint(err))
	}

	narrowedResult := resultPlayer{
		result.PlayerId,
		result.Balance,
	}

	return c.JSON(http.StatusOK, narrowedResult)
}

func ResetDB(c echo.Context) error {
	m := "DB was cleared. "
	log.Println(m)
	session, err := mgo.Dial(mongo_address)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(dbName).DropDatabase()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprint(err))
		log.Println(err)
	}
	return c.String(http.StatusOK, m)
}
