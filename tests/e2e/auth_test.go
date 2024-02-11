package e2e

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    controller "organization_management/pkg/controllers"
)

func TestRegisterUserEndpoint(t *testing.T) {
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    reqBody := map[string]string{
        "name":     "John Doe",
        "email":    "john@example.com",
        "password": "password123",
    }
    requestBodyBytes, _ := json.Marshal(reqBody)

    req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBodyBytes))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    c.Request = req

    controller.RegisterUser()(c)

    assert.Equal(t, http.StatusCreated, w.Code)

    var responseBody map[string]interface{}
    err = json.Unmarshal(w.Body.Bytes(), &responseBody)
    if err != nil {
        t.Fatal(err)
    }
    assert.Equal(t, "success", responseBody["message"])
    assert.NotNil(t, responseBody["data"])
}


func TestEmailValidation(t *testing.T) {
    reqBody := map[string]string{
        "name":     "John Doe",
        "email":    "invalid_email", // Invalid email format
        "password": "password123",
    }
    requestBodyBytes, _ := json.Marshal(reqBody)
    req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBodyBytes))
    if err != nil {
        t.Fatal(err)
    }

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = req

    handler := controller.RegisterUser()

    handler(c)

    if status := w.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
    }

    var responseBody map[string]interface{}
    if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
        t.Fatal(err)
    }

    expectedErrorMessage := "Invalid email address"
    actualErrorMessage, ok := responseBody["error"].(string)
    if !ok {
        t.Errorf("error message not found in response body")
    }

    if actualErrorMessage != expectedErrorMessage {
        t.Errorf("handler returned unexpected error message: got %s want %s", actualErrorMessage, expectedErrorMessage)
    }
}

func TestRegisterUserWithMissingField(t *testing.T) {
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    reqBody := map[string]string{
        // Missing "name" and "email" fields
        "password": "password123",
    }
    requestBodyBytes, _ := json.Marshal(reqBody)

    req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBodyBytes))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    c.Request = req

    handler := controller.RegisterUser()

    handler(c)

    assert.Equal(t, http.StatusBadRequest, w.Code)

    var responseBody map[string]interface{}
    err = json.Unmarshal(w.Body.Bytes(), &responseBody)
    if err != nil {
        t.Fatal(err)
    }
    expectedErrorMessage := "Key: 'User.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'User.Email' Error:Field validation for 'Email' failed on the 'required' tag"
    actualErrorMessage, ok := responseBody["error"].(string)
    if !ok {
        t.Errorf("error message not found in response body")
    }
    assert.Equal(t, expectedErrorMessage, actualErrorMessage)
}