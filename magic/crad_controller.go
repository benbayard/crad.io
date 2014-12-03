package crad

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sort"
	"strconv"
	"strings"
	// "fmt"
)

type CradController struct {
	Crads   map[string]Crad
	Cmcs    map[float64][]*Crad
	CradAry []string
}

func (cc *CradController) Cmc(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	cmcString := ps.ByName("cmc")

	cmc, err := strconv.ParseFloat(cmcString, 64)

	crads, ok := cc.Cmcs[cmc]

	if !ok {
		http.Error(w, "Does not exist", http.StatusNotFound)
		return
	}

	js, err := json.Marshal(crads)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func (cc *CradController) Show(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	name := ps.ByName("crad")

	crad, ok := cc.Crads[name]

	if !ok {
		http.Error(w, "Does not exist", http.StatusNotFound)
		return
	}

	js, err := json.Marshal(crad)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func (cc *CradController) Search(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	term := ps.ByName("search")
	termLen := len(term)
	cradAryLen := len(cc.CradAry)

	index := sort.Search(cradAryLen, func(i int) bool {
		oldTermLen := termLen
		cradName := cc.CradAry[i]
		if termLen > len(cradName) {
			termLen = len(cradName)
		}
		pos := strings.ToLower(cc.CradAry[i][0:termLen]) >= strings.ToLower(term)
		termLen = oldTermLen
		return pos
	})

	if index < cradAryLen {
		finalIndex := index + 5
		if finalIndex > cradAryLen {
			finalIndex = cradAryLen
		}
		crads := cc.CradAry[index:finalIndex]

		js, err := json.Marshal(crads)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(js)
	} else {
		http.Error(w, "Does not exist", http.StatusNotFound)
		return
	}
}
