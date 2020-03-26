package sts

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

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
