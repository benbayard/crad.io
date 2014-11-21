package main

import (
	"fmt"
	"github.com/benbayard/crad.io/magic"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	router := httprouter.New()
	crads, cmcs := crad.GetCrads()

	cc := &crad.CradController{
		Crads: crads,
		Cmcs:  cmcs,
	}

	router.GET("/crad/:crad", cc.Show)
	router.GET("/cmc/:cmc",   cc.Cmc)

	a := crad.CreateUser("Ben Bayard", "bjbayard@gmail.com", "bjbayard", "testpassword")

	fmt.Printf("%#v", a)

	log.Fatal(http.ListenAndServe(":8080", router))

}
