package crad

import (
	"encoding/json"
	// "fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

type UserController struct{}

func (uc *UserController) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	DatabaseConnect()
	// fmt.Printf("Global Connection: ", GlobalConnection)

	_, userCollection := GlobalConnection.GetDB()

	var result []*User
	iter := userCollection.Find(nil).Iter()
	err := iter.All(&result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func (uc *UserController) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpUser := make(map[string]map[string]string)
	defer clear([]byte(tmpUser["user"]["password"]))

	json.Unmarshal(body, &tmpUser)

	user, token := UserNew(tmpUser["user"]["name"], tmpUser["user"]["email"], tmpUser["user"]["username"], tmpUser["user"]["password"])

	if len(token) < 5 {
		// js, _ := json.Marshal(user.Errors)
		w.Header().Set("Status", string(http.StatusNotAcceptable))
		// w.Write(js)
		// return
	}

	js, err := json.Marshal(user)
	w.Write(js)
	return
	// fmt.Printf("Request Body: %#v \n", tmpUser)

}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	creds := make(map[string]string)
	defer clear([]byte(creds["password"]))

	json.Unmarshal(body, &creds)

	password, pExists := creds["password"]
	email, eExists := creds["email"]

	if (pExists && eExists) == false {
		http.Error(w, "Email or Password not included", http.StatusNotAcceptable)
		return
	}

	user := UserByEmail(email)

	err = user.CorrectPassword(password)

	if err != nil {
		http.Error(w, "Email or Password is incorrect", http.StatusNotAcceptable)
		return
	}

	type UserToken struct {
		User  User   `json:"user"`
		Token string `json:"token"`
	}

	newUser := UserToken{user, user.createToken()}

	js, err := json.Marshal(newUser)

	w.Write(js)

}

func (uc *UserController) Show(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId := ps.ByName("id")

	user := UserById(userId)

	js, err := json.Marshal(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func (uc *UserController) ValidToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")

	user := UserByUsername(username)
	// okay so we have the  user, now we just have to approve the user.
	valid, err := user.ValidToken(r)

	if err != nil || !valid {
		http.Error(w, "Incorrect Token", http.StatusForbidden)
		return
	}

	js, err := json.Marshal(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func (uc *UserController) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	// fmt.Println("Line 128")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println("Line 134")

	tmpUser := make(map[string]map[string]string)
	json.Unmarshal(body, &tmpUser)

	// fmt.Println("Line 138")

	userId := ps.ByName("id")

	user := UserById(userId)
	// okay so we have the  user, now we just have to approve the user.
	valid, err := user.ValidToken(r)

	if err != nil || !valid {
		http.Error(w, "Incorrect Token", http.StatusForbidden)
		return
	}

	newUser, _ := tmpUser["user"]

	if newUser != nil {
		user.Save(newUser)
	}

	// fmt.Printf("User: %#v \n", user)

	js, err := json.Marshal(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)

}
