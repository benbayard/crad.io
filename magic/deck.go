package crad

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type Deck struct {
	Id          bson.ObjectId       `json:"id" bson:"_id,omitempty"`
	Name        string              `json:"name" bson:"name"`
	Description string              `json:"description" bson: "name"`
	Crads       map[string]DeckCrad `json:"crads" bson: "crads"`
	Format      *string             `json:"format" bson:"format"`
	Errors      []string            `json:"errors" bson:"-"`
	Valid       bool                `json:"-" bson:"-"`
}

type DeckCrad struct {
	Quantity int    `json:"quantity" bson:"quantity"`
	CradName string `json:"name" bson: "cradname"`
}

var (
	formats map[string]string = map[string]string{
		"legacy":     "legacy",
		"commander":  "commander",
		"cradmander": "cradmander",
		"standard":   "standard",
		"modern":     "modern",
		"vintage":    "vintage",
	}
)

const (
	MinDeckNameLength = 5
)

func (u *User) DeckNew(name string, description string, format string) (deck *Deck) {
	foundFormat := formats[format]
	deck = &Deck{
		Id:          bson.NewObjectId(),
		Name:        name,
		Description: description,
		Format:      &foundFormat,
	}

	deck.Validate()

	if deck.Valid == true {
		u.Decks = append(u.Decks, *deck)
		u.Validate()

		if u.Valid {
			session, userCollection := databaseConnect()
			defer session.Close()

			err := userCollection.UpdateId(u.Id, u)

			if err != nil {
				u.Valid = false
				u.Errors = append(u.Errors, err.Error())
			}
		}

		return deck
	}

	for _, element := range deck.Errors {
		u.Errors = append(u.Errors, element)
	}

	return &Deck{}

}

func (d *Deck) Validate() {
	vn := d.validateName()
	vf := d.validateFormat()

	d.Valid = (vn && vf)
}

func (d *Deck) validateName() bool {
	isValid := true
	if len(d.Name) < MinDeckNameLength {
		d.Errors = append(d.Errors, fmt.Sprintf("Name must be at least %d characters long.", MinNameLength))
		isValid = false
	}

	return isValid
}

func (d *Deck) validateFormat() bool {
	isValid := true
	format := formats[*d.Format]

	if format == "" {
		d.Errors = append(d.Errors, "Invalid Format")
		isValid = false
	}

	return isValid
}

func (u *User) DeckByName(name string) int {
	for index, d := range u.Decks {
		if d.Name == name {
			return index
		}
	}

	return -1
}
