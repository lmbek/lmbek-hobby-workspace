package validate

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPerformHealthCheckHTTP(t *testing.T) {
	// Mock HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	err := performHealthCheck(ts.URL)
	if err != nil {
		t.Errorf("expected no error for healthy service, got %v", err)
	}

	// Mock unhealthy HTTP server
	tsErr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer tsErr.Close()

	err = performHealthCheck(tsErr.URL)
	if err == nil {
		t.Error("expected error for unhealthy service, got nil")
	}
}

func TestPerformHealthCheckTCP(t *testing.T) {
	// Start a TCP listener
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start TCP listener: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().String()
	err = performHealthCheck("tcp://" + addr)
	if err != nil {
		t.Errorf("expected no error for healthy TCP service, got %v", err)
	}

	// Test non-existent TCP service
	err = performHealthCheck("tcp://127.0.0.1:1") // Unlikely to have anything on port 1
	if err == nil {
		t.Error("expected error for non-existent TCP service, got nil")
	}
}

func TestPerformHealthCheckUnsupported(t *testing.T) {
	err := performHealthCheck("invalid://localhost")
	if err == nil {
		t.Error("expected error for unsupported protocol, got nil")
	}
}
