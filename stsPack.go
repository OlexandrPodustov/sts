package sts

import (
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
	Name      string
	Points    int
	Timestamp time.Time
}

//e.GET("/take?playerId=P1&points=300", fund)
func Take(c echo.Context) error {
	// Get playerID and points from the query string
	playerId := c.QueryParam("playerId")
	points := c.QueryParam("points")
	pointsConverted, err := strconv.Atoi(points)
	if err != nil {
		panic(err)
	}
	//pp := (-1) * pointsConverted

	session, err := mgo.Dial(mongo_address)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//todo: check whether user exists or not, before inserting a row into collection
	//or we just can implement a different logic of retreiving balance
	collection := session.DB(dbName).C(playersColl)
	//check if player has sufficient amount of points
	result := Player{}
	err = collection.Find(bson.M{"name": playerId}).One(&result)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Unable to retrieve data from database, maybe there is no such player")
		log.Fatal(err)
	}
	var res string
	des := strconv.Itoa(result.Points)
	if result.Points-pointsConverted >= 0 {
		err = collection.Insert(&Player{Name: playerId, Points: -pointsConverted, Timestamp: time.Now()})
		if err != nil {
			return c.String(http.StatusInternalServerError, "Unable to insert data into database")
			log.Fatal(err)
		}
		res = "player " + playerId + " points from db " + des + " points wanted to charge " + points
		log.Println("player ", playerId, " points from db ", result.Points, " points wanted to charge ", points)
	} else {
		res = "Take can't be processed, cause of insufficient amount of points " + playerId + " current amount " + des
		log.Println("Take can't be processed, cause of insufficient amount of points ", playerId, "current amount", result.Points)
	}

	return c.String(http.StatusOK, res)
}

//e.GET("/fund?playerId=P1&points=300", fund)
func Fund(c echo.Context) error {
	// Get playerID and points from the query string
	playerId := c.QueryParam("playerId")
	points := c.QueryParam("points")
	pointsConverted, err := strconv.Atoi(points)
	if err != nil {
		panic(err)
	}

	session, err := mgo.Dial(mongo_address)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//todo: check whether user exists or not, before inserting a row into collection
	//or we just can implement a different logic of retreiving balance
	collection := session.DB(dbName).C(playersColl)
	err = collection.Insert(&Player{Name: playerId, Points: pointsConverted, Timestamp: time.Now()})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Unable to insert data into database")
		log.Fatal(err)
	}
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

	session, err := mgo.Dial(mongo_address)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	collection := session.DB(dbName).C(playersColl)

	result := Player{}
	//current implementation as the result of points gives the amount from the first row.
	//todo: change One with All, then range the map, calculate overall points, give the proper response
	err = collection.Find(bson.M{"name": playerId}).One(&result)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Unable to retrieve data from database, maybe there is no such player")
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, result)
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
		return c.String(http.StatusInternalServerError, "Unable to drop database")
		log.Fatal(err)
	}
	return c.String(http.StatusOK, m)
}
