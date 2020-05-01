package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/matancredi/superheroes-api/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateSuper(t *testing.T) {

	err := refreshTables()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		name         string
		errorMessage string
	}{
		{
			inputJSON:    `{"name":"Batman"}`,
			statusCode:   201,
			name:         "Batman",
			errorMessage: "",
		},
		{
			inputJSON:    `{"name":"Batman"}`,
			statusCode:   500,
			name:         "Batman",
			errorMessage: "Super already registered",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/supers", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateSuper)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["name"], v.name)
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetSupers(t *testing.T) {

	err := refreshTables()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedSupers()
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/supers", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetSupers)
	handler.ServeHTTP(rr, req)

	var supers []models.Super
	err = json.Unmarshal([]byte(rr.Body.String()), &supers)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)

	// There are two supers registered, so it has to return 2
	assert.Equal(t, len(supers), 2)
}

func TestGetSuperByID(t *testing.T) {

	err := refreshTables()
	if err != nil {
		log.Fatal(err)
	}
	super, err := seedOneSuper()
	if err != nil {
		log.Fatal(err)
	}
	superSample := []struct {
		uuid         string
		statusCode   int
		name         string
		alignment    string
		errorMessage string
	}{
		{
			uuid:       strconv.Itoa(int(super.Uuid)),
			alignment:  super.Biography.Alignment,
			statusCode: 200,
			name:       super.Name,
		},
		{
			uuid:       "588",
			statusCode: 500,
		},
	}
	for _, v := range superSample {

		req, err := http.NewRequest("GET", "/supers", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"uuid": v.uuid})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetSuperById)
		handler.ServeHTTP(rr, req)

		newSuper := models.Super{}
		err = json.Unmarshal([]byte(rr.Body.String()), &newSuper)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, super.Name, newSuper.Name)
			assert.Equal(t, super.Biography.Alignment, newSuper.Biography.Alignment)
		}
	}
}

func TestGetSuperByName(t *testing.T) {

	err := refreshTables()
	if err != nil {
		log.Fatal(err)
	}
	super, err := seedOneSuper()
	if err != nil {
		log.Fatal(err)
	}
	superSample := []struct {
		uuid         string
		statusCode   int
		name         string
		alignment    string
		errorMessage string
	}{
		{
			uuid:       strconv.Itoa(int(super.Uuid)),
			alignment:  super.Biography.Alignment,
			statusCode: 200,
			name:       super.Name,
		},
		{
			name:       "banana",
			statusCode: 500,
		},
	}
	for _, v := range superSample {

		req, err := http.NewRequest("GET", "/supers/search", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"name": v.name})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetSuperByName)
		handler.ServeHTTP(rr, req)

		newSuper := models.Super{}
		err = json.Unmarshal([]byte(rr.Body.String()), &newSuper)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, super.Uuid, newSuper.Uuid)
			assert.Equal(t, super.Biography.Alignment, newSuper.Biography.Alignment)
		}
	}
}

func TestDeleteSuper(t *testing.T) {

	err := refreshTables()
	if err != nil {
		log.Fatal(err)
	}

	seedOneSuper()

	superSample := []struct {
		uuid         string
		statusCode   int
		errorMessage string
	}{
		{
			uuid:         "1",
			statusCode:   200,
			errorMessage: "",
		},
	}
	for _, v := range superSample {

		req, err := http.NewRequest("GET", "/supers", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"uuid": v.uuid})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteSuper)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
