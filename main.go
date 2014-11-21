package main

import (
	"fmt"
	"github.com/benbayard/crad.io/magic"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)

	crads, cmcs := crad.GetCrads()

	cc := &crad.CradController{
		Crads: crads,
		Cmcs:  cmcs,
	}

	router.GET("/crads/cmcs/:cmc", cc.Cmc)
	// router.GET("/crads/:crad", cc.get)

	// fmt.Printf("%#v", cmcs[1])

	log.Fatal(http.ListenAndServe(":8080", router))

}
