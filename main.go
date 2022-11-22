package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	surrealdb "github.com/surrealdb/surrealdb.go"
)

var DB *surrealdb.DB

type User struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     int    `json:"age"`
}

// Create a map[string]interface{} from a User struct
func (u *User) toMap() map[string]interface{} {
	var m map[string]interface{}
	data, _ := json.Marshal(u)
	json.Unmarshal(data, &m)
	return m
}

var Users []User

func main() {
	fmt.Println("SurrelDB Prototype running...")

	// Setup env information
	var err error
	err = godotenv.Load(".env")
	handleErr(err)

	// Connect to the websocket for SurrealDB
	DB, err = surrealdb.New("ws://" + os.Getenv("URL") + ":" + os.Getenv("PORT") + "/" + os.Getenv("METHOD"))
	handleErr(err)

	// Signing in to SurrealDB with the defined user
	_, err = DB.Signin(map[string]interface{}{
		"user": os.Getenv("USR"),
		"pass": os.Getenv("PASS"),
	})
	handleErr(err)

	// Use the "test" namespace, and the "test" database
	_, err = DB.Use("test", "test")
	handleErr(err)

	// Create two dummy records into company
	createRecord(User{
		Name:    "Wiaan",
		Surname: "Botha",
		Age:     24,
	})
	createRecord(User{
		Name:    "Guy",
		Surname: "Gibson",
		Age:     22 + 8,
	})

	// Fetch all the records from company
	records := getRecords()

	// Compare the two
	fmt.Printf("Records added to the DB: %v\n", Users)
	fmt.Printf("Records found in the DB: %v\n", records)

	// Delete all the records from company by ID
	for _, record := range records {
		deleteRecord(record)
	}
}

func createRecord(u User) {
	rec, err := DB.Create("company", u.toMap())
	handleErr(err)
	addUser(rec)
}

func getRecords() (u []User) {
	any, err := DB.Select("company")
	handleErr(err)
	data, _ := json.Marshal(any)
	json.Unmarshal(data, &u)
	return
}

func deleteRecord(record User) {
	fmt.Printf("Deleting %v with ID: %v\n", record.Name, record.Id)
	_, err := DB.Delete(record.Id)
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

func addUser(a any) (u User) {
	// any -> map[string]interface{}
	var m []map[string]interface{}
	data, _ := json.Marshal(a)
	json.Unmarshal(data, &m)

	// For each hit, add it to the slice
	for _, v := range m {
		data, _ := json.Marshal(v)
		json.Unmarshal(data, &u)
		Users = append(Users, u)
	}
	return
}
