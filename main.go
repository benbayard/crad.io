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

	uc := &crad.UserController{}
	dc := &crad.DeckController{
		Crads: crads,
	}

	router.GET("/crad/:crad", cc.Show)
	router.GET("/cmc/:cmc", cc.Cmc)

	router.GET( "/user", uc.Index)
	router.POST("/user", uc.Create)

	router.GET( "/user/:id", uc.Show)
	router.PUT( "/user/:id", uc.Update)
	router.POST("/user/:id", uc.Update)

	router.POST("/login", uc.Login)

	router.POST( "/decks/:username", dc.Create)
	router.PUT(  "/decks/:username", dc.AddCrad)
	router.PATCH("/decks/:username", dc.EditCrad)

	// user, token := crad.UserNew("Ben B", "bjbayard3@gmail.com", "bjbayard3", "testpassword")

	// fmt.Printf("%#v \n", user)
	// fmt.Printf("Token: \n")
	// fmt.Printf("%#v \n", token)

	// valid, err := user.ValidToken(token)

	// fmt.Printf("Valid Token For User: %#v, %#v \n", valid, err)
	fmt.Println("Starting Server...")

	log.Fatal(http.ListenAndServe(":8080", router))

}
