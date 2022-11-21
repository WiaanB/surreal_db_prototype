package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	surrealdb "github.com/surrealdb/surrealdb.go"
)

var DB *surrealdb.DB

func main() {
	fmt.Println("SurrelDB Prototype running...")

	var err error
	fmt.Println("Setting up the end information...")
	err = godotenv.Load(".env")
	handleErr(err)

	DB, err = surrealdb.New("ws://" + os.Getenv("URL") + ":" + os.Getenv("PORT") + "/" + os.Getenv("METHOD"))
	handleErr(err)

	_, err = DB.Signin(map[string]interface{}{
		"user": os.Getenv("USR"),
		"pass": os.Getenv("PASS"),
	})
	handleErr(err)

	_, err = DB.Use("test", "test")
	handleErr(err)

	createRecord(map[string]interface{}{
		"name":    "Wiaan",
		"surname": "Botha",
		"age":     24,
	})
	createRecord(map[string]interface{}{
		"name":    "Guy",
		"surname": "Gibson",
		"age":     22 + 8,
	})

	getRecords()

	deleteRecords()
}

func createRecord(data map[string]interface{}) {
	_, err := DB.Create("company", data)
	handleErr(err)
}

func getRecords() {
	any, err := DB.Select("company")
	handleErr(err)
	fmt.Println(any)
}

func deleteRecords() {
	_, err := DB.Delete("company")
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
