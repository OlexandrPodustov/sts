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

type playerWithBackers struct {
	PlayerId string
	Backers  []backer
}

type backer struct {
	BackerId     string
	Balance      int
	InterestAmt  int
	InterestRate int
}

type winnPlayer struct {
	PlayerId string
	Prize    int
}

type tournament struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	TournamentId int
	Deposit      int
	Player       playerWithBackers
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
	playerid := c.QueryParam("playerId")
	tournamentId, err := strconv.Atoi(tournam)
	if err != nil {
		panic(err)
	}
	//allQueryValues := c.QueryParams()

	var sliceOfBakers []backer
	for key, val := range c.QueryParams() {
		if key == "backerId" {
			for _, v := range val {
				sliceOfBakers = append(sliceOfBakers, backer{BackerId: v})
			}
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
		newPlayer := playerWithBackers{PlayerId: playerid, Backers: sliceOfBakers}
		//commented due to incomplete status. Start.
		//collec2 := session.DB(dbName).C(playersColl)
		//if len(sliceOfBakers) > 0 {
		//	//retrieving the deposit amount of the tournament
		//	amt := result.Deposit / (len(sliceOfBakers) + 1)
		//
		//	//check the balance should be done during charging the player.
		//	//Between two actions (check here and redeem there) the player or redeemer can perform other action?
		//	//If yes implementation should be changed.
		//	log.Println(result.Deposit, len(sliceOfBakers)+1, amt)
		//	playerBalance := player{}
		//	err = collec2.Find(bson.M{"playerid": playerid}).Sort("-timestamp").One(&playerBalance)
		//	//checking the balance of the playerid, and of backers
		// todo: probably spin goroutines with request to db, charge only if all backers and player will have sufficient amount of points
		//	if err != nil {
		//		log.Println(err)
		//		return c.String(http.StatusInternalServerError, fmt.Sprint(err))
		//	}
		//	log.Println(playerBalance.Balance >= amt)
		//} else {
		//	log.Println("there are no backers, checking the balance of the playerid")
		//	playerBalance := player{}
		//	err = collec2.Find(bson.M{"playerid": playerid}).Sort("-timestamp").One(&playerBalance)
		//	log.Println(result.Deposit, playerBalance.Balance, playerBalance.Balance >= result.Deposit)
		//}
		//todo: add check whether such a playerid is already in this tournament or not. If not - then update.
		err = collection.Insert(
			&tournament{
				TournamentId: result.TournamentId,
				Deposit:      result.Deposit,
				Player:       newPlayer})

		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, fmt.Sprint(err))
		}
		log.Println("Tournament found, joining")
	} else if p := err.Error(); p == "not found" {
		log.Println("Tournament not found, doing nothing. Can't join unexisting tournament")
	} else {
		log.Println(err)
		return c.String(http.StatusInternalServerError, fmt.Sprint(err))
	}

	return c.String(http.StatusOK, "tournamentId:"+tournam+", playerId:"+playerid)
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
