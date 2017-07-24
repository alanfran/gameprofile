package profile

import (
	"strconv"
	"testing"
	"time"

	pg "gopkg.in/pg.v4"
)

var nt *PostgresStore
var db *pg.DB

var (
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbDatabase = "test"
)

func init() {
	// load test configuration
	db = pg.Connect(&pg.Options{
		User:     dbUser,
		Password: dbPassword,
		Database: dbDatabase,
	})
	// verify connection
	_, err := db.Exec(`SELECT 1`)
	if err != nil {
		panic("Error connecting to the database.")
	}

	db.Exec(`
		DROP TABLE profiles;
		DROP TABLE punishments;
	`)

	nt = NewPostgresStore(db)
}

func TestProfilePg(t *testing.T) {
	id := "example steamid"
	coins := int64(64)
	inventory := map[string]string{"key": "value"}
	equip := map[string]string{"head": "cowboy hat"}

	err := nt.PutProfile(Profile{ID: id, Coins: coins, Inventory: inventory, Equipment: equip})
	if err != nil {
		t.Error(err)
	}
	p, err := nt.GetProfile(id)
	if err != nil {
		t.Error(err)
	}

	if p.Coins != coins || p.Inventory["key"] != inventory["key"] || p.Equipment["head"] != equip["head"] {
		t.Error("coins: " + strconv.Itoa(int(p.Coins)) + " inv: " + p.Inventory["key"])
	}
}

func TestPunishmentPg(t *testing.T) {
	steamid := "example id"
	by := "someAdmin"
	typ := "rubber hammer"
	reason := "being bad"
	date := time.Now()
	expires := time.Now().AddDate(0, 0, 1)

	p := Punishment{
		PlayerID: steamid,
		By:       by,
		Type:     typ,
		Reason:   reason,
		Date:     date,
		Expires:  expires,
	}
	err := nt.PutPunishment(p)
	if err != nil {
		t.Error(err)
	}

	punishments, err := nt.GetPunishments(steamid)
	if err != nil {
		t.Error(err)
	}

	// if punishment[typ] does not exist, error
	p1, ok := punishments[typ]
	if !ok {
		t.Error("Punishment not in GetPunishments return.")
	}

	// if data is not consistent, error
	if p1.PlayerID != steamid || p1.By != by || p1.Type != typ || p1.Reason != reason {
		t.Error("Bad data returned.")
	}

	err = nt.DelPunishment(p1.ID)
	if err != nil {
		t.Error("Failed to delete punishment from db.", err)
	}
}
