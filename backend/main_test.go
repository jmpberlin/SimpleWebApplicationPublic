package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		want         string
	}{
		{
			name:         "returns environment variable when set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "custom",
			want:         "custom",
		},
		{
			name:         "returns default when env var not set",
			key:          "NONEXISTENT_VAR",
			defaultValue: "default",
			envValue:     "",
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			got := getEnv(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPort(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{
			name:     "returns default port when PORT not set",
			envValue: "",
			want:     "8081",
		},
		{
			name:     "returns custom port when PORT is set",
			envValue: "9000",
			want:     "9000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("PORT", tt.envValue)
				defer os.Unsetenv("PORT")
			} else {
				os.Unsetenv("PORT")
			}

			got := getPort()
			if got != tt.want {
				t.Errorf("getPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	db = mockDB
	return mockDB, mock
}

func TestHelloHandler(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	// Set up expectation: when querying for ID 1, return "Hello from the database!"
	rows := sqlmock.NewRows([]string{"message"}).AddRow("Hello from the database!")
	mock.ExpectQuery("SELECT message FROM messages WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	helloHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expectedBody := "Hello from the database!"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, w.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestByeHandler(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"message"}).AddRow("Goodbye from the database!")
	mock.ExpectQuery("SELECT message FROM messages WHERE id = \\$1").
		WithArgs(2).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/bye", nil)
	w := httptest.NewRecorder()

	byeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expectedBody := "Goodbye from the database!"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, w.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestLandingPageHandler(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"message"}).AddRow("Knock knock! Who is there? The database!")
	mock.ExpectQuery("SELECT message FROM messages WHERE id = \\$1").
		WithArgs(3).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	langingPageHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expectedBody := "Knock knock! Who is there? The database!"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, w.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestImpressumHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/impressum", nil)
	w := httptest.NewRecorder()

	impressumHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expectedBody := "Copyright (c) 2025 by github.com/jmpberlin"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, w.Body.String())
	}
}

func TestHandlerDatabaseError(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	// Simulate database error
	mock.ExpectQuery("SELECT message FROM messages WHERE id = \\$1").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	helloHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	expectedBody := "Database error\n"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, w.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetMessageByID(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	tests := []struct {
		name        string
		id          int
		mockMessage string
		mockError   error
		wantErr     bool
	}{
		{
			name:        "successfully retrieves message",
			id:          1,
			mockMessage: "Test message",
			mockError:   nil,
			wantErr:     false,
		},
		{
			name:        "returns error when message not found",
			id:          999,
			mockMessage: "",
			mockError:   sql.ErrNoRows,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockError != nil {
				mock.ExpectQuery("SELECT message FROM messages WHERE id = \\$1").
					WithArgs(tt.id).
					WillReturnError(tt.mockError)
			} else {
				rows := sqlmock.NewRows([]string{"message"}).AddRow(tt.mockMessage)
				mock.ExpectQuery("SELECT message FROM messages WHERE id = \\$1").
					WithArgs(tt.id).
					WillReturnRows(rows)
			}

			got, err := getMessageByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("getMessageByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.mockMessage {
				t.Errorf("getMessageByID() = %v, want %v", got, tt.mockMessage)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
