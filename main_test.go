package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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
	// unmarshal parses the json-encoded data and stores the result in the value pointed to by m
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

func TestGetUser(t *testing.T) {
	clearTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/user/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func addUsers(count int) {
	if count < 1 {
		count = 1
	}
	// create a request
	for i := 0; i < count; i++ {
		statement := fmt.Sprintf("INSERT INTO users(name, age) VALUES(`%s`,`%d`)", "User "+strconv.Itoa(i+1), ((i + 1) * 10))
		// execute the request
		a.DB.Exec(statement)
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

func TestUpdateUser(t *testing.T) {
	clearTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/user/1", nil)
	response := executeRequest(req)

	var originalUser map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalUser)

	payload := []byte(`{"name":"test user - updated name", "age": 21}`)

	req, _ = http.NewRequest("PUT", "/user/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalUser["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalUser["id"], m["id"])
	}

	if m["name"] != originalUser["name"] {
		t.Errorf("Expected the name to change from `%v` to `%v`. Got %v", originalUser["name"], m["id"], m["id"])
	}

	if m["age"] != originalUser["age"] {
		t.Errorf("Expected the age to change from `%v` to `%v`. Got %v", originalUser["age"], m["age"], m["age"])
	}
}

func TestDeleteUser(t *testing.T) {
	clearTable()
	addUsers(1)
	// get the user data which was just created
	req, _ := http.NewRequest("GET", "/user/1", nil)
	// get the resopnse
	response := executeRequest(req)
	// make sure we have a correct response code
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/user/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/user/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

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
