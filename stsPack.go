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
	dbName                 = "stsDB"
	playersColl            = "PlayersCollection"
	tournamentsColl        = "TournamentsCollection"
	tournamentsPlayersColl = "TournamentsPlayersCollection"
)

var mongo_address = os.Getenv("MONGO_ADDRESS")

type player struct {
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
type winnPlayer struct {
	PlayerId string
	Prize    int
}

type tournament struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	TournamentId int
	Deposit      int
	Players      []resultPlayer
	Winners      []winnPlayer
	Timestamp    time.Time
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

	result := player{}
	err = collection.Find(bson.M{"playerid": playerId}).Sort("-timestamp").One(&result)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, fmt.Sprint(err))
	}
	var response string
	currentPointsString := strconv.Itoa(result.Balance)
	if calculatedPoints := result.Balance - pointsToCharge; calculatedPoints >= 0 {
		err = collection.Insert(&player{
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

	result := player{}
	err = collection.Find(bson.M{"playerid": playerId}).Sort("-timestamp").One(&result)
	if err == nil {
		log.Println("player found, updating...")
		err = collection.Insert(&player{
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
		err = collection.Insert(&player{
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
	tournam := c.QueryParam("tournamentId")
	tournamentId, err := strconv.Atoi(tournam)
	if err != nil {
		panic(err)
	}

	depos := c.QueryParam("deposit")
	deposit, err := strconv.Atoi(depos)
	if err != nil {
		panic(err)
	}

	session, err := mgo.Dial(mongo_address)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	collection := session.DB(dbName).C(tournamentsColl)

	result := tournament{}
	err = collection.Find(bson.M{"tournamentid": tournamentId}).Sort("-timestamp").One(&result)
	if err == nil {
		log.Println("Tournament found, doing nothing. Can't insert tournament with the same id")
		/*err = collection.Insert(&tournament{TournamentId: tournamentId, Deposit: deposit})

		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, fmt.Sprint(err))
		}*/
	} else if p := err.Error(); p == "not found" {
		log.Println("Tournament wasn't found, creating one")
		err = collection.Insert(&tournament{TournamentId: tournamentId, Deposit: deposit})

		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, fmt.Sprint(err))
		}
	} else {
		log.Println(err)
		return c.String(http.StatusInternalServerError, fmt.Sprint(err))
	}

	return c.String(http.StatusOK, "tournamentId:"+tournam+", deposit:"+depos)
}

//not completed
func JoinTournament(c echo.Context) error {
	// Get tournamentId and players from the query string
	tournam := c.QueryParam("tournamentId")
	player := c.QueryParam("playerId")
	tournamentId, err := strconv.Atoi(tournam)
	if err != nil {
		panic(err)
	}
	//playerId, err := strconv.Atoi(player)
	//if err != nil {
	//	panic(err)
	//}
	allQueryValues := c.QueryParams()

	var sliceOfBakers []string
	var allBakers string
	for key, val := range allQueryValues {
		if key == "backerId" {
			sliceOfBakers = append(sliceOfBakers, val...)
		}
	}

	session, err := mgo.Dial(mongo_address)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	collection := session.DB(dbName).C(tournamentsColl)

	result := tournament{}
	err = collection.Find(bson.M{"tournamentid": tournamentId}).Sort("-timestamp").One(&result)
	if err == nil {
		log.Println(result)
		newPlayer := resultPlayer{PlayerId: player}
		var players []resultPlayer
		players = append(result.Players, newPlayer)
		err = collection.Insert(&tournament{TournamentId: result.TournamentId, Players: players})
		log.Println("Tournament found, joining")
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, fmt.Sprint(err))
		}
	} else if p := err.Error(); p == "not found" {
		log.Println("Tournament not found, doing nothing. Can't join unexisting tournament")
	} else {
		log.Println(err)
		return c.String(http.StatusInternalServerError, fmt.Sprint(err))
	}

	for _, v := range sliceOfBakers {
		allBakers += v + ", "
	}
	return c.String(http.StatusOK, "tournamentId:"+tournam+", playerId:"+player+", allQueryValues:"+allBakers)
}

// e.POST("/resultTournament", resultTournament)
//not completed
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

	result := player{}
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
