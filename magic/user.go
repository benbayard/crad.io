package crad

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"regexp"
	"time"
)

type User struct {
	Id                bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	Name              string        `json:"name"      bson:"name"`
	Email             string        `json:"email"     bson:"email"`
	Username          string        `json:"username"  bson:"username"`
	EncryptedPassword []byte        `json:"-"         bson:"encrypted_password"`
	Errors            []string      `json:"errors"    bson:"-"`
	Valid             bool          `json:"-"         bson:"-"`
	Decks             []Deck        `json:"decks"     bson:"decks"`
}

const (
	MinNameLength     = 5
	MinUsernameLength = 8
	EmailRegex        = "^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$"
	TokenSign         = "I love Things That Work This Well! WHAT"
)

var GlobalConnection DBConnection

// func databaseConnect() (*mgo.Session, *mgo.Collection) {
// 	session, err := mgo.Dial("localhost:27017")
// 	index := mgo.Index{
// 		Key:        []string{"email"},
// 		Unique:     true,
// 		Background: true,
// 	}
// 	unIndex := mgo.Index{
// 		Key:        []string{"username"},
// 		Unique:     true,
// 		Background: true,
// 	}

// 	deckIndex := mgo.Index{
// 		Key:        []string{"_id decks.name"},
// 		Unique:     true,
// 		Background: true,
// 	}

// 	// username
// 	userCollection := session.DB("crad").C("users")
// 	err = userCollection.EnsureIndex(index)
// 	err = userCollection.EnsureIndex(unIndex)
// 	err = userCollection.EnsureIndex(deckIndex)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return session, userCollection
// }

func DatabaseConnect() {
	session, err := mgo.Dial("localhost:27017")
	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		Background: true,
	}
	unIndex := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		Background: true,
	}

	deckIndex := mgo.Index{
		Key:        []string{"_id decks.name"},
		Unique:     true,
		Background: true,
	}

	// username
	userCollection := session.DB("crad").C("users")
	err = userCollection.EnsureIndex(index)
	err = userCollection.EnsureIndex(unIndex)
	err = userCollection.EnsureIndex(deckIndex)
	if err != nil {
		panic(err)
	}

	GlobalConnection.Session = session
	GlobalConnection.Collection = userCollection
}

// you use clear because it removed the password in plain text from memory.
func clear(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}

func UserNew(name string, email string, username string, password string) (*User, string) {
	defer clear([]byte(password))
	_, userCollection := GlobalConnection.GetDB()
	// // defer session.Close()

	if validatePassword(password) == false {
		user := &User{}
		user.Errors = append(user.Errors, "Password is not long enough")
		return user, ""
	}

	encryptedPassword, _ := encryptPassword([]byte(password))

	user := &User{
		Name:              name,
		Email:             email,
		Username:          username,
		EncryptedPassword: encryptedPassword,
		Valid:             true,
		Id:                bson.NewObjectId(),
	}

	user.Validate()

	if !user.Valid {
		return user, ""
	}
	err := userCollection.Insert(user)

	if err != nil {
		user.Valid = false
		user.Errors = append(user.Errors, err.Error())
	}

	return user, user.createToken()
}

func validatePassword(password string) bool {
	defer clear([]byte(password))
	if len(password) < 6 {
		return false
	} else {
		return true
	}
}

func (u *User) Validate() bool {
	vn := u.validateName()
	ve := u.validateEmail()
	vu := u.validateUsername()

	u.Valid = (vn && ve && vu)

	return u.Valid
}

func (u *User) validateName() bool {
	isValid := true
	if len(u.Name) < MinNameLength {
		u.Errors = append(u.Errors, fmt.Sprintf("Name must be at least %d characters long.", MinNameLength))
		isValid = false
	}

	if match, _ := regexp.MatchString("\\d", u.Name); match == true {
		u.Errors = append(u.Errors, "Name must not include numbers.")
		isValid = false
	}

	return isValid
}

func (u *User) validateEmail() bool {
	isValid := true

	if match, _ := regexp.MatchString(EmailRegex, u.Email); match == false {
		u.Errors = append(u.Errors, "Email \""+u.Email+"\" is invalid. Please use a valid email.")
		isValid = false
	}

	return isValid
}

func (u *User) validateUsername() bool {
	isValid := true
	if len(u.Username) < MinUsernameLength {
		u.Errors = append(u.Errors, fmt.Sprintf("Username must be at least %d characters long.", MinUsernameLength))
		isValid = false
	}

	return isValid
}

func encryptPassword(password []byte) ([]byte, error) {
	defer clear(password)
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// func decodeFunction(token *jwt.Token) ([]byte, error) {
//   return []byte(TokenSign), nil
// }

func decodeToken(r *http.Request) (interface{}, error) {
	// jwt.ParseFromRequest(r, func(token *jwt.Token) ([]byte, error) {
	//     return publicKey, nil
	//   })
	token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
		return []byte(TokenSign), nil
	})

	// fmt.Printf("Token: %#v \n", token)
	// fmt.Printf("Error: %#v \n", err)

	if err != nil || token.Valid == false {
		return "", err
	}

	return token.Claims["CustomUserInfo"], err
}

func (u *User) createToken() string {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	// expire after 3 days
	token.Claims["CustomUserInfo"] = u.Id

	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(TokenSign))

	if err != nil {
		u.Valid = false
		fmt.Printf("Error: \n")
		fmt.Printf("%#v \n", err)
	}

	return tokenString
}

func (u *User) ValidToken(r *http.Request) (bool, error) {
	decoded, err := decodeToken(r)
	id, _ := decoded.(string)

	if err != nil {
		return false, err
	}

	// fmt.Printf("Decoded: %#v \n", decoded)

	// fmt.Printf("(%#v == %#v)", bson.ObjectIdHex(id), u.Id)

	return (bson.ObjectIdHex(id) == u.Id), err
}

func UserById(id string) User {
	_, userCollection := GlobalConnection.GetDB()
	// defer session.Close()

	user := User{}

	userCollection.FindId(bson.ObjectIdHex(id)).One(&user)

	return user
}

func UserByUsername(username string) User {
	_, userCollection := GlobalConnection.GetDB()
	// defer session.Close()

	user := User{}

	userCollection.Find(bson.M{"username": username}).One(&user)

	return user
}

func UserByEmail(email string) User {
	user := User{}
	_, userCollection := GlobalConnection.GetDB()
	// defer session.Close()

	userCollection.Find(bson.M{"email": email}).One(&user)
	return user
}

func (user *User) CorrectPassword(password string) error {
	defer clear([]byte(password))
	return bcrypt.CompareHashAndPassword(user.EncryptedPassword, []byte(password))
}

func (user *User) Save(changes map[string]string) {
	_, userCollection := GlobalConnection.GetDB()
	// defer session.Close()

	changed := false

	if changes["email"] != "" {
		user.Email = changes["email"]
		changed = true
	}

	if changes["username"] != "" {
		user.Username = changes["username"]
		changed = true

	}

	if changes["name"] != "" {
		user.Name = changes["name"]
		changed = true
	}

	if changed == true {
		user.Validate()

		if !user.Valid {
			return
		}

		err := userCollection.UpdateId(user.Id, user)

		if err != nil {
			user.Valid = false
			user.Errors = append(user.Errors, err.Error())
		}
	}

	fmt.Printf("User: %#v \n", user)
}
