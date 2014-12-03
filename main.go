package main

import (
	"fmt"
	"github.com/benbayard/crad.io/magic"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "assets/html/index.html")
}

func ServeFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	http.ServeFile(w, r, "assets/"+ps.ByName("file"))
}

// var GlobalConnection crad.DBConnection

func main() {
	router := httprouter.New()
	crads, cmcs, cradAry := crad.GetCrads()

	cc := &crad.CradController{
		Crads:   crads,
		Cmcs:    cmcs,
		CradAry: cradAry,
	}

	// crad.DatabaseConnect()

	uc := &crad.UserController{}
	dc := &crad.DeckController{
		Crads: crads,
	}

	router.GET("/crad/:crad", cc.Show)
	router.GET("/crads/:search", cc.Search)

	router.GET("/cmc/:cmc", cc.Cmc)

	router.GET("/user", uc.Index)
	router.POST("/user", uc.Create)

	router.GET("/user/:id", uc.Show)
	router.PUT("/user/:id", uc.Update)
	router.POST("/user/:id", uc.Update)

	router.POST("/login", uc.Login)
	router.POST("/admin/:username", uc.ValidToken)

	router.POST("/decks/:username", dc.Create)
	router.PUT("/decks/:username", dc.AddCrad)
	router.PATCH("/decks/:username", dc.EditCrad)
	router.GET("/decks/:username/:deckname", dc.Show)

	router.GET("/app/*path", Index)
	router.GET("/", Index)

	router.GET("/assets/*file", ServeFile)

	fmt.Println("Starting Server...")

	log.Fatal(http.ListenAndServe(":8080", router))

}
