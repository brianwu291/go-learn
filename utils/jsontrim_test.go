package utils

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// TestTrimJSON tests the main TrimJSON function
func TestTrimJSON(t *testing.T) {
	// Test case 1: Trim specific paths in a simple object
	input := `{
		"name": "  John Doe  ",
		"email": " john@example.com ",
		"address": "  123 Main St  ",
		"age": 30
	}`

	expected := `{"name":"John Doe","email":"john@example.com","address":"  123 Main St  ","age":30}`

	result, modified, err := TrimJSON(input, TrimOptions{
		Paths: []string{"name", "email"},
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if !modified {
		t.Error("Expected modified=true, got false")
	}

	// Normalize JSON for comparison
	var expectedObj, resultObj interface{}
	json.Unmarshal([]byte(expected), &expectedObj)
	json.Unmarshal([]byte(result), &resultObj)

	expectedJSON, _ := json.Marshal(expectedObj)
	resultJSON, _ := json.Marshal(resultObj)

	if string(expectedJSON) != string(resultJSON) {
		t.Errorf("Expected: %s, got: %s", expectedJSON, resultJSON)
	}

	// Test case 2: Test with nested paths
	input = `{
		"user": {
			"name": " Jane Smith ",
			"contact": {
				"email": " jane@example.com ",
				"phone": " 555-1234 "
			}
		},
		"tags": [" tag1 ", "tag2", " tag3 "],
		"items": [
			{"id": 1, "name": " Item 1 "},
			{"id": 2, "name": " Item 2 "}
		]
	}`

	expected = `{"user":{"name":"Jane Smith","contact":{"email":"jane@example.com","phone":" 555-1234 "}},"tags":[" tag1 ","tag2"," tag3 "],"items":[{"id":1,"name":" Item 1 "},{"id":2,"name":" Item 2 "}]}`

	result, modified, err = TrimJSON(input, TrimOptions{
		Paths: []string{"user.name", "user.contact.email"},
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if !modified {
		t.Error("Expected modified=true, got false")
	}

	// Normalize JSON for comparison
	json.Unmarshal([]byte(expected), &expectedObj)
	json.Unmarshal([]byte(result), &resultObj)

	expectedJSON, _ = json.Marshal(expectedObj)
	resultJSON, _ = json.Marshal(resultObj)

	if string(expectedJSON) != string(resultJSON) {
		t.Errorf("Expected: %s, got: %s", expectedJSON, resultJSON)
	}

	// Test case 3: Test with array indices
	input = `{
		"users": [
			{"name": " User 1 "},
			{"name": " User 2 "},
			{"name": " User 3 "}
		]
	}`

	expected = `{"users":[{"name":"User 1"},{"name":" User 2 "},{"name":" User 3 "}]}`

	result, modified, err = TrimJSON(input, TrimOptions{
		Paths: []string{"users.0.name"},
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if !modified {
		t.Error("Expected modified=true, got false")
	}

	// Normalize JSON for comparison
	json.Unmarshal([]byte(expected), &expectedObj)
	json.Unmarshal([]byte(result), &resultObj)

	expectedJSON, _ = json.Marshal(expectedObj)
	resultJSON, _ = json.Marshal(resultObj)

	if string(expectedJSON) != string(resultJSON) {
		t.Errorf("Expected: %s, got: %s", expectedJSON, resultJSON)
	}

	// Test case 4: Path doesn't exist (should not modify)
	input = `{"name": " John Doe "}`

	_, modified, err = TrimJSON(input, TrimOptions{
		Paths: []string{"email"},
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if modified {
		t.Error("Expected modified=false, got true")
	}

	// Test case 5: Path exists but not a string
	input = `{"name": " John Doe ", "age": 30}`

	_, modified, err = TrimJSON(input, TrimOptions{
		Paths: []string{"age"},
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if modified {
		t.Error("Expected modified=false, got true")
	}

	// Test case 6: No paths provided (should not modify)
	input = `{"name": " John Doe "}`

	_, modified, err = TrimJSON(input, TrimOptions{
		Paths: []string{},
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if modified {
		t.Error("Expected modified=false, got true")
	}

	// Test case 7: Invalid JSON
	input = `{invalid json}`

	_, _, err = TrimJSON(input, TrimOptions{
		Paths: []string{"name"},
	})

	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

// TestTrimJSONBytes tests the byte slice version
func TestTrimJSONBytes(t *testing.T) {
	input := []byte(`{"name": " John Doe ", "age": 30}`)
	expected := []byte(`{"name":"John Doe","age":30}`)

	result, modified, err := TrimJSONBytes(input, TrimOptions{
		Paths: []string{"name"},
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if !modified {
		t.Error("Expected modified=true, got false")
	}

	// Normalize JSON for comparison
	var expectedObj, resultObj interface{}
	json.Unmarshal(expected, &expectedObj)
	json.Unmarshal(result, &resultObj)

	expectedJSON, _ := json.Marshal(expectedObj)
	resultJSON, _ := json.Marshal(resultObj)

	if string(expectedJSON) != string(resultJSON) {
		t.Errorf("Expected: %s, got: %s", expectedJSON, resultJSON)
	}

	// Test with no modifications
	input = []byte(`{"name":"John","age":30}`)

	result, modified, err = TrimJSONBytes(input, TrimOptions{
		Paths: []string{"name"},
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if modified {
		t.Error("Expected modified=false, got true")
	}

	// Original bytes should be returned
	if !bytes.Equal(input, result) {
		t.Error("Expected original bytes to be returned")
	}
}

// TestUnmarshalAndTrim tests unmarshaling and trimming in one step
func TestUnmarshalAndTrim(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	input := []byte(`{"name": " John Doe ", "email": " john@example.com ", "age": 30}`)

	var user User
	err := UnmarshalAndTrim(input, &user, TrimOptions{
		Paths: []string{"name", "email"},
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if user.Name != "John Doe" {
		t.Errorf("Expected Name='John Doe', got '%s'", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Errorf("Expected Email='john@example.com', got '%s'", user.Email)
	}
	if user.Age != 30 {
		t.Errorf("Expected Age=30, got %d", user.Age)
	}

	// Test with one specific path
	input = []byte(`{"name": " John Doe ", "email": " john@example.com ", "age": 30}`)

	var user2 User
	err = UnmarshalAndTrim(input, &user2, TrimOptions{
		Paths: []string{"name"},
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if user2.Name != "John Doe" {
		t.Errorf("Expected Name='John Doe', got '%s'", user2.Name)
	}
	if user2.Email != " john@example.com " {
		t.Errorf("Expected Email=' john@example.com ', got '%s'", user2.Email)
	}
}

// TestReadAndTrim tests reading from an io.Reader and trimming
func TestReadAndTrim(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	input := `{"name": " John Doe ", "email": " john@example.com "}`
	reader := strings.NewReader(input)

	var user User
	err := ReadAndTrim(reader, &user, TrimOptions{
		Paths: []string{"name", "email"},
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if user.Name != "John Doe" {
		t.Errorf("Expected Name='John Doe', got '%s'", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Errorf("Expected Email='john@example.com', got '%s'", user.Email)
	}
}

// TestAdvancedPathSyntax tests more complex path expressions supported by gjson
func TestAdvancedPathSyntax(t *testing.T) {
	input := `{
		"users": [
			{
				"name": " User 1 ",
				"emails": [" email1@example.com ", " email2@example.com "]
			},
			{
				"name": " User 2 ",
				"emails": [" user2@example.com "]
			}
		],
		"metadata": {
			"tags": {
				"important": " urgent ",
				"category": " test "
			}
		}
	}`

	// Test array wildcard
	result, modified, err := TrimJSON(input, TrimOptions{
		Paths: []string{"users.#.name", "metadata.tags.important"},
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if !modified {
		t.Error("Expected modified=true, got false")
	}

	// Verify specific paths were trimmed
	parsedResult := make(map[string]interface{})
	json.Unmarshal([]byte(result), &parsedResult)

	users := parsedResult["users"].([]interface{})
	for _, u := range users {
		user := u.(map[string]interface{})
		name := user["name"].(string)
		if strings.HasPrefix(name, " ") || strings.HasSuffix(name, " ") {
			t.Errorf("Expected trimmed name, got '%s'", name)
		}
	}

	metadata := parsedResult["metadata"].(map[string]interface{})
	tags := metadata["tags"].(map[string]interface{})
	important := tags["important"].(string)
	if important != "urgent" {
		t.Errorf("Expected 'urgent', got '%s'", important)
	}

	category := tags["category"].(string)
	if category != " test " {
		t.Errorf("Expected ' test ', got '%s'", category)
	}
}
