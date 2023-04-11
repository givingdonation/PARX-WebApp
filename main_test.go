package main

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHash(t *testing.T) {
	rand.Seed(time.Now().Unix())
	var passwords []string
	for i := 0; i < 50; i++ {
		var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321!@#$%^&*()"
		str := make([]byte, rand.Intn(16)+1)
		for k := range str {
			str[k] = chars[rand.Intn(len(chars))]
		}
		passwords = append(passwords, string(str))
		result := hashPswd(passwords[i])
		expected := hash(passwords[i])
		if result != expected {
			t.Errorf("Incorrect hash for %s FAILED. Expected %s, got %s\n", passwords[i], expected, result)
			return
		}
	}
	t.Logf("All hashes as expected.")
}

func TestEvents(t *testing.T) {
	var Events []eventInfo
	for i := 0; i < 10; i++ {
		Events = append(Events, eventInfo{
			Points:              0,
			EventDescription:    "asdf",
			EventDate:           "2017-06-01",
			RoomNumber:          0,
			AdvisorNames:        "asdf",
			Location:            "asdf",
			LocationDescription: "asdf",
			Sport:               "asdf",
			SportDescription:    "asdf",
			EventImage:          "https://imgs.search.brave.com/ToRVheIVFOHdWRebW6v6BriMZf_slwrqoAXvU-I62CY/rs:fit:1200:1200:1/g:ce/aHR0cHM6Ly90aGV3/b3dzdHlsZS5jb20v/d3AtY29udGVudC91/cGxvYWRzLzIwMTUv/MDEvbmF0dXJlLWlt/YWdlcy4uanBn",
			StudentName:         "asdf",
			StudentNumber:       0,
			StudentAttended:     true,
		})
	}
	httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		err := tplExec(writer, "teacher_events.gohtml", Events)
		if err != nil {
			t.Errorf("Failed to load events")
		}
	}))
	t.Logf("All events loaded.")
}
func TestHome(t *testing.T) {
	//SEMI-SCUFFED WAY OF MAKING THE USER NOT BE ABLE TO ACCESS HOME IF NOT LOGGED IN, CONSIDER USING COOKIES

	//Here we should populate the rest of the userInfo struct with sql queries and load whatever else we need for the home page.
	//Also, we need to find out how to get signup to upload to db and login to get
	//We can probably just do different interactions for get/post requests to the home, same way we did

	httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		err := tplExec(writer, "home.gohtml", homeData{
			Name:          "a",
			Grade:         9,
			Points:        10,
			Grade9Points:  []int{1, 2, 3, 4, 5},
			Grade10Points: []int{1, 2, 3, 4, 5},
			Grade11Points: []int{1, 2, 3, 4, 5},
			Grade12Points: []int{1, 2, 3, 4, 5},
		})
		if err != nil {
			t.Errorf("Failed to load home")
		}
	}))
	t.Logf("All points loaded.")
}
func TestTplExec(t *testing.T) {
	templateNames := []string{"error.gohtml", "login.gohtml", "signup.gohtml", "teacher_events.gohtml"}
	for i := 0; i < len(templateNames); i++ {
		httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			err := tplExec(writer, templateNames[i], nil)
			if err != nil {
				t.Errorf("Unable to load template %s", templateNames[i])
				return
			}
		}))
	}
	t.Logf("All templates loaded.")
}

func TestDataValidation(t *testing.T) {
	rand.Seed(time.Now().Unix())
	testData := userData{}
	var expected bool
	requestMethod := "signup"

	for i := 0; i < 1000000; i++ {
		if rand.Intn(2) != 0 {
			requestMethod = "login"
		}

		expected = false

		var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321!@#$%^&*()"
		str := make([]byte, rand.Intn(16)+1)
		for k := range str {
			str[k] = chars[rand.Intn(len(chars))]
		}
		testData.passwordHash = hashPswd(string(str))
		testData.IdNumber = rand.Intn(9999998) + 1

		if requestMethod == "signup" {
			testData.Grade = rand.Intn(4) + 9
			testData.Name = string(str)
		}

		punishment := rand.Intn(11)

		if punishment == 1 && requestMethod == "signup" {
			testData.Grade = 13
		} else if punishment == 3 && requestMethod == "signup" {
			testData.Name = ""
		} else if punishment == 5 {
			testData.passwordHash = hashPswd("")
		} else if punishment == 7 {
			testData.IdNumber = 999999999999
		} else {
			expected = true
		}

		if checkData(requestMethod, &testData) != expected {
			t.Errorf("DID NOT STOP A BAD INPUT. Expected %t using method %s on iteration %d using random number %d", expected, requestMethod, i, punishment)
			return
		}
	}
	t.Logf("STOPPED BAD INPUTS")
}

func hash(pwd string) string {

	hasher := sha256.New()
	hasher.Write([]byte(pwd))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}
