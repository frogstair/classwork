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
	// Load in all environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Could not find .env file!")
	}
}

// TestCreateUser tests creating mutliple users
func TestCreateUser(t *testing.T) {

	// Get the database
	db := database.GetPostgres()
	defer database.Disconnect()

	// Test user to create
	type user struct {
		email string
		fname string
		lname string
		passw string
	}

	// List of different users
	testCases := []struct {
		user   user
		passes bool
	}{
		{
			// Valid user
			user: user{
				email: "testmail@gmail.com",
				fname: "Adam",
				lname: "Tester",
				passw: "thisisatest123",
			},
			passes: true,
		},
		{
			// Invalid email
			user: user{
				email: "testmailgmail.com",
				fname: "Adam",
				lname: "Tester",
				passw: "thisisatest123",
			},
			passes: false,
		},
		{
			// Empty firstname
			user: user{
				email: "testmail@gmail.com",
				fname: "      ",
				lname: "Tester",
				passw: "thisisatest123",
			},
			passes: false,
		},
		{
			// Last name too short
			user: user{
				email: "testmail@gmail.com",
				fname: "Adam",
				lname: "T",
				passw: "thisisatest123",
			},
			passes: false,
		},
		{
			// Password too short
			user: user{
				email: "testmail@gmail.com",
				fname: "Adam",
				lname: "T",
				passw: "short",
			},
			passes: false,
		},
	}

	// For each test case
	for i, c := range testCases {
		rUser := models.RegisterUser{
			Email:     c.user.email,
			Password:  c.user.passw,
			FirstName: c.user.fname,
			LastName:  c.user.lname,
		}

		// Register the user
		code, resp := rUser.Register(db)
		// Check error
		if code != 201 {
			// Check if there was meant to be an error
			if c.passes {
				t.Fatalf("Test case %d: error %s", i, resp.Error)
			}
			// If it was meant to fail then no reason to delete it so skip to next loop
			continue
		}
		// If it passed but wasnt meant to
		if code == 201 && !c.passes {
			t.Fatalf("Test case %d: error %s", i, "test succeeded when shouldnt have")
		}

		// Since the struct returned from the Register function
		// is a private struct, I convert it to a hashmap
		// and get the values
		// The struct that is converted is the Response struct
		// which has a Data field, which contains user information.
		// The user information contains an ID field which I will access
		// and use to retreive all the information about the user
		// I cannot make a cast because type information about the struct is lost when
		// its reference is placed in a different struct
		m := structs.Map(resp)
		id := m["Data"].(map[string]interface{})["ID"]

		// Get the user from the database
		user := new(models.User)
		err := db.Where("id = ?", id).First(user).Error
		defer db.Delete(user)

		// Check for any errors
		if err != nil {
			if util.IsNotFoundErr(err) {
				t.Fatalf("Test case %d: error %s", i, "ID doesnt match, user wasnt created")
			}
			t.Fatalf("Test case %d: error %s", i, err.Error())
		}

		// Check that all data was preserved
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
