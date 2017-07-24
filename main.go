package main

import (
	"./profile"
)

// port 80

/* routes:

/:steamid

/:steamid/coins

/:steamid/punishments



*/

type App struct {
	profiles profile.Storer
}
