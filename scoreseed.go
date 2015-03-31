package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	r "gopkg.in/dancannon/gorethink.v0"
)

//ScoreEntry for scores
type ScoreEntry struct {
	ID         string `gorethink:"id,omitempty"`
	PlayerName string
	Score      int
}

func main() {
	fmt.Println("Connecting to RethinkDB")

	session, err := r.Connect(r.ConnectOpts{
		Address:  "127.0.0.2:28015",
		Database: "players",
	})

	if err != nil {
		log.Fatal("Could not connect")
	}

	err = r.Db("players").TableDrop("scores").Exec(session)
	err = r.Db("players").TableCreate("scores").Exec(session)
	if err != nil {
		log.Fatal("Could not create table")
	}

	err = r.Db("players").Table("scores").IndexCreate("Score").Exec(session)
	if err != nil {
		log.Fatal("Could not create index")
	}

	for i := 0; i < 1000; i++ {
		player := new(ScoreEntry)
		player.ID = strconv.Itoa(i)
		player.PlayerName = fmt.Sprintf("Player %d", i)
		player.Score = rand.Intn(100)
		_, err := r.Table("scores").Insert(player).RunWrite(session)
		if err != nil {
			log.Fatal(err)
		}
	}

	for {
		var scoreentry ScoreEntry
		pl := rand.Intn(1000)
		sc := rand.Intn(6) - 2
		res, err := r.Table("scores").Get(strconv.Itoa(pl)).Run(session)
		if err != nil {
			log.Fatal(err)
		}

		err = res.One(&scoreentry)
		scoreentry.Score = scoreentry.Score + sc
		_, err = r.Table("scores").Update(scoreentry).RunWrite(session)
		time.Sleep(100 * time.Millisecond)
	}
}
