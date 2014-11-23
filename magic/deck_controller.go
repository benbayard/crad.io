package crad

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

type DeckController struct {
	Crads map[string]Crad
}

func (dc *DeckController) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpDeck := make(map[string]map[string]string)
	json.Unmarshal(body, &tmpDeck)

	userId := ps.ByName("username")

	user := UserByUsername(userId)
	// okay so we have the  user, now we just have to approve the user.
	valid, err := user.ValidToken(r)

	if err != nil || !valid {
		http.Error(w, "Incorrect Token", http.StatusForbidden)
		return
	}

	newDeck := tmpDeck["deck"]

	if newDeck == nil {
		http.Error(w, "Incorrect Data Structure", http.StatusNotAcceptable)
	}

	_ = user.DeckNew(newDeck["name"], newDeck["description"], newDeck["format"])

	js, err := json.Marshal(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)

}

func (dc *DeckController) AddCrad(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpDeck := make(map[string]map[string]interface{})

	// fmt.Printf("Deck: %#v \n", body)

	err = json.Unmarshal(body, &tmpDeck)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId := ps.ByName("username")

	user := UserByUsername(userId)
	// okay so we have the  user, now we just have to approve the user.
	valid, err := user.ValidToken(r)

	if err != nil || !valid {
		http.Error(w, "Incorrect Token", http.StatusForbidden)
		return
	}

	fmt.Printf("Deck: %#v \n", tmpDeck["deck"])
	deckIndex := user.DeckByName(tmpDeck["deck"]["name"].(string))

	if deckIndex == -1 {
		http.Error(w, "Does not exist", http.StatusNotFound)
		return
	}
	deck := user.Decks[deckIndex]

	areCrads := tmpDeck["deck"]["crads"]

	fmt.Printf("Are Crads: %#v \n", areCrads)

	crads := areCrads.([]interface{})

	fmt.Printf("Crads: %#v \n", crads)

	for _, tmpCrad := range crads {
		crad := tmpCrad.(map[string]interface{})
		q := int(crad["quantity"].(float64))
		n := crad["name"].(string)
		fmt.Println("CRADDINGIGGINGIGNGINGIGNISDIFGN")
		if q > 4 {
			q = 4
		}
		if q < 1 {
			q = 1
		}
		_, okay := dc.Crads[n]
		if okay != true {
			fmt.Println("NOT OKAY CRAD")
			deck.Errors = append(deck.Errors, "Crad "+n+" does not exist")
		} else {
			deck.Crads[n] = DeckCrad{
				Quantity: q,
				CradName: n,
			}
		}
	}

	fmt.Printf("Deck Crads: %#v \n", deck.Crads)

	user.Decks[deckIndex] = deck

	user.Validate()

	fmt.Printf("Users: %#v \n", user)

	if user.Valid {
		session, userCollection := databaseConnect()
		defer session.Close()

		err := userCollection.UpdateId(user.Id, user)

		if err != nil {
			user.Valid = false
			user.Errors = append(user.Errors, err.Error())
		}
	}

	js, err := json.Marshal(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)

}

func (dc *DeckController) EditCrad(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
