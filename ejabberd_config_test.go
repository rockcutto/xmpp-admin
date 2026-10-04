package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEjabberdConfigWithInclude(t *testing.T) {
	dir := t.TempDir()
	include := filepath.Join(dir, "invites.yml")
	main := filepath.Join(dir, "ejabberd.yml")

	if err := os.WriteFile(include, []byte(`
modules:
  mod_invites:
    max_invites: 3
    token_expire_seconds: 604800
    landing_page: auto
    access_create_account: create_account_invite
  mod_register:
    allow_modules:
      - mod_invites
listen:
  - port: 5281
    ip: 127.0.0.1
    module: ejabberd_http
    request_handlers:
      /api: mod_http_api
`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(main, []byte(`
hosts:
  - example.org
include_config_file: invites.yml
`), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadEjabberdConfig(main)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PrimaryHost != "example.org" {
		t.Fatalf("host: got %q", cfg.PrimaryHost)
	}
	if !cfg.InvitesEnabled || cfg.InviteMaxInvites != "3" || cfg.InviteTTLSeconds != 604800 {
		t.Fatalf("invite config not merged: %+v", cfg)
	}
	if cfg.InviteAccessRule != "create_account_invite" {
		t.Fatalf("invite access rule missing: %+v", cfg)
	}
	if !cfg.RegisterEnabled || len(cfg.RegisterAllowModules) != 1 || cfg.RegisterAllowModules[0] != "mod_invites" {
		t.Fatalf("mod_register config missing: %+v", cfg)
	}
	if !cfg.HTTPAPIListenerFound || cfg.HTTPAPIURL != "http://127.0.0.1:5281/api" {
		t.Fatalf("API listener not discovered: %+v", cfg)
	}
	if len(cfg.IncludedFiles) != 1 {
		t.Fatalf("expected one included file, got %v", cfg.IncludedFiles)
	}
}

func TestLoadEjabberdConfigHostOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ejabberd.yml")
	if err := os.WriteFile(path, []byte(`
hosts:
  - first.example
  - second.example
modules:
  mod_invites:
    max_invites: 2
host_config:
  second.example:
    modules:
      mod_invites:
        max_invites: 7
        site_name: Second
`), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadEjabberdConfigForHost(path, "second.example")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PrimaryHost != "second.example" || cfg.InviteMaxInvites != "7" || cfg.InviteSiteName != "Second" {
		t.Fatalf("host override not applied: %+v", cfg)
	}
}

func TestIncludeCycleDoesNotLoop(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.yml")
	b := filepath.Join(dir, "b.yml")
	if err := os.WriteFile(a, []byte("hosts: [example.org]\ninclude_config_file: b.yml\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(b, []byte("include_config_file: a.yml\nmodules:\n  mod_invites: {}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := loadEjabberdConfig(a)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.InvitesEnabled {
		t.Fatal("included config should still be merged")
	}
}


func TestIncludeAllowOnlyAndDisallow(t *testing.T) {
	dir := t.TempDir()
	include := filepath.Join(dir, "extra.yml")
	main := filepath.Join(dir, "ejabberd.yml")

	if err := os.WriteFile(include, []byte(`
listen:
  - port: 9999
    module: ejabberd_http
    request_handlers:
      /api: mod_http_api
modules:
  mod_invites:
    max_invites: 9
`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(main, []byte(`
hosts: [example.org]
include_config_file:
  extra.yml:
    allow_only: [modules]
`), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadEjabberdConfig(main)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.InviteMaxInvites != "9" {
		t.Fatalf("allowed modules were not merged: %+v", cfg)
	}
	if cfg.HTTPAPIListenerFound {
		t.Fatal("listen should have been filtered by allow_only")
	}

	if err := os.WriteFile(main, []byte(`
hosts: [example.org]
include_config_file:
  extra.yml:
    disallow: [modules]
`), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err = loadEjabberdConfig(main)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.InvitesEnabled {
		t.Fatal("modules should have been filtered by disallow")
	}
	if !cfg.HTTPAPIListenerFound {
		t.Fatal("listen should remain available when only modules are disallowed")
	}
}


func TestIPv6HTTPAPIListener(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ejabberd.yml")
	if err := os.WriteFile(path, []byte(`
hosts: [example.org]
listen:
  - port: 5281
    ip: "::1"
    module: ejabberd_http
    request_handlers:
      /api: mod_http_api
modules:
  mod_invites: {}
`), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := loadEjabberdConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HTTPAPIURL != "http://[::1]:5281/api" {
		t.Fatalf("unexpected IPv6 API URL: %q", cfg.HTTPAPIURL)
	}
}


func TestTLSHTTPAPIListenerUsesVHostName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ejabberd.yml")
	if err := os.WriteFile(path, []byte(`
hosts: [example.org]
listen:
  - port: 5443
    ip: "::"
    module: ejabberd_http
    tls: true
    request_handlers:
      /api: mod_http_api
modules:
  mod_invites: {}
`), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := loadEjabberdConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HTTPAPIURL != "https://example.org:5443/api" {
		t.Fatalf("unexpected TLS API URL: %q", cfg.HTTPAPIURL)
	}
}

func TestRepeatedIncludeWithDifferentFilters(t *testing.T) {
	dir := t.TempDir()
	include := filepath.Join(dir, "shared.yml")
	main := filepath.Join(dir, "ejabberd.yml")

	if err := os.WriteFile(include, []byte(`
modules:
  mod_invites:
    max_invites: 11
listen:
  - port: 5281
    ip: 127.0.0.1
    module: ejabberd_http
    request_handlers:
      /api: mod_http_api
`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(main, []byte(`
hosts: [example.org]
include_config_file:
  - shared.yml:
      allow_only: [modules]
  - shared.yml:
      allow_only: [listen]
`), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadEjabberdConfig(main)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.InvitesEnabled || cfg.InviteMaxInvites != "11" {
		t.Fatalf("modules from repeated include missing: %+v", cfg)
	}
	if !cfg.HTTPAPIListenerFound {
		t.Fatalf("listen from repeated include missing: %+v", cfg)
	}
}


func TestNestedRelativeIncludeUsesMainConfigDirectory(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "sub")
	if err := os.MkdirAll(subdir, 0700); err != nil {
		t.Fatal(err)
	}

	main := filepath.Join(dir, "ejabberd.yml")
	first := filepath.Join(subdir, "first.yml")
	second := filepath.Join(dir, "second.yml")

	if err := os.WriteFile(second, []byte(`
modules:
  mod_invites:
    max_invites: 13
`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(first, []byte(`
include_config_file: second.yml
`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(main, []byte(`
hosts: [example.org]
include_config_file: sub/first.yml
`), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadEjabberdConfig(main)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.InvitesEnabled || cfg.InviteMaxInvites != "13" {
		t.Fatalf("nested include should resolve relative to main config directory: %+v", cfg)
	}
}


func TestMainTopLevelOptionReplacesIncludedOption(t *testing.T) {
	dir := t.TempDir()
	include := filepath.Join(dir, "invites.yml")
	main := filepath.Join(dir, "ejabberd.yml")

	if err := os.WriteFile(include, []byte(`
modules:
  mod_invites:
    max_invites: 3
    token_expire_seconds: 604800
  mod_register:
    allow_modules: [mod_invites]
listen:
  - port: 5281
    ip: 127.0.0.1
    module: ejabberd_http
    request_handlers:
      /api: mod_http_api
`), 0600); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(main, []byte(`
hosts: [example.org]
include_config_file: invites.yml
modules:
  mod_invites:
    access_create_account: create_account_invite
`), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadEjabberdConfig(main)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.InvitesEnabled || cfg.InviteAccessRule != "create_account_invite" {
		t.Fatalf("main modules option not applied: %+v", cfg)
	}
	if cfg.InviteMaxInvites != "infinity" || cfg.InviteTTLSeconds != 432000 {
		t.Fatalf("included modules must not partially redefine main modules: %+v", cfg)
	}
	if cfg.RegisterEnabled {
		t.Fatalf("included mod_register must not survive replacement of top-level modules: %+v", cfg)
	}
	if !cfg.HTTPAPIListenerFound {
		t.Fatal("unrelated included listen option should still be present")
	}
}
