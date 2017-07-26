package main

import "github.com/alanfran/gameprofile/profile"

func main() {
	boltStore, err := profile.NewBoltStore("bolt.db")
	if err != nil {
		panic(err)
	}

	a := NewApp(boltStore)

	a.Run(":80")
}
