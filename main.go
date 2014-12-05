package main

import (
	"fmt"
	"github.com/benbayard/crad.io/magic"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"compress/gzip"
	"io"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}
 
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
 

type gzipHandler struct {
	handler http.Handler
}

func (gzh *gzipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n|      Method      |      Path      |\n| %#v | %#v |\n", r.Method, r.URL.Path)
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		gzh.handler.ServeHTTP(w, r)
		return
	}
	w.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(w)
	defer gz.Close()
	gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
	gzh.handler.ServeHTTP(gzr, r)	
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "assets/html/index.html")
}

func ServeFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	http.ServeFile(w, r, "assets/"+ps.ByName("file"))
}


func main() {
	router := httprouter.New()
	crads, cmcs, cradAry := crad.GetCrads()

	cc := &crad.CradController{
		Crads:   crads,
		Cmcs:    cmcs,
		CradAry: cradAry,
	}

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

	gzh := &gzipHandler{
		handler: router,
	}

	log.Fatal(http.ListenAndServe(":8080", gzh))

}
