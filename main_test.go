package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize("root", "godfather", "rest_api_example")

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)

}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestCreateUser(t *testing.T) {
	clearTable()
	// create a payload with fake user data
	payload := []byte(`{"name": "test user", "age": 30}`)
	// create a POST request to the /user endpoint and attach our userdata payload
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	// execute the request and capture the request
	response := executeRequest(req)
	// make sure that the response code is a a successful one (201)
	checkResponseCode(t, http.StatusCreated, response.Code)
	// create a map variable
	var m map[string]interface{}
	// unmarshal parses the json-encoded data and stores the result in the value pointed to by v
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test user" {
		t.Errorf("Expected user name to be `test user`. Got `%v'", m["name"])
	}
	if m["age"] != "30" {
		t.Errorf("Expected user age to be `30`. Got `%v'", m["age"])
	}
	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	//floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected user ID to be `1`. Got `%v`", m["id"])
	}
}

// The following test, tests two things, the status code 404 and if the reponse contains the expected error message
func TestGetNonExistentUser(t *testing.T) {
	// clear the users table
	clearTable()
	// create a request for user 45 on the users table
	req, _ := http.NewRequest("GET", "/user/45", nil)
	// execute the request and capture the result
	response := executeRequest(req)
	// make sure that the result is a status not found one
	checkResponseCode(t, http.StatusNotFound, response.Code)
	// create an empty map
	var m map[string]string
	// parses the encoded json and stores the value in the second variable
	json.Unmarshal(response.Body.Bytes(), &m)
	// fail the test if the error is not what is expected.
	if m["error"] != "User not found" {
		t.Errorf("Expected the error key of the response to be set to `User not found`. Got `%s`", m["error"])
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func clearTable() {
	a.DB.Exec("DELETE FROM users")
	a.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS users(
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	age INT NOT NULL
)`
