package main

import (
	"context"
	"fmt"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSplitFullJID(t *testing.T) {
	user, host, resource := splitFullJID("alice@example.org/mobile")
	if user != "alice" || host != "example.org" || resource != "mobile" {
		t.Fatalf("unexpected split: user=%q host=%q resource=%q", user, host, resource)
	}

	user, host, resource = splitFullJID("room@conference.example.org")
	if user != "room" || host != "conference.example.org" || resource != "" {
		t.Fatalf("unexpected room split: user=%q host=%q resource=%q", user, host, resource)
	}
}

func TestSessionsForHostFiltersOtherVhosts(t *testing.T) {
	rows := sessionsForHost([]string{
		"bob@other.example/web",
		"alice@example.org/mobile",
		"alice@example.org/desktop",
	}, "example.org")
	if len(rows) != 2 {
		t.Fatalf("got %d rows, want 2: %+v", len(rows), rows)
	}
	for _, row := range rows {
		_, host, _ := splitFullJID(row.JID)
		if host != "example.org" {
			t.Fatalf("foreign vhost leaked into rows: %+v", row)
		}
	}
}

func TestReadOnlyOperationsAPIWire(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/registered_users":
			var input map[string]string
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				t.Fatal(err)
			}
			if input["host"] != "example.org" {
				t.Fatalf("registered_users host = %q", input["host"])
			}
			_, _ = w.Write([]byte(`["alice","bob"]`))
		case "/connected_users":
			_, _ = w.Write([]byte(`["alice@example.org/mobile","other@other.org/web"]`))
		case "/muc_online_rooms":
			var input map[string]string
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				t.Fatal(err)
			}
			if input["service"] != "global" {
				t.Fatalf("muc_online_rooms service = %q", input["service"])
			}
			_, _ = w.Write([]byte(`["room@conference.example.org"]`))
		case "/status":
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer api.Close()

	client := &EjabberdClient{baseURL: api.URL, client: api.Client()}
	ctx := context.Background()

	users, err := client.RegisteredUsers(ctx, "example.org")
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 || users[0] != "alice" {
		t.Fatalf("unexpected users: %v", users)
	}

	sessions, err := client.ConnectedUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 2 {
		t.Fatalf("unexpected sessions: %v", sessions)
	}

	rooms, err := client.MUCOnlineRooms(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(rooms) != 1 || rooms[0] != "room@conference.example.org" {
		t.Fatalf("unexpected rooms: %v", rooms)
	}

	if err := client.Status(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestAPIProblemExplainsForbiddenAccess(t *testing.T) {
	summary, technical := apiProblem(langRU, "users", fmt.Errorf("ejabberd API registered_users failed: HTTP 403"))
	if summary != "У XMPP Admin нет права читать список пользователей." {
		t.Fatalf("unexpected summary: %q", summary)
	}
	if !strings.Contains(technical, "HTTP 403") {
		t.Fatalf("technical detail must preserve HTTP status: %q", technical)
	}
}
