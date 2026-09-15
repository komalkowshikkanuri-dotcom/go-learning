package main

import (
    "net/http/httptest"
    "testing"
)

func TestGetUser(t *testing.T) {
	users = []User{
        	{ID: 1, Name: "John", Age: 25},
        	{ID: 2, Name: "Alice", Age: 30},
    	}
	handler := setupRoutes()

	req := httptest.NewRequest("GET", "/users/1", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != 200 {
    		t.Error("expected status 200")
	}
}
