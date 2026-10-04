package main

import (
	"encoding/json"
	"testing"
)

func TestFlexibleBool(t *testing.T) {
	cases := []struct {
		raw  string
		want bool
	}{
		{"true", true},
		{"false", false},
		{"\"true\"", true},
		{"\"false\"", false},
		{"1", true},
		{"0", false},
	}
	for _, tc := range cases {
		var got FlexibleBool
		if err := json.Unmarshal([]byte(tc.raw), &got); err != nil {
			t.Fatalf("unmarshal %s: %v", tc.raw, err)
		}
		if bool(got) != tc.want {
			t.Fatalf("unmarshal %s: got %v want %v", tc.raw, got, tc.want)
		}
	}
}

func TestNativeInviteTuple(t *testing.T) {
	raw := `[
		"token123","false","2026-10-04T10:00:00Z","2026-10-09T10:00:00Z",
		"account_only","alice@example.org","","bob",
		"xmpp:bob@example.org?register;preauth=token123",
		"https://example.org/invites/token123"
	]`
	var invite NativeInvite
	if err := json.Unmarshal([]byte(raw), &invite); err != nil {
		t.Fatal(err)
	}
	if bool(invite.Valid) {
		t.Fatal("string false must decode as false")
	}
	if invite.AccountName != "bob" || invite.Inviter != "alice@example.org" {
		t.Fatalf("unexpected invite: %+v", invite)
	}
}

func TestNativeInviteObject(t *testing.T) {
	raw := `{
		"token":"token123",
		"valid":true,
		"created_at":"2026-10-04T10:00:00Z",
		"expires":"2026-10-09T10:00:00Z",
		"type":"account_only",
		"inviter":"alice@example.org",
		"invitee":"",
		"account_name":"",
		"token_uri":"xmpp:example.org?register;preauth=token123",
		"landing_page":"https://example.org/invites/token123"
	}`
	var invite NativeInvite
	if err := json.Unmarshal([]byte(raw), &invite); err != nil {
		t.Fatal(err)
	}
	if !bool(invite.Valid) {
		t.Fatal("boolean true must decode as true")
	}
}

func TestGeneratedInviteTupleAndObject(t *testing.T) {
	cases := []string{
		`["xmpp:example.org?register;preauth=x","https://example.org/invites/x"]`,
		`{"invite_uri":"xmpp:example.org?register;preauth=x","landing_page":"https://example.org/invites/x"}`,
	}
	for _, raw := range cases {
		var invite GeneratedInvite
		if err := json.Unmarshal([]byte(raw), &invite); err != nil {
			t.Fatal(err)
		}
		if invite.InviteURI == "" || invite.LandingPage == "" {
			t.Fatalf("unexpected generated invite: %+v", invite)
		}
	}
}
