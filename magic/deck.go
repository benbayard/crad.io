package crad

import (
	"gopkg.in/mgo.v2/bson"
)

type Deck struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson: "name"`
	Crads       []DeckCrad    `json:"crads" bson: "crads"`
	Format      *string       `json:"format" bson:"format"`
  Errors      []string      `json:"errors" bson:"-"`
  Valid       bool          `json:"-" bson:"-"`
}

type DeckCrad struct {
	Quantity int    `json:"quantity" bson:"quantity"`
	Crad     *Crad  `json:"crad" bson:"-"`
	CradName string `json:"-" bson: "cradname"`
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

const(
  MinNameLength= 5
)

func (u *User) DeckNew(name string, description string, format string) {

}

func (d *DeckCrad) Validate() {
  vn := d.validateName()

  d.Valid := 
}

func (d *DeckCrad) validateName() bool {
  if len(d.Name) < MinNameLength {
    d.Errors = append(d.Errors, fmt.Sprintf("Name must be at least %d characters long.", MinNameLength))
    isValid = false
  }

  return isValid
}


