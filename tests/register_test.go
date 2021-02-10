package tests

import (
	"classwork/database"
	"classwork/models"
	"classwork/util"
	"log"
	"testing"

	"github.com/fatih/structs"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Could not find .env file!")
	}
}

// TestCreateUser tests creating mutliple users
func TestCreateUser(t *testing.T) {

	db := database.GetPostgres()
	defer database.Disconnect()

	type user struct {
		email string
		fname string
		lname string
		passw string
	}

	testCases := []struct {
		user   user
		passes bool
	}{
		{
			user: user{
				email: "testmail@gmail.com",
				fname: "Adam",
				lname: "Tester",
				passw: "thisisatest123",
			},
			passes: true,
		},
		{
			user: user{
				email: "testmailgmail.com",
				fname: "Adam",
				lname: "Tester",
				passw: "thisisatest123",
			},
			passes: false,
		},
		{
			user: user{
				email: "testmail@gmail.com",
				fname: "      ",
				lname: "Tester",
				passw: "thisisatest123",
			},
			passes: false,
		},
		{
			user: user{
				email: "testmail@gmail.com",
				fname: "Adam",
				lname: "T",
				passw: "thisisatest123",
			},
			passes: false,
		},
		{
			user: user{
				email: "testmail@gmail.com",
				fname: "Adam",
				lname: "T",
				passw: "short",
			},
			passes: false,
		},
	}

	for i, c := range testCases {
		rUser := models.RegisterUser{
			Email:     c.user.email,
			Password:  c.user.passw,
			FirstName: c.user.fname,
			LastName:  c.user.lname,
		}

		code, resp := rUser.Register(db)
		if code != 201 {
			if c.passes {
				t.Fatalf("Test case %d: error %s", i, resp.Error)
			}
			return
		}
		if code == 201 && !c.passes {
			t.Fatalf("Test case %d: error %s", i, "test succeeded when shouldnt have")
		}

		m := structs.Map(resp)
		id := m["Data"].(map[string]interface{})["ID"]

		user := new(models.User)
		err := db.Where("id = ?", id).First(user).Error
		defer db.Delete(user)

		if err != nil {
			if util.IsNotFoundErr(err) {
				t.Fatalf("Test case %d: error %s", i, "first name doesnt match")
			}
			t.Fatalf("Test case %d: error %s", i, err.Error())
		}

		if user.FirstName != c.user.fname {
			t.Fatalf("Test case %d: error %s", i, "first name doesnt match")
		}
		if user.LastName != c.user.lname {
			t.Fatalf("Test case %d: error %s", i, "last name doesnt match")
		}
		if user.Email != c.user.email {
			t.Fatalf("Test case %d: error %s", i, "first name doesnt match")
		}
		if !util.Compare(user.Password, c.user.passw) {
			t.Fatalf("Test case %d: error %s", i, "password does not match")
		}
	}
}
