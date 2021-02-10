package tests

import (
	"classwork/database"
	m "classwork/models"
	"classwork/util"
	"log"
	"testing"

	"github.com/fatih/structs"
	"github.com/joho/godotenv"
)

//
// A HEADMASTER USER MUST EXIST BEFOREHAND
//

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Could not find .env file!")
	}
}

func TestAddSchool(t *testing.T) {

	db := database.GetPostgres()
	defer database.Disconnect()

	user := new(m.User)
	const id = "1oIpNHE7WOEmdlCGuIsYu4R7LGo" // <- REPLACE ID
	db.Where("id = ?", id).First(user)

	tests := []struct {
		school m.NewSchool
		passes bool
	}{
		{
			school: m.NewSchool{
				Name: "Test school 1",
			},
			passes: true,
		},
		{
			school: m.NewSchool{
				Name: "Tes",
			},
			passes: false,
		},
	}

	for i, c := range tests {
		code, resp := c.school.Add(db, user)
		if code != 201 {
			if c.passes {
				t.Fatalf("Test case %d: error %s", i, resp.Error)
			}
			return
		}
		if code == 201 && !c.passes {
			t.Fatalf("Test case %d: error %s", i, "test succeeded when shouldnt have")
		}

		mp := structs.Map(resp)
		id := mp["Data"].(map[string]interface{})["ID"]

		school := new(m.School)
		err := db.Where("id = ?", id).First(school).Error
		if err != nil {
			if util.IsNotFoundErr(err) {
				t.Fatalf("Test case %d: error %s", i, "school wasnt created")
			}
			t.Fatalf("Test case %d: error %s", i, err.Error())
		}
		defer db.Delete(school)

		if school.Name != c.school.Name {
			t.Fatalf("Test case %d: error %s", i, "name doesnt match")
		}
	}

}
