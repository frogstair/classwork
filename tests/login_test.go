package tests

import (
	"classwork/database"
	m "classwork/models"
	"log"
	"testing"

	"github.com/fatih/structs"
	"github.com/joho/godotenv"
)

// Load in the environment
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Could not find .env file!")
	}
}

func TestLogin(t *testing.T) {
	// Get the database and disconnect when done
	db := database.GetPostgres()
	defer database.Disconnect()

	// The testing struct
	type login struct {
		email    string
		password string
	}

	// User must be registered before logging in
	type register struct {
		email string
		fname string
		lname string
		passw string
	}

	// List of test cases
	testcases := []struct {
		login    login
		register register
		passes   bool
	}{
		{
			login: login{
				email:    "daniladudkin412@gmail.comm",
				password: "OdcpbY4KTcpNUtHQ1oPI",
			},
			register: register{
				email: "daniladudkin412@gmail.comm",
				passw: "OdcpbY4KTcpNUtHQ1oPI",
				fname: "Danila",
				lname: "Dudkin",
			},
			passes: true,
		},
		{
			login: login{
				email:    "daniladudkin412@gmail.com",
				password: "tl08PDebEK5bOGryFj8a",
			},
			register: register{
				email: "daniladudkin412@gmail.co",
				passw: "tl08PDebEK5bOGryFj8a",
				fname: "Danila",
				lname: "Dudkin",
			},
			passes: false,
		},
		{
			login: login{
				email:    "daniladudkin412@gmail.commm",
				password: "xTVJ1kpluOuhPi25oDbD",
			},
			register: register{
				email: "daniladudkin412@gmail.commm",
				passw: "xTVJ1kpluOuhPi25oDb",
				fname: "Danila",
				lname: "Dudkin",
			},
			passes: false,
		},
	}

	for i, c := range testcases {

		// Register the user
		// 100% guaranteed success
		regUser := m.RegisterUser{
			FirstName: c.register.fname,
			LastName:  c.register.lname,
			Email:     c.register.email,
			Password:  c.register.passw,
		}
		_, resp := regUser.Register(db)
		mp := structs.Map(resp)
		id := mp["Data"].(map[string]interface{})["ID"]

		// Get the user
		user := new(m.User)
		db.Where("id = ?", id).First(user)
		defer db.Delete(user)

		// Login the user
		logUser := m.LoginUser{
			Email:    c.login.email,
			Password: c.login.password,
		}
		code, resp, tok := logUser.Login(db)
		// Check if test was supposed to be passed
		if code != 200 {
			if c.passes {
				t.Fatalf("Test case %d: error %s", i, resp.Error)
			}
			return
		}
		if code == 200 && !c.passes {
			t.Fatalf("Test case %d: error %s", i, "test succeeded when shouldnt have")
		}

		// Parse the token the user got along with the claims
		_, resp = m.ParseToken(tok, db)
		mp = structs.Map(resp)
		logid := mp["Data"].(map[string]interface{})["ID"]

		// Check the ID to see if it matches
		if logid != id && c.passes {
			t.Fatalf("Test case %d: error %s", i, "user id not present in token")
		}

	}

}
