package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/urfave/negroni"

	"github.com/edwinnduti/file-info/model"
)

// declare database, err and their types
var db *sql.DB
var err error
var config model.Config

// init function
func init() {
	fmt.Println("Getting configs...")
	err = godotenv.Load()
	sendResponse(err, http.StatusInternalServerError, "Error getting configs")

	fmt.Println("We are getting the env values")

	// secret keys
	config.Host = os.Getenv("HOST")
	config.Dbport = os.Getenv("DBPORT")
	config.Dbusername = os.Getenv("USER")
	config.Dbname = os.Getenv("DBNAME")
	config.Passwd = os.Getenv("PASSWORD")

	fmt.Println("Config keys captured successfully!")
}

// Open up our database connection.
func getConnection() *sql.DB {
	db, err = sql.Open("mysql", config.Dbusername+":"+config.Passwd+"@tcp("+config.Host+":"+config.Dbport+")/"+config.Dbname+"?charset=utf8")
	sendResponse(err, http.StatusInternalServerError, "Error connecting to database")
	// if locally running, use this:
	// db ,err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/dbname")
	return db
}

// send response to client
func sendResponse(err error, statusCode int, message string) {
	if err != nil {
		log.Fatalf("%v\n", err)
		var resWriter http.ResponseWriter

		resWriter.WriteHeader(http.StatusOK)
		response := model.Response{
			Code:    statusCode,
			Message: message,
		}
		// give response to client
		json.NewEncoder(resWriter).Encode(response)

	}
}

// homePage is the handler for the root path
func HomePage(w http.ResponseWriter, r *http.Request) {
	// response status code
	w.WriteHeader(http.StatusOK)

	// response message
	w.Write([]byte("Hello world, welcome to the file info API."))
}

// UploadFile is the handler for the upload path
func UploadFile(w http.ResponseWriter, r *http.Request) {
	// get the fileheader from the request and dont save the file
	_, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}

	// connect to database
	db := getConnection()

	// close connection later
	defer db.Close()

	// set the sql insert statement
	sqlStatement := `INSERT INTO property (name, size, extension, type) VALUES (?, ?, ?, ?)`

	// new file properties
	property := model.Property{
		Name:      fileHeader.Filename,
		Size:      fileHeader.Size,
		Type:      fileHeader.Header.Get("Content-Type"),
		Extension: filepath.Ext(fileHeader.Filename),
	}

	result, err := db.Exec(sqlStatement, property.Name, property.Size, property.Extension, property.Type)
	sendResponse(err, http.StatusInternalServerError, "Error executing the insert statement")

	// set id to the new file
	property.ID, err = result.LastInsertId()
	sendResponse(err, http.StatusInternalServerError, "Error returning last inserted id")

	// response status code
	w.WriteHeader(http.StatusOK)

	// response message is the file properties
	json.NewEncoder(w).Encode(property)

}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", HomePage).Methods("GET")
	router.HandleFunc("/upload", UploadFile).Methods("POST")

	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "8080"
	}

	// establish logger

	n := negroni.Classic()
	n.UseHandler(router)
	var server = &http.Server{
		Addr:    ":" + PORT,
		Handler: router,
	}

	server.ListenAndServe()

}
