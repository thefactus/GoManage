package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTasks(t *testing.T) {
	t.Run("expects to get empty data", func(t *testing.T) {
		clearTableTasksAndLists()

		payload := []byte(`{"title": "Tasks"}`)
		req, _ := http.NewRequest("POST", "/api/v1/lists", bytes.NewReader(payload))
		response := executeRequest(req)

		checkResponseCode(t, http.StatusCreated, response.Code)

		req, _ = http.NewRequest("GET", "/api/v1/lists/1/tasks", nil)
		response = executeRequest(req)

		checkResponseCode(t, http.StatusOK, response.Code)

		var actual Response
		json.Unmarshal(response.Body.Bytes(), &actual)

		expected := Response{
			Data:    []interface{}{},
			Message: "Tasks retrieved successfully",
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected '%v'. Got '%v'", expected, actual)
		}
	})

	t.Run("expects to get a task", func(t *testing.T) {
		clearTableTasksAndLists()

		payload := []byte(`{"title": "Tasks"}`)
		req, _ := http.NewRequest("POST", "/api/v1/lists", bytes.NewReader(payload))
		response := executeRequest(req)

		checkResponseCode(t, http.StatusCreated, response.Code)

		payload = []byte(`{"title": "Task 1"}`)
		req, _ = http.NewRequest("POST", "/api/v1/lists/1/tasks", bytes.NewReader(payload))
		response = executeRequest(req)

		checkResponseCode(t, http.StatusCreated, response.Code)

		req, _ = http.NewRequest("GET", "/api/v1/lists/1/tasks", nil)
		response = executeRequest(req)

		checkResponseCode(t, http.StatusOK, response.Code)

		var actual Response
		json.Unmarshal(response.Body.Bytes(), &actual)

		for _, item := range actual.Data {
			task, ok := item.(map[string]interface{})
			if !ok {
				t.Errorf("Expected task item to be a map. Got '%v'", item)
			}

			assert.Equal(t, float64(1), task["id"], "Expected id to be 1")
			assert.Equal(t, "Task 1", task["title"], "Expected title to be 'Task 1'")
		}
	})

	t.Run("expects to create a task", func(t *testing.T) {
		clearTableTasksAndLists()

		payload := []byte(`{"title": "Tasks"}`)
		req, _ := http.NewRequest("POST", "/api/v1/lists", bytes.NewReader(payload))
		response := executeRequest(req)

		checkResponseCode(t, http.StatusCreated, response.Code)

		payload = []byte(`{"title": "Task 1"}`)
		req, _ = http.NewRequest("POST", "/api/v1/lists/1/tasks", bytes.NewReader(payload))
		response = executeRequest(req)

		checkResponseCode(t, http.StatusCreated, response.Code)

		var result map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &result)
		if err != nil {
			t.Errorf("Error unmarshalling response: %v", err)
			return
		}

		// Check for message
		assert.Equal(t, "Task created successfully", result["message"],
			"Expected message to be 'Task created successfully'")
		assert.Equal(t, float64(1), result["data"].(map[string]interface{})["id"],
			"Expected id to be 1")
		assert.Equal(t, "Task 1", result["data"].(map[string]interface{})["title"],
			"Expected title to be 'Task 1'")
		assert.Equal(t, float64(1), result["data"].(map[string]interface{})["list_id"],
			"Expected list_id to be 1")
	})

	t.Run("expects to update a task", func(t *testing.T) {
		clearTableTasksAndLists()

		payload := []byte(`{"title": "Tasks"}`)
		req, _ := http.NewRequest("POST", "/api/v1/lists", bytes.NewReader(payload))
		response := executeRequest(req)

		checkResponseCode(t, http.StatusCreated, response.Code)

		payload = []byte(`{"title": "Task 1"}`)
		req, _ = http.NewRequest("POST", "/api/v1/lists/1/tasks", bytes.NewReader(payload))
		response = executeRequest(req)

		checkResponseCode(t, http.StatusCreated, response.Code)

		payload = []byte(`{"title": "Task 1 Updated"}`)
		req, _ = http.NewRequest("PUT", "/api/v1/lists/1/tasks/1", bytes.NewReader(payload))
		response = executeRequest(req)

		checkResponseCode(t, http.StatusOK, response.Code)

		var result map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &result)
		if err != nil {
			t.Errorf("Error unmarshalling response: %v", err)
			return
		}

		// Check for message
		assert.Equal(t, "Task updated successfully", result["message"],
			"Expected message to be 'Task updated successfully'")
		assert.Equal(t, float64(1), result["data"].(map[string]interface{})["id"],
			"Expected id to be 1")
		assert.Equal(t, "Task 1 Updated", result["data"].(map[string]interface{})["title"],
			"Expected title to be 'Task 1 Updated'")
		assert.Equal(t, float64(1), result["data"].(map[string]interface{})["list_id"],
			"Expected list_id to be 1")
	})

	t.Run("expects to mark task as done", func(t *testing.T) {
		clearTableTasksAndLists()

		payload := []byte(`{"title": "Tasks"}`)
		req, _ := http.NewRequest("POST", "/api/v1/lists", bytes.NewReader(payload))
		response := executeRequest(req)

		checkResponseCode(t, http.StatusCreated, response.Code)

		payload = []byte(`{"title": "Task 1"}`)
		req, _ = http.NewRequest("POST", "/api/v1/lists/1/tasks", bytes.NewReader(payload))
		response = executeRequest(req)

		checkResponseCode(t, http.StatusCreated, response.Code)

		req, _ = http.NewRequest("PATCH", "/api/v1/lists/1/tasks/1/done", nil)
		response = executeRequest(req)

		checkResponseCode(t, http.StatusOK, response.Code)

		var result map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &result)
		if err != nil {
			t.Errorf("Error unmarshalling response: %v", err)
			return
		}

		assert.Equal(t, "Task marked as done successfully", result["message"],
			"Expected message to be 'Task marked as done successfully'")
		assert.Equal(t, float64(1), result["data"].(map[string]interface{})["id"],
			"Expected id to be 1")
		assert.Equal(t, "Task 1", result["data"].(map[string]interface{})["title"],
			"Expected title to be 'Task 1'")
		assert.Equal(t, float64(1), result["data"].(map[string]interface{})["list_id"],
			"Expected list_id to be 1")
		assert.Equal(t, true, result["data"].(map[string]interface{})["done"],
			"Expected done to be true")
	})

	t.Run("expects to mark task as undone", func(t *testing.T) {
		clearTableTasksAndLists()

		payload := []byte(`{"title": "Tasks"}`)
		req, _ := http.NewRequest("POST", "/api/v1/lists", bytes.NewReader(payload))
		response := executeRequest(req)

		checkResponseCode(t, http.StatusCreated, response.Code)

		payload = []byte(`{"title": "Task 1"}`)
		req, _ = http.NewRequest("POST", "/api/v1/lists/1/tasks", bytes.NewReader(payload))
		response = executeRequest(req)

		checkResponseCode(t, http.StatusCreated, response.Code)

		req, _ = http.NewRequest("PATCH", "/api/v1/lists/1/tasks/1/done", nil)
		response = executeRequest(req)

		checkResponseCode(t, http.StatusOK, response.Code)

		req, _ = http.NewRequest("PATCH", "/api/v1/lists/1/tasks/1/undone", nil)
		response = executeRequest(req)

		checkResponseCode(t, http.StatusOK, response.Code)

		var result map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &result)
		if err != nil {
			t.Errorf("Error unmarshalling response: %v", err)
			return
		}

		assert.Equal(t, "Task marked as undone successfully", result["message"],
			"Expected message to be 'Task marked as undone successfully'")
		assert.Equal(t, float64(1), result["data"].(map[string]interface{})["id"],
			"Expected id to be 1")
		assert.Equal(t, "Task 1", result["data"].(map[string]interface{})["title"],
			"Expected title to be 'Task 1'")
		assert.Equal(t, float64(1), result["data"].(map[string]interface{})["list_id"],
			"Expected list_id to be 1")
		assert.Equal(t, false, result["data"].(map[string]interface{})["done"],
			"Expected done to be false")
	})

	t.Run("expects to delete a task", func(t *testing.T) {
		clearTableTasksAndLists()

		payload := []byte(`{"title": "Tasks"}`)
		req, _ := http.NewRequest("POST", "/api/v1/lists", bytes.NewReader(payload))
		response := executeRequest(req)

		checkResponseCode(t, http.StatusCreated, response.Code)

		payload = []byte(`{"title": "Task 1"}`)
		req, _ = http.NewRequest("POST", "/api/v1/lists/1/tasks", bytes.NewReader(payload))
		response = executeRequest(req)

		checkResponseCode(t, http.StatusCreated, response.Code)

		req, _ = http.NewRequest("DELETE", "/api/v1/lists/1/tasks/1", nil)
		response = executeRequest(req)

		checkResponseCode(t, http.StatusOK, response.Code)

		var result map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &result)
		if err != nil {
			t.Errorf("Error unmarshalling response: %v", err)
			return
		}

		// Check for message
		assert.Equal(t, "Task deleted successfully", result["message"],
			"Expected message to be 'Task deleted successfully'")
	})

	t.Run("expects to delete a list and all its tasks", func(t *testing.T) {
		clearTableTasksAndLists()

		payload := []byte(`{"title": "Tasks"}`)
		req, _ := http.NewRequest("POST", "/api/v1/lists", bytes.NewReader(payload))
		response := executeRequest(req)

		checkResponseCode(t, http.StatusCreated, response.Code)

		payload = []byte(`{"title": "Task 1"}`)
		req, _ = http.NewRequest("POST", "/api/v1/lists/1/tasks", bytes.NewReader(payload))
		response = executeRequest(req)

		checkResponseCode(t, http.StatusCreated, response.Code)

		req, _ = http.NewRequest("DELETE", "/api/v1/lists/1", nil)
		response = executeRequest(req)

		checkResponseCode(t, http.StatusOK, response.Code)

		var result map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &result)
		if err != nil {
			t.Errorf("Error unmarshalling response: %v", err)
			return
		}

		// Check for message
		assert.Equal(t, "List deleted successfully", result["message"],
			"Expected message to be 'List deleted successfully'")
	})

}
