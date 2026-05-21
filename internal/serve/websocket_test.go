package serve

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/gorilla/websocket"
)

// TestGenerateClientID guards lvt#331: the prior implementation used
// `string(rune(clientCounter))` which converted the counter to a UTF-8
// codepoint (clientCounter=1 → "\x01" SOH), making IDs unprintable. The
// fix uses strconv.FormatUint to produce the decimal string. This test
// asserts (a) IDs contain no non-printable bytes, (b) consecutive IDs
// are unique, and (c) the suffix is the decimal counter (not a rune).
func TestGenerateClientID(t *testing.T) {
	// Generate a small sequence; uniqueness + printability is the contract.
	seen := make(map[string]bool, 8)
	for i := 0; i < 8; i++ {
		id := generateClientID()
		for _, r := range id {
			if !unicode.IsPrint(r) {
				t.Fatalf("ID %q contains non-printable rune %U; lvt#331 regression", id, r)
			}
		}
		if seen[id] {
			t.Fatalf("duplicate ID %q at iteration %d", id, i)
		}
		seen[id] = true
		// Suffix is "-<decimal>"; the last segment should parse as a positive int.
		if dash := strings.LastIndex(id, "-"); dash >= 0 {
			suffix := id[dash+1:]
			if suffix == "" {
				t.Fatalf("ID %q has empty suffix after '-'", id)
			}
			for _, r := range suffix {
				if r < '0' || r > '9' {
					t.Fatalf("ID %q suffix %q is not all decimal digits; lvt#331 regression", id, suffix)
				}
			}
		} else {
			t.Fatalf("ID %q has no '-' separator", id)
		}
	}
}

func TestWebSocketManager_CreateAndClose(t *testing.T) {
	wsm := NewWebSocketManager()

	if wsm.ClientCount() != 0 {
		t.Errorf("Expected 0 clients, got %d", wsm.ClientCount())
	}

	wsm.Close()

	if wsm.ClientCount() != 0 {
		t.Errorf("Expected 0 clients after close, got %d", wsm.ClientCount())
	}
}

func TestWebSocketManager_ReloadClients(t *testing.T) {
	wsm := NewWebSocketManager()
	defer wsm.Close()

	server := httptest.NewServer(http.HandlerFunc(wsm.HandleWebSocket))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	time.Sleep(100 * time.Millisecond)

	if wsm.ClientCount() != 1 {
		t.Errorf("Expected 1 client, got %d", wsm.ClientCount())
	}

	testData := map[string]interface{}{
		"type": "test",
		"data": "hello",
	}

	wsm.ReloadClients(testData)

	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	var received map[string]interface{}
	if err := json.Unmarshal(message, &received); err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	if received["type"] != "test" {
		t.Errorf("Expected type=test, got %v", received["type"])
	}
	if received["data"] != "hello" {
		t.Errorf("Expected data=hello, got %v", received["data"])
	}
}

func TestWebSocketManager_MultipleClients(t *testing.T) {
	wsm := NewWebSocketManager()
	defer wsm.Close()

	server := httptest.NewServer(http.HandlerFunc(wsm.HandleWebSocket))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	conn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect client 1: %v", err)
	}
	defer conn1.Close()

	conn2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect client 2: %v", err)
	}
	defer conn2.Close()

	time.Sleep(100 * time.Millisecond)

	if wsm.ClientCount() != 2 {
		t.Errorf("Expected 2 clients, got %d", wsm.ClientCount())
	}

	wsm.ReloadClients(map[string]string{"message": "broadcast"})

	for i, conn := range []*websocket.Conn{conn1, conn2} {
		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, message, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("Client %d failed to read: %v", i+1, err)
		}

		var received map[string]string
		if err := json.Unmarshal(message, &received); err != nil {
			t.Fatalf("Client %d failed to unmarshal: %v", i+1, err)
		}

		if received["message"] != "broadcast" {
			t.Errorf("Client %d got wrong message: %v", i+1, received)
		}
	}
}

func TestWebSocketManager_ClientDisconnect(t *testing.T) {
	wsm := NewWebSocketManager()
	defer wsm.Close()

	server := httptest.NewServer(http.HandlerFunc(wsm.HandleWebSocket))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if wsm.ClientCount() != 1 {
		t.Errorf("Expected 1 client, got %d", wsm.ClientCount())
	}

	conn.Close()
	time.Sleep(200 * time.Millisecond)

	if wsm.ClientCount() != 0 {
		t.Errorf("Expected 0 clients after disconnect, got %d", wsm.ClientCount())
	}
}
