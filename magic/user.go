package crad

import (
	"code.google.com/p/go.crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
  "gopkg.in/mgo.v2"
	"regexp"
  "github.com/dgrijalva/jwt-go"
  "time"
)

type User struct {
	Id                bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	Name              string        `json:"name"      bson:"name"`
	Email             string        `json:"email"     bson:"email"`
	Username          string        `json:"username"  bson:"username"`
	EncryptedPassword []byte        `bson:"encrypted_password"`
	Errors            []string      `json:"errors"`
	Valid             bool          `json:"valid"`
}

const (
	MinNameLength     = 5
	MinUsernameLength = 8
	EmailRegex        = "^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$"
  TokenSign         = "I love Things That Work This Well! WHAT"
)

func databaseConnect() *mgo.Session {
  session, err := mgo.Dial("localhost:27017")
  if err != nil {
    panic(err)
  }
  return session
}

// you use clear because it removed the password in plain text from memory.
func clear(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}

func CreateUser(name string, email string, username string, password string) (*User, string) {
	defer clear([]byte(password))
  session := databaseConnect()
  defer session.Close()
  userCollection := session.DB("crad").C("users")


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
	}

  user.Validate()

  if !user.Valid {
    return user, ""
  } else {
    userCollection.Insert(user)
  }

	return user, user.createToken()
}

func (u *User) createToken() string {
  token := jwt.New(jwt.GetSigningMethod("HS256"))
  token.Claims["userId"] = u.Id
  // expire after 3 days
  token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
  // Sign and get the complete encoded token as a string
  tokenString, _ := token.SignedString(TokenSign)

  return tokenString  

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
	u.validateName()
	u.validateEmail()
	u.validateUsername()

	return u.Valid
}

func (u *User) validateName() {
	isValid := true
	if len(u.Name) < MinNameLength {
		u.Errors = append(u.Errors, "Username must be at least "+string(MinNameLength)+"characters long.")
		isValid = false
	}

	if match, _ := regexp.MatchString("\\d", u.Name); match == true {
		u.Errors = append(u.Errors, "Username must not include numbers.")
		isValid = false
	}

	u.Valid = isValid
}

func (u *User) validateEmail() {
	isValid := true

	if match, _ := regexp.MatchString(EmailRegex, u.Email); match == false {
		u.Errors = append(u.Errors, "Email \""+u.Email+"\" is invalid. Please use a valid email.")
		isValid = false
	}

	u.Valid = isValid
}

func (u *User) validateUsername() {
	isValid := false
	if len(u.Username) < MinUsernameLength {
		u.Errors = append(u.Errors, "Username must be at least "+string(MinUsernameLength)+"characters long.")
		isValid = false
	}

	u.Valid = isValid
}

func encryptPassword(password []byte) ([]byte, error) {
	defer clear(password)
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}
