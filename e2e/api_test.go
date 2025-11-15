package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func getBaseURL() string {
	url := os.Getenv("TEST_URL")
	if url == "" {
		log.Fatal("TEST_URL environment variable is required")
	}
	return url
}

func TestAPI(t *testing.T) {
	// Генерируем уникальные имена для каждого запуска тестов
	timestamp := time.Now().Unix()
	teamName := fmt.Sprintf("team_%d", timestamp)
	user1 := fmt.Sprintf("user1_%d", timestamp)
	user2 := fmt.Sprintf("user2_%d", timestamp)
	user3 := fmt.Sprintf("user3_%d", timestamp)
	user4 := fmt.Sprintf("user4_%d", timestamp)
	prID := fmt.Sprintf("pr_%d", timestamp)

	t.Run("Starting API tests", func(t *testing.T) {
		testFullWorkflow(t, teamName, user1, user2, user3, user4, prID)
	})

	// Тесты для проверки ошибок
	t.Run("Error cases", func(t *testing.T) {
		testErrorCases(t, timestamp)
	})
}

func testFullWorkflow(t *testing.T, teamName, user1, user2, user3, user4, prID string) {
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

	var createTeamResponse map[string]interface{}
	parseAndCheckResponse(t, resp, &createTeamResponse)

	// Проверяем структуру ответа
	team, ok := createTeamResponse["team"].(map[string]interface{})
	if !ok {
		t.Fatal("POST /team/add: response should contain 'team' object")
	}
	if team["team_name"] != teamName {
		t.Errorf("Expected team_name %s, got %s", teamName, team["team_name"])
	}

	members, ok := team["members"].([]interface{})
	if !ok || len(members) != 4 {
		t.Errorf("Expected 4 members, got %d", len(members))
	}
	closeBody(t, resp)

	// 2. Get team - GET /team/get
	resp = makeRequest(t, "GET", "/team/get?team_name="+teamName, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /team/get: Expected 200, got %d", resp.StatusCode)
	}

	var getTeamResponse map[string]interface{}
	parseAndCheckResponse(t, resp, &getTeamResponse)

	// Проверяем структуру ответа
	if getTeamResponse["team_name"] != teamName {
		t.Errorf("GET /team/get: Expected team_name %s, got %s", teamName, getTeamResponse["team_name"])
	}

	members, ok = getTeamResponse["members"].([]interface{})
	if !ok {
		t.Fatal("GET /team/get: members field is missing or invalid")
	}

	if len(members) != 4 {
		t.Errorf("GET /team/get: Expected 4 members, got %d", len(members))
	}

	// Проверяем структуру каждого члена команды
	for i, member := range members {
		memberMap := member.(map[string]interface{})
		requiredFields := []string{"user_id", "username", "is_active"}
		for _, field := range requiredFields {
			if _, exists := memberMap[field]; !exists {
				t.Errorf("Member %d missing required field: %s", i, field)
			}
		}
		if !memberMap["is_active"].(bool) {
			t.Errorf("All members should be active initially")
		}
	}
	closeBody(t, resp)

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

	var createPRResponse map[string]interface{}
	parseAndCheckResponse(t, resp, &createPRResponse)

	// Проверяем структуру PR
	pr, ok := createPRResponse["pr"].(map[string]interface{})
	if !ok {
		t.Fatal("POST /pullRequest/create: response should contain 'pr' object")
	}

	// Проверяем обязательные поля PR
	if pr["pull_request_id"] != prID {
		t.Errorf("Expected pull_request_id %s, got %s", prID, pr["pull_request_id"])
	}
	if pr["pull_request_name"] != "Add search" {
		t.Errorf("Expected pull_request_name 'Add search', got %s", pr["pull_request_name"])
	}
	if pr["author_id"] != user1 {
		t.Errorf("Expected author_id %s, got %s", user1, pr["author_id"])
	}
	if pr["status"] != "OPEN" {
		t.Errorf("Expected status 'OPEN', got %s", pr["status"])
	}

	// Проверяем назначенных ревьюверов
	reviewers, ok := pr["assigned_reviewers"].([]interface{})
	if !ok {
		t.Fatal("assigned_reviewers field is missing or invalid")
	}
	if len(reviewers) > 2 {
		t.Errorf("Expected 0-2 assigned reviewers, got %d", len(reviewers))
	}

	// Проверяем что автор не среди ревьюверов
	for _, reviewer := range reviewers {
		if reviewer == user1 {
			t.Error("Author should not be assigned as reviewer")
		}
	}

	// Сохраняем первоначальных ревьюверов для последующих проверок
	initialReviewers := make([]string, len(reviewers))
	for i, r := range reviewers {
		initialReviewers[i] = r.(string)
	}
	closeBody(t, resp)

	// 4. Get user reviews - GET /users/getReview (для одного из ревьюверов)
	reviewerToCheck := initialReviewers[0]
	resp = makeRequest(t, "GET", "/users/getReview?user_id="+reviewerToCheck, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /users/getReview: Expected 200, got %d", resp.StatusCode)
	}

	var userReviewsResponse map[string]interface{}
	parseAndCheckResponse(t, resp, &userReviewsResponse)

	// Проверяем структуру ответа
	if userReviewsResponse["user_id"] != reviewerToCheck {
		t.Errorf("Expected user_id %s, got %s", reviewerToCheck, userReviewsResponse["user_id"])
	}

	pullRequests, ok := userReviewsResponse["pull_requests"].([]interface{})
	if !ok {
		t.Fatal("pull_requests field is missing or invalid")
	}

	// Проверяем что PR есть в списке ревью
	foundPR := false
	for _, prShort := range pullRequests {
		prMap := prShort.(map[string]interface{})
		if prMap["pull_request_id"] == prID {
			foundPR = true
			// Проверяем структуру PullRequestShort
			requiredFields := []string{"pull_request_id", "pull_request_name", "author_id", "status"}
			for _, field := range requiredFields {
				if _, exists := prMap[field]; !exists {
					t.Errorf("PullRequestShort missing required field: %s", field)
				}
			}
			if prMap["status"] != "OPEN" {
				t.Errorf("PR status should be OPEN, got %s", prMap["status"])
			}
			break
		}
	}

	if !foundPR {
		t.Error("PR should be in user's review list")
	}
	closeBody(t, resp)

	// 5. Deactivate user - POST /users/setIsActive
	userToDeactivate := initialReviewers[0]
	userData := map[string]interface{}{
		"user_id":   userToDeactivate,
		"is_active": false,
	}
	userJSON, _ := json.Marshal(userData)

	resp = makeRequest(t, "POST", "/users/setIsActive", userJSON)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("POST /users/setIsActive: Expected 200, got %d", resp.StatusCode)
	}

	var deactivateResponse map[string]interface{}
	parseAndCheckResponse(t, resp, &deactivateResponse)

	// Проверяем структуру ответа
	user, ok := deactivateResponse["user"].(map[string]interface{})
	if !ok {
		t.Fatal("POST /users/setIsActive: response should contain 'user' object")
	}

	if user["user_id"] != userToDeactivate {
		t.Errorf("Expected user_id %s, got %s", userToDeactivate, user["user_id"])
	}
	if user["is_active"] != false {
		t.Error("User should be deactivated")
	}
	if user["team_name"] != teamName {
		t.Errorf("Expected team_name %s, got %s", teamName, user["team_name"])
	}
	closeBody(t, resp)

	// 6. Reassign reviewer - POST /pullRequest/reassign
	reassignData := map[string]string{
		"pull_request_id": prID,
		"old_user_id":     userToDeactivate,
	}
	reassignJSON, _ := json.Marshal(reassignData)

	resp = makeRequest(t, "POST", "/pullRequest/reassign", reassignJSON)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("POST /pullRequest/reassign: Expected 200, got %d", resp.StatusCode)
	}

	var reassignResponse map[string]interface{}
	parseAndCheckResponse(t, resp, &reassignResponse)

	// Проверяем структуру ответа
	if reassignResponse["replaced_by"] == "" {
		t.Error("replaced_by field should contain new reviewer ID")
	}
	newReviewer := reassignResponse["replaced_by"].(string)

	reassignedPR, ok := reassignResponse["pr"].(map[string]interface{})
	if !ok {
		t.Fatal("POST /pullRequest/reassign: response should contain 'pr' object")
	}

	// Проверяем что старый ревьювер удален, а новый добавлен
	updatedReviewers, ok := reassignedPR["assigned_reviewers"].([]interface{})
	if !ok {
		t.Fatal("assigned_reviewers field is missing after reassign")
	}

	// Проверяем что деактивированный пользователь удален из ревьюверов
	for _, reviewer := range updatedReviewers {
		if reviewer == userToDeactivate {
			t.Error("Deactivated user should be removed from reviewers")
		}
	}

	// Проверяем что новый ревьювер добавлен
	foundNewReviewer := false
	for _, reviewer := range updatedReviewers {
		if reviewer == newReviewer {
			foundNewReviewer = true
			break
		}
	}
	if !foundNewReviewer {
		t.Error("New reviewer should be in assigned_reviewers")
	}
	closeBody(t, resp)

	// 7. Проверяем что у нового ревьювера появился PR в списке
	resp = makeRequest(t, "GET", "/users/getReview?user_id="+newReviewer, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /users/getReview for new reviewer: Expected 200, got %d", resp.StatusCode)
	}

	var newReviewerResponse map[string]interface{}
	parseAndCheckResponse(t, resp, &newReviewerResponse)

	newReviewerPRs, ok := newReviewerResponse["pull_requests"].([]interface{})
	if !ok {
		t.Fatal("pull_requests field missing for new reviewer")
	}

	foundPRInNewReviewer := false
	for _, prShort := range newReviewerPRs {
		prMap := prShort.(map[string]interface{})
		if prMap["pull_request_id"] == prID {
			foundPRInNewReviewer = true
			break
		}
	}
	if !foundPRInNewReviewer {
		t.Error("PR should be assigned to new reviewer")
	}
	closeBody(t, resp)

	// 8. Merge PR - POST /pullRequest/merge
	mergeData := map[string]string{
		"pull_request_id": prID,
	}
	mergeJSON, _ := json.Marshal(mergeData)

	resp = makeRequest(t, "POST", "/pullRequest/merge", mergeJSON)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("POST /pullRequest/merge: Expected 200, got %d", resp.StatusCode)
	}

	var mergeResponse map[string]interface{}
	parseAndCheckResponse(t, resp, &mergeResponse)

	// Проверяем структуру ответа
	mergedPR, ok := mergeResponse["pr"].(map[string]interface{})
	if !ok {
		t.Fatal("POST /pullRequest/merge: response should contain 'pr' object")
	}

	if mergedPR["pull_request_id"] != prID {
		t.Errorf("Expected pull_request_id %s, got %s", prID, mergedPR["pull_request_id"])
	}
	if mergedPR["status"] != "MERGED" {
		t.Errorf("Expected status 'MERGED', got %s", mergedPR["status"])
	}
	if mergedPR["mergedAt"] == nil {
		t.Error("mergedAt field should be set after merge")
	}
	closeBody(t, resp)

	// 9. Verify merged PR not in review list - GET /users/getReview
	resp = makeRequest(t, "GET", "/users/getReview?user_id="+newReviewer, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /users/getReview after merge: Expected 200, got %d", resp.StatusCode)
	}

	var finalReviewsResponse map[string]interface{}
	parseAndCheckResponse(t, resp, &finalReviewsResponse)

	finalPRs, ok := finalReviewsResponse["pull_requests"].([]interface{})
	if !ok {
		t.Fatal("pull_requests field missing in final check")
	}

	// Проверяем что мерженый PR не в списке на ревью
	for _, prShort := range finalPRs {
		prMap := prShort.(map[string]interface{})
		if prMap["pull_request_id"] == prID {
			t.Error("Merged PR should not be in review list")
		}
	}
	closeBody(t, resp)

	t.Log("All endpoints teste successfully")
}

func testErrorCases(t *testing.T, timestamp int64) {
	errorTimestamp := time.Now().UnixNano()
	errorTeamName := fmt.Sprintf("error_team_%d", errorTimestamp)
	errorUser1 := fmt.Sprintf("error_user1_%d", errorTimestamp)
	errorUser2 := fmt.Sprintf("error_user2_%d", errorTimestamp)
	errorUser3 := fmt.Sprintf("error_user3_%d", errorTimestamp)
	errorPR := fmt.Sprintf("error_pr_%d", errorTimestamp)

	teamData := map[string]interface{}{
		"team_name": errorTeamName,
		"members": []map[string]interface{}{
			{"user_id": errorUser1, "username": "ErrorUser1", "is_active": true},
			{"user_id": errorUser2, "username": "ErrorUser2", "is_active": true},
			{"user_id": errorUser3, "username": "ErrorUser3", "is_active": true},
		},
	}
	teamJSON, _ := json.Marshal(teamData)

	resp := makeRequest(t, "POST", "/team/add", teamJSON)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create team for error tests: Expected 201, got %d", resp.StatusCode)
	}
	closeBody(t, resp)

	prData := map[string]string{
		"pull_request_id":   errorPR,
		"pull_request_name": "Error Test PR",
		"author_id":         errorUser1,
	}
	prJSON, _ := json.Marshal(prData)

	resp = makeRequest(t, "POST", "/pullRequest/create", prJSON)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create PR for error tests: Expected 201, got %d", resp.StatusCode)
	}
	closeBody(t, resp)

	// 1. Попытка создать команду с существующим именем - TEAM_EXISTS
	existingTeamData := map[string]interface{}{
		"team_name": errorTeamName,
		"members": []map[string]interface{}{
			{"user_id": "new_user", "username": "NewUser", "is_active": true},
		},
	}
	existingTeamJSON, _ := json.Marshal(existingTeamData)

	resp = makeRequest(t, "POST", "/team/add", existingTeamJSON)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("POST /team/add with existing name: Expected 400, got %d", resp.StatusCode)
	} else {
		var errorResponse map[string]interface{}
		parseAndCheckResponse(t, resp, &errorResponse)
		checkErrorCode(t, errorResponse, "TEAM_EXISTS")
		closeBody(t, resp)
	}

	// 2. Попытка получить несуществующую команду - NOT_FOUND
	resp = makeRequest(t, "GET", "/team/get?team_name=nonexistent_team_123", nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("GET /team/get nonexistent: Expected 404, got %d", resp.StatusCode)
	} else {
		var errorResponse map[string]interface{}
		parseAndCheckResponse(t, resp, &errorResponse)
		checkErrorCode(t, errorResponse, "NOT_FOUND")
		closeBody(t, resp)
	}

	// 3. Попытка создать PR с существующим ID - PR_EXISTS
	existingPRData := map[string]string{
		"pull_request_id":   errorPR,
		"pull_request_name": "Duplicate PR",
		"author_id":         errorUser1,
	}
	existingPRJSON, _ := json.Marshal(existingPRData)

	resp = makeRequest(t, "POST", "/pullRequest/create", existingPRJSON)
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("POST /pullRequest/create with existing ID: Expected 409, got %d", resp.StatusCode)
	} else {
		var errorResponse map[string]interface{}
		parseAndCheckResponse(t, resp, &errorResponse)
		checkErrorCode(t, errorResponse, "PR_EXISTS")
		closeBody(t, resp)
	}

	// 4. Попытка создать PR с несуществующим автором - NOT_FOUND
	nonexistentAuthorData := map[string]string{
		"pull_request_id":   "new_pr_123",
		"pull_request_name": "PR with bad author",
		"author_id":         "nonexistent_user_123",
	}
	nonexistentAuthorJSON, _ := json.Marshal(nonexistentAuthorData)

	resp = makeRequest(t, "POST", "/pullRequest/create", nonexistentAuthorJSON)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("POST /pullRequest/create with nonexistent author: Expected 404, got %d", resp.StatusCode)
	} else {
		var errorResponse map[string]interface{}
		parseAndCheckResponse(t, resp, &errorResponse)
		checkErrorCode(t, errorResponse, "NOT_FOUND")
		closeBody(t, resp)
	}

	// 5. Попытка мержить несуществующий PR - NOT_FOUND
	nonexistentMergeData := map[string]string{
		"pull_request_id": "nonexistent_pr_123",
	}
	nonexistentMergeJSON, _ := json.Marshal(nonexistentMergeData)

	resp = makeRequest(t, "POST", "/pullRequest/merge", nonexistentMergeJSON)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("POST /pullRequest/merge with nonexistent PR: Expected 404, got %d", resp.StatusCode)
	} else {
		var errorResponse map[string]interface{}
		parseAndCheckResponse(t, resp, &errorResponse)
		checkErrorCode(t, errorResponse, "NOT_FOUND")
		closeBody(t, resp)
	}

	// 6. Попытка переназначить несуществующий PR - NOT_FOUND
	nonexistentReassignData := map[string]string{
		"pull_request_id": "nonexistent_pr_123",
		"old_user_id":     errorUser1,
	}
	nonexistentReassignJSON, _ := json.Marshal(nonexistentReassignData)

	resp = makeRequest(t, "POST", "/pullRequest/reassign", nonexistentReassignJSON)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("POST /pullRequest/reassign with nonexistent PR: Expected 404, got %d", resp.StatusCode)
	} else {
		var errorResponse map[string]interface{}
		parseAndCheckResponse(t, resp, &errorResponse)
		checkErrorCode(t, errorResponse, "NOT_FOUND")
		closeBody(t, resp)
	}

	// 7. Попытка переназначить несуществующего пользователя - NOT_FOUND
	nonexistentUserReassignData := map[string]string{
		"pull_request_id": errorPR,
		"old_user_id":     "nonexistent_user_123",
	}
	nonexistentUserReassignJSON, _ := json.Marshal(nonexistentUserReassignData)

	resp = makeRequest(t, "POST", "/pullRequest/reassign", nonexistentUserReassignJSON)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("POST /pullRequest/reassign with nonexistent user: Expected 404, got %d", resp.StatusCode)
	} else {
		var errorResponse map[string]interface{}
		parseAndCheckResponse(t, resp, &errorResponse)
		checkErrorCode(t, errorResponse, "NOT_FOUND")
		closeBody(t, resp)
	}

	// 8. Попытка переназначить пользователя, который не назначен ревьювером - NOT_ASSIGNED
	specialTeam := fmt.Sprintf("special_team_%d", errorTimestamp)
	specialUser1 := fmt.Sprintf("special_user1_%d", errorTimestamp) // автор
	specialUser2 := fmt.Sprintf("special_user2_%d", errorTimestamp)
	specialUser3 := fmt.Sprintf("special_user3_%d", errorTimestamp)
	specialUser4 := fmt.Sprintf("special_user4_%d", errorTimestamp)

	specialTeamData := map[string]interface{}{
		"team_name": specialTeam,
		"members": []map[string]interface{}{
			{"user_id": specialUser1, "username": "Special1", "is_active": true},
			{"user_id": specialUser2, "username": "Special2", "is_active": true},
			{"user_id": specialUser3, "username": "Special3", "is_active": true},
			{"user_id": specialUser4, "username": "Special4", "is_active": true},
		},
	}
	specialTeamJSON, _ := json.Marshal(specialTeamData)

	resp = makeRequest(t, "POST", "/team/add", specialTeamJSON)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create special team: Expected 201, got %d", resp.StatusCode)
	}
	closeBody(t, resp)

	// Создаем PR - назначатся 2 ревьювера из 3 доступных
	specialPR := fmt.Sprintf("special_pr_%d", errorTimestamp)
	specialPRData := map[string]string{
		"pull_request_id":   specialPR,
		"pull_request_name": "Special PR",
		"author_id":         specialUser1,
	}
	specialPRJSON, _ := json.Marshal(specialPRData)

	resp = makeRequest(t, "POST", "/pullRequest/create", specialPRJSON)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create special PR: Expected 201, got %d", resp.StatusCode)
	}

	var specialPRResponse map[string]interface{}
	parseAndCheckResponse(t, resp, &specialPRResponse)
	specialPRObj := specialPRResponse["pr"].(map[string]interface{})
	assignedReviewers := specialPRObj["assigned_reviewers"].([]interface{})
	closeBody(t, resp)

	// Находим пользователя который не назначен ревьювером
	allUsers := []string{specialUser2, specialUser3, specialUser4}
	var notAssignedUser string

	for _, user := range allUsers {
		isAssigned := false
		for _, reviewer := range assignedReviewers {
			if reviewer == user {
				isAssigned = true
				break
			}
		}
		if !isAssigned {
			notAssignedUser = user
			break
		}
	}

	if notAssignedUser == "" {
		t.Fatal("Should have at least one not assigned user in team of 4")
	}

	specialNotAssignedData := map[string]string{
		"pull_request_id": specialPR,
		"old_user_id":     notAssignedUser,
	}
	specialNotAssignedJSON, _ := json.Marshal(specialNotAssignedData)

	resp = makeRequest(t, "POST", "/pullRequest/reassign", specialNotAssignedJSON)
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("POST /pullRequest/reassign with not assigned user: Expected 409, got %d", resp.StatusCode)
	} else {
		var errorResponse map[string]interface{}
		parseAndCheckResponse(t, resp, &errorResponse)
		checkErrorCode(t, errorResponse, "NOT_ASSIGNED")
		closeBody(t, resp)
	}

	// 9. Попытка переназначить для мерженого PR - PR_MERGED
	mergedTestPR := fmt.Sprintf("merged_test_pr_%d", errorTimestamp)
	mergedPRData := map[string]string{
		"pull_request_id":   mergedTestPR,
		"pull_request_name": "PR for merge test",
		"author_id":         errorUser1,
	}
	mergedPRJSON, _ := json.Marshal(mergedPRData)

	resp = makeRequest(t, "POST", "/pullRequest/create", mergedPRJSON)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create PR for merge test: Expected 201, got %d", resp.StatusCode)
	}

	// Получаем назначенных ревьюверов
	var createdPRResponse map[string]interface{}
	parseAndCheckResponse(t, resp, &createdPRResponse)
	prObj := createdPRResponse["pr"].(map[string]interface{})
	reviewers := prObj["assigned_reviewers"].([]interface{})
	if len(reviewers) == 0 {
		t.Fatal("Expected at least one reviewer for merge test")
	}
	reviewerToReplace := reviewers[0].(string)
	closeBody(t, resp)

	// Мержим PR
	mergeTestData := map[string]string{
		"pull_request_id": mergedTestPR,
	}
	mergeTestJSON, _ := json.Marshal(mergeTestData)

	resp = makeRequest(t, "POST", "/pullRequest/merge", mergeTestJSON)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to merge PR for test: Expected 200, got %d", resp.StatusCode)
	}
	closeBody(t, resp)

	// Пытаемся переназначить ревьювера для мерженого PR - PR_MERGED
	mergedReassignData := map[string]string{
		"pull_request_id": mergedTestPR,
		"old_user_id":     reviewerToReplace,
	}
	mergedReassignJSON, _ := json.Marshal(mergedReassignData)

	resp = makeRequest(t, "POST", "/pullRequest/reassign", mergedReassignJSON)
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("POST /pullRequest/reassign on merged PR: Expected 409, got %d", resp.StatusCode)
	} else {
		var errorResponse map[string]interface{}
		parseAndCheckResponse(t, resp, &errorResponse)
		checkErrorCode(t, errorResponse, "PR_MERGED")
		closeBody(t, resp)
	}

	// 10. Попытка переназначить когда нет кандидатов - NO_CANDIDATE
	timestamp2 := time.Now().UnixNano()
	isolatedTeam := fmt.Sprintf("isolated_team_%d", timestamp2)
	userA := fmt.Sprintf("userA_%d", timestamp2)
	userB := fmt.Sprintf("userB_%d", timestamp2)
	userC := fmt.Sprintf("userC_%d", timestamp2)

	isolatedTeamData := map[string]interface{}{
		"team_name": isolatedTeam,
		"members": []map[string]interface{}{
			{"user_id": userA, "username": "UserA", "is_active": true},
			{"user_id": userB, "username": "UserB", "is_active": true},
			{"user_id": userC, "username": "UserC", "is_active": true},
		},
	}
	isolatedTeamJSON, _ := json.Marshal(isolatedTeamData)

	resp = makeRequest(t, "POST", "/team/add", isolatedTeamJSON)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create isolated team: Expected 201, got %d", resp.StatusCode)
	}
	closeBody(t, resp)

	isolatedPR := fmt.Sprintf("isolated_pr_%d", timestamp2)
	isolatedPRData := map[string]string{
		"pull_request_id":   isolatedPR,
		"pull_request_name": "Isolated PR",
		"author_id":         userA,
	}
	isolatedPRJSON, _ := json.Marshal(isolatedPRData)

	resp = makeRequest(t, "POST", "/pullRequest/create", isolatedPRJSON)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create isolated PR: Expected 201, got %d", resp.StatusCode)
	}

	var isolatedPRResponse map[string]interface{}
	parseAndCheckResponse(t, resp, &isolatedPRResponse)
	isolatedPRObj := isolatedPRResponse["pr"].(map[string]interface{})
	isolatedReviewers := isolatedPRObj["assigned_reviewers"].([]interface{})
	if len(isolatedReviewers) != 2 {
		t.Fatalf("Expected 2 reviewers for team of 3, got %d", len(isolatedReviewers))
	}
	closeBody(t, resp)

	// Деактивируем userC
	deactivateData := map[string]interface{}{
		"user_id":   userC,
		"is_active": false,
	}
	deactivateJSON, _ := json.Marshal(deactivateData)

	resp = makeRequest(t, "POST", "/users/setIsActive", deactivateJSON)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to deactivate userC: Expected 200, got %d", resp.StatusCode)
	}
	closeBody(t, resp)

	// Теперь пытаемся переназначить userB - NO_CANDIDATE
	noCandidateData := map[string]string{
		"pull_request_id": isolatedPR,
		"old_user_id":     userB,
	}
	noCandidateJSON, _ := json.Marshal(noCandidateData)

	resp = makeRequest(t, "POST", "/pullRequest/reassign", noCandidateJSON)
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("POST /pullRequest/reassign with no candidates: Expected 409, got %d", resp.StatusCode)
	} else {
		var errorResponse map[string]interface{}
		parseAndCheckResponse(t, resp, &errorResponse)
		checkErrorCode(t, errorResponse, "NO_CANDIDATE")
		closeBody(t, resp)
	}

	// 11. Ошибки для users endpoints
	// Попытка получить ревью для несуществующего пользователя - NOT_FOUND
	resp = makeRequest(t, "GET", "/users/getReview?user_id=nonexistent_user_123", nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("GET /users/getReview with nonexistent user: Expected 404, got %d", resp.StatusCode)
	} else {
		var errorResponse map[string]interface{}
		parseAndCheckResponse(t, resp, &errorResponse)
		checkErrorCode(t, errorResponse, "NOT_FOUND")
		closeBody(t, resp)
	}

	// Попытка деактивировать несуществующего пользователя - NOT_FOUND
	nonexistentUserData := map[string]interface{}{
		"user_id":   "nonexistent_user_456",
		"is_active": false,
	}
	nonexistentUserJSON, _ := json.Marshal(nonexistentUserData)

	resp = makeRequest(t, "POST", "/users/setIsActive", nonexistentUserJSON)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("POST /users/setIsActive with nonexistent user: Expected 404, got %d", resp.StatusCode)
	} else {
		var errorResponse map[string]interface{}
		parseAndCheckResponse(t, resp, &errorResponse)
		checkErrorCode(t, errorResponse, "NOT_FOUND")
		closeBody(t, resp)
	}

	// 12. Проверка отсутствия обязательных параметров
	// GET /team/get без team_name
	resp = makeRequest(t, "GET", "/team/get", nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("GET /team/get without team_name: Expected 400, got %d", resp.StatusCode)
	} else {
		closeBody(t, resp)
	}

	// GET /users/getReview без user_id
	resp = makeRequest(t, "GET", "/users/getReview", nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("GET /users/getReview without user_id: Expected 400, got %d", resp.StatusCode)
	} else {
		closeBody(t, resp)
	}

	t.Log("All error cases tested successfully")
}

func closeBody(t *testing.T, resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		t.Logf("Failed to close response body: %v", err)
	}
}

func checkErrorCode(t *testing.T, errorResponse map[string]interface{}, expectedCode string) {
	errorObj, ok := errorResponse["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("Error response should contain error object, got: %v", errorResponse)
	}

	code, exists := errorObj["code"]
	if !exists {
		t.Fatalf("Error object should contain code field, got: %v", errorObj)
	}

	if code != expectedCode {
		t.Errorf("Expected error code '%s', got '%s'", expectedCode, code)
	}

	message, exists := errorObj["message"]
	if !exists {
		t.Fatalf("Error object should contain message field, got: %v", errorObj)
	}

	if message == "" {
		t.Error("Error message should not be empty")
	}
}

func makeRequest(t *testing.T, method, path string, body []byte) *http.Response {
	req, err := http.NewRequest(method, getBaseURL()+path, bytes.NewBuffer(body))
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

func parseAndCheckResponse(t *testing.T, resp *http.Response, target interface{}) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if len(body) == 0 {
		t.Fatal("Empty response body")
	}

	if err := json.Unmarshal(body, target); err != nil {
		t.Fatalf("Failed to parse JSON response: %v\nResponse body: %s", err, string(body))
	}
}
