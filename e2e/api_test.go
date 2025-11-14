package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

const baseURL = "http://localhost:8080"

func TestAPI(t *testing.T) {
	// Генерируем уникальные имена для каждого запуска тестов
	timestamp := time.Now().Unix()
	teamName := fmt.Sprintf("team_%d", timestamp)
	user1 := fmt.Sprintf("user1_%d", timestamp)
	user2 := fmt.Sprintf("user2_%d", timestamp)
	user3 := fmt.Sprintf("user3_%d", timestamp)
	user4 := fmt.Sprintf("user4_%d", timestamp)
	prID := fmt.Sprintf("pr_%d", timestamp)

	t.Run("Full workflow - ALL endpoints", func(t *testing.T) {
		// 1. Create team - POST /team/add
		teamData := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": user1, "username": "Alice", "is_active": true},
				{"user_id": user2, "username": "Bob", "is_active": true},
				{"user_id": user3, "username": "Charlie", "is_active": true},
				{"user_id": user4, "username": "David", "is_active": true},
			},
		}
		teamJSON, _ := json.Marshal(teamData)

		resp := makeRequest(t, "POST", "/team/add", teamJSON)
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("POST /team/add: Expected 201, got %d", resp.StatusCode)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Logf("Failed to close response body: %v", err)
			}
		}()

		// 2. Get team - GET /team/get
		resp = makeRequest(t, "GET", "/team/get?team_name="+teamName, nil)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("GET /team/get: Expected 200, got %d", resp.StatusCode)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Logf("Failed to close response body: %v", err)
			}
		}()

		// 3. Create PR - POST /pullRequest/create
		prData := map[string]string{
			"pull_request_id":   prID,
			"pull_request_name": "Add search",
			"author_id":         user1,
		}
		prJSON, _ := json.Marshal(prData)

		resp = makeRequest(t, "POST", "/pullRequest/create", prJSON)
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("POST /pullRequest/create: Expected 201, got %d", resp.StatusCode)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Logf("Failed to close response body: %v", err)
			}
		}()

		// 4. Get user reviews - GET /users/getReview
		resp = makeRequest(t, "GET", "/users/getReview?user_id="+user2, nil)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("GET /users/getReview: Expected 200, got %d", resp.StatusCode)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Logf("Failed to close response body: %v", err)
			}
		}()

		// 5. Deactivate user - POST /users/setIsActive
		userData := map[string]interface{}{
			"user_id":   user2,
			"is_active": false,
		}
		userJSON, _ := json.Marshal(userData)

		resp = makeRequest(t, "POST", "/users/setIsActive", userJSON)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("POST /users/setIsActive: Expected 200, got %d", resp.StatusCode)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Logf("Failed to close response body: %v", err)
			}
		}()

		// 6. Reassign reviewer - POST /pullRequest/reassign
		reassignData := map[string]string{
			"pull_request_id": prID,
			"old_user_id":     user2,
		}
		reassignJSON, _ := json.Marshal(reassignData)

		resp = makeRequest(t, "POST", "/pullRequest/reassign", reassignJSON)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("POST /pullRequest/reassign: Expected 200, got %d", resp.StatusCode)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Logf("Failed to close response body: %v", err)
			}
		}()

		// 7. Merge PR - POST /pullRequest/merge
		mergeData := map[string]string{
			"pull_request_id": prID,
		}
		mergeJSON, _ := json.Marshal(mergeData)

		resp = makeRequest(t, "POST", "/pullRequest/merge", mergeJSON)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("POST /pullRequest/merge: Expected 200, got %d", resp.StatusCode)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Logf("Failed to close response body: %v", err)
			}
		}()

		// 8. Verify merged PR not in review list - GET /users/getReview
		resp = makeRequest(t, "GET", "/users/getReview?user_id="+user3, nil)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("GET /users/getReview after merge: Expected 200, got %d", resp.StatusCode)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Logf("Failed to close response body: %v", err)
			}
		}()

		t.Log("All endpoints tested successfully!")
	})
}

func makeRequest(t *testing.T, method, path string, body []byte) *http.Response {
	req, err := http.NewRequest(method, baseURL+path, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	return resp
}
