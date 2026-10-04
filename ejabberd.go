package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type EjabberdClient struct {
	baseURL  string
	username string
	password string
	client   *http.Client
}

type FlexibleBool bool

func (b *FlexibleBool) UnmarshalJSON(data []byte) error {
	raw := strings.TrimSpace(string(data))
	switch raw {
	case "true", `"true"`, "1", `"1"`:
		*b = true
		return nil
	case "false", `"false"`, "0", `"0"`, "null", `""`:
		*b = false
		return nil
	default:
		return fmt.Errorf("invalid boolean value %q", raw)
	}
}

type NativeInvite struct {
	Token       string       `json:"token"`
	Valid       FlexibleBool `json:"valid"`
	CreatedAt   string       `json:"created_at"`
	Expires     string       `json:"expires"`
	Type        string       `json:"type"`
	Inviter     string       `json:"inviter"`
	Invitee     string       `json:"invitee"`
	AccountName string       `json:"account_name"`
	TokenURI    string       `json:"token_uri"`
	LandingPage string       `json:"landing_page"`
}

func (i *NativeInvite) UnmarshalJSON(data []byte) error {
	type alias NativeInvite
	var object alias
	if len(data) > 0 && data[0] == '{' {
		if err := json.Unmarshal(data, &object); err != nil {
			return err
		}
		*i = NativeInvite(object)
		return nil
	}

	var tuple []json.RawMessage
	if err := json.Unmarshal(data, &tuple); err != nil {
		return err
	}
	if len(tuple) != 10 {
		return fmt.Errorf("unexpected invite tuple length %d", len(tuple))
	}
	fields := []any{
		&i.Token, &i.Valid, &i.CreatedAt, &i.Expires, &i.Type,
		&i.Inviter, &i.Invitee, &i.AccountName, &i.TokenURI, &i.LandingPage,
	}
	for idx, dst := range fields {
		if err := json.Unmarshal(tuple[idx], dst); err != nil {
			return fmt.Errorf("decode invite field %d: %w", idx, err)
		}
	}
	return nil
}

type GeneratedInvite struct {
	InviteURI   string `json:"invite_uri"`
	LandingPage string `json:"landing_page"`
}

func (i *GeneratedInvite) UnmarshalJSON(data []byte) error {
	type alias GeneratedInvite
	var object alias
	if len(data) > 0 && data[0] == '{' {
		if err := json.Unmarshal(data, &object); err != nil {
			return err
		}
		*i = GeneratedInvite(object)
		return nil
	}
	var tuple []string
	if err := json.Unmarshal(data, &tuple); err != nil {
		return err
	}
	if len(tuple) != 2 {
		return fmt.Errorf("unexpected generated invite tuple length %d", len(tuple))
	}
	i.InviteURI = tuple[0]
	i.LandingPage = tuple[1]
	return nil
}

func (c *EjabberdClient) call(ctx context.Context, command string, input any, output any) error {
	body, err := json.Marshal(input)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/"+command, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.username != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("ejabberd API unavailable: %w", err)
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ejabberd API %s failed: HTTP %d", command, resp.StatusCode)
	}
	if output == nil || len(bytes.TrimSpace(payload)) == 0 {
		return nil
	}
	if err := json.Unmarshal(payload, output); err != nil {
		return fmt.Errorf("ejabberd API %s returned invalid JSON: %w", command, err)
	}
	return nil
}

func (c *EjabberdClient) ListInvites(ctx context.Context, host string) ([]NativeInvite, error) {
	var invites []NativeInvite
	if err := c.call(ctx, "list_invites", map[string]string{"host": host}, &invites); err != nil {
		return nil, err
	}
	return invites, nil
}

func (c *EjabberdClient) GenerateInvite(ctx context.Context, host, username string) (GeneratedInvite, error) {
	var invite GeneratedInvite
	command := "generate_invite"
	input := map[string]string{"host": host}
	if username != "" {
		command = "generate_invite_with_username"
		input["username"] = username
	}
	if err := c.call(ctx, command, input, &invite); err != nil {
		return GeneratedInvite{}, err
	}
	return invite, nil
}

func (c *EjabberdClient) ExpireInvite(ctx context.Context, host, token string) error {
	var raw json.RawMessage
	if err := c.call(ctx, "expire_invite_by_token", map[string]string{
		"host":  host,
		"token": token,
	}, &raw); err != nil {
		return err
	}
	if len(raw) == 0 {
		return nil
	}

	var code int
	if err := json.Unmarshal(raw, &code); err == nil {
		if code == 0 {
			return nil
		}
		return fmt.Errorf("ejabberd refused to expire invite: code %d", code)
	}

	var value string
	if err := json.Unmarshal(raw, &value); err == nil {
		switch strings.ToLower(strings.TrimSpace(value)) {
		case "ok", "0":
			return nil
		default:
			if parsed, err := strconv.Atoi(value); err == nil && parsed == 0 {
				return nil
			}
			return fmt.Errorf("ejabberd refused to expire invite")
		}
	}
	return fmt.Errorf("unexpected ejabberd response to expire_invite_by_token")
}
