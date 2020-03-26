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

var mongoDB *mgo.Session

func Init() {
	mongoAddress := os.Getenv("MONGO_ADDRESS")
	if mongoAddress == "" {
		log.Fatal("empty mongo db address")
	}

	session, err := mgo.Dial(mongoAddress)
	if err != nil {
		log.Fatal(err)
	}

	err = session.Ping()
	if err != nil {
		log.Fatal(err)
	}

	mongoDB = session
}

func Take(c echo.Context) error {
	playerID := c.QueryParam("playerId")
	points := c.QueryParam("points")
	pointsToCharge, err := strconv.Atoi(points)
	if err != nil {
		panic(err)
	}

	collection := mongoDB.DB(dbName).C(playersColl)

	result := player{}
	err = collection.Find(bson.M{"playerid": playerID}).Sort("-timestamp").One(&result)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, fmt.Sprint(err))
	}
	var response string
	currentPointsString := strconv.Itoa(result.Balance)
	if calculatedPoints := result.Balance - pointsToCharge; calculatedPoints >= 0 {
		err = collection.Insert(&player{
			PlayerId:  playerID,
			PointsOld: result.Balance,
			Balance:   calculatedPoints,
			Timestamp: time.Now()})

		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, fmt.Sprint(err))
		}
		response = "player " + playerID + " points from db " + currentPointsString + " points wanted to charge " + points
	} else {
		response = "Take can't be processed. Insufficient amount of points: player - " + playerID + " points wanted to charge " + points
	}

	return c.String(http.StatusOK, response)
}

func Fund(c echo.Context) error {
	playerID := c.QueryParam("playerId")
	points := c.QueryParam("points")
	pointsToFund, err := strconv.Atoi(points)
	if err != nil {
		panic(err)
	}

	collection := mongoDB.DB(dbName).C(playersColl)

	result := player{}
	err = collection.Find(bson.M{"playerid": playerID}).Sort("-timestamp").One(&result)
	if err == nil {
		log.Println("player found, updating...")
		err = collection.Insert(&player{
			PlayerId:  playerID,
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
			PlayerId:  playerID,
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

	return c.String(http.StatusOK, "done")
}

func AnnounceTournament(c echo.Context) error {
	tournam := c.QueryParam("tournamentId")
	tournamentID, err := strconv.Atoi(tournam)
	if err != nil {
		panic(err)
	}

	depos := c.QueryParam("deposit")
	deposit, err := strconv.Atoi(depos)
	if err != nil {
		panic(err)
	}

	collection := mongoDB.DB(dbName).C(tournamentsColl)

	result := tournament{}
	err = collection.Find(bson.M{"tournamentid": tournamentID}).Sort("-timestamp").One(&result)
	if err == nil {
		log.Println("Tournament found, doing nothing. Can't insert tournament with the same id")
		/*err = collection.Insert(&tournament{TournamentId: tournamentId, Deposit: deposit})

		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, fmt.Sprint(err))
		}*/
	} else if p := err.Error(); p == "not found" {
		log.Println("Tournament wasn't found, creating one")
		err = collection.Insert(&tournament{TournamentId: tournamentID, Deposit: deposit})

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

func JoinTournament(c echo.Context) error {
	playerID := c.QueryParam("playerId")

	tournam := c.QueryParam("tournamentId")
	tournamentID, err := strconv.Atoi(tournam)
	if err != nil {
		panic(err)
	}

	var sliceOfBakers []backer
	for key, val := range c.QueryParams() {
		if key == "backerId" {
			for _, v := range val {
				sliceOfBakers = append(sliceOfBakers, backer{BackerId: v})
			}
		}
	}

	collection := mongoDB.DB(dbName).C(tournamentsColl)

	result := tournament{}
	err = collection.Find(bson.M{"tournamentid": tournamentID}).Sort("-timestamp").One(&result)
	if err == nil {
		newPlayer := playerWithBackers{PlayerId: playerID, Backers: sliceOfBakers}
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
		//smth like
		//1. create and fill a slice of player and backers with their IDs
		var sliceOfBakersAndPlayers []string
		for _, id := range sliceOfBakers {
			sliceOfBakersAndPlayers = append(sliceOfBakersAndPlayers, id.BackerId)
		}
		sliceOfBakersAndPlayers = append(sliceOfBakersAndPlayers, playerID)
		//2. range on that slice with spinning goroutines
		for _, id := range sliceOfBakersAndPlayers {
			go func(id string) {
				//go to db - and check if each one have sufficient amount of points with lock by player
				//send info to channel
				//process all responces. if all of them are true - redeem points
				//release locks by player
			}(id)
		}
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

	return c.String(http.StatusOK, "tournamentId:"+tournam+", playerId:"+playerID)
}

func ResultTournament(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")

	return c.String(http.StatusOK, "name:"+name+", email:"+email)
}

func Balance(c echo.Context) error {
	playerID := c.QueryParam("playerId")

	collection := mongoDB.DB(dbName).C(playersColl)

	result := player{}
	err := collection.Find(bson.M{"playerid": playerID}).Sort("-timestamp").One(&result)
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
	err := mongoDB.DB(dbName).DropDatabase()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprint(err))
	}

	m := "DB was cleared"
	log.Println(m)

	return c.String(http.StatusOK, m)
}
