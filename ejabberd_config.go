package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type EjabberdConfigSnapshot struct {
	Path                 string
	IncludedFiles        []string
	Hosts                []string
	PrimaryHost          string
	InvitesEnabled       bool
	InviteAccessRule     string
	InviteMaxInvites     string
	InviteTTLSeconds     int
	InviteLandingPage    string
	InviteTemplatesDir   string
	InviteSiteName       string
	InviteDBType         string
	InviteWebchatURL     string
	RegisterEnabled      bool
	RegisterAllowModules []string
	HTTPAPIURL           string
	HTTPAPIListenerFound bool
}

func loadEjabberdConfig(path string) (EjabberdConfigSnapshot, error) {
	return loadEjabberdConfigForHost(path, "")
}

func loadEjabberdConfigForHost(path, hostOverride string) (EjabberdConfigSnapshot, error) {
	rootPath, err := filepath.Abs(path)
	if err != nil {
		rootPath = filepath.Clean(path)
	}
	root, includes, err := loadMergedYAML(rootPath, filepath.Dir(rootPath), map[string]bool{}, 0)
	if err != nil {
		return EjabberdConfigSnapshot{}, err
	}

	out := EjabberdConfigSnapshot{
		Path:                path,
		IncludedFiles:       includes,
		InviteAccessRule:    "none",
		InviteMaxInvites:    "infinity",
		InviteTTLSeconds:    432000,
		InviteLandingPage:   "none",
		InviteTemplatesDir:  "",
		InviteWebchatURL:    "auto",
		InviteDBType:        scalarOrDefault(mapValue(root, "default_db"), "mnesia"),
	}

	if hosts := mapValue(root, "hosts"); hosts != nil {
		out.Hosts = stringSequence(hosts)
	}
	out.PrimaryHost = strings.TrimSpace(hostOverride)
	if out.PrimaryHost == "" && len(out.Hosts) > 0 {
		out.PrimaryHost = out.Hosts[0]
	}
	out.InviteSiteName = out.PrimaryHost

	modules := mapValue(root, "modules")
	defaultDB := out.InviteDBType

	if out.PrimaryHost != "" {
		if hc := mapValue(root, "host_config"); hc != nil {
			if hostCfg := mapValue(hc, out.PrimaryHost); hostCfg != nil {
				if n := mapValue(hostCfg, "default_db"); n != nil {
					defaultDB = scalarOrDefault(n, defaultDB)
				}
				if hostModules := mapValue(hostCfg, "modules"); hostModules != nil {
					modules = hostModules
				}
			}
		}
	}

	out.InviteDBType = defaultDB
	readModulesIntoSnapshot(modules, &out)

	if out.PrimaryHost != "" {
		if ahc := mapValue(root, "append_host_config"); ahc != nil {
			if hostCfg := mapValue(ahc, out.PrimaryHost); hostCfg != nil {
				if n := mapValue(hostCfg, "default_db"); n != nil {
					out.InviteDBType = scalarOrDefault(n, out.InviteDBType)
				}
				if hostModules := mapValue(hostCfg, "modules"); hostModules != nil {
					applyModuleOverrides(hostModules, &out)
				}
			}
		}
	}

	findHTTPAPI(root, &out)
	return out, nil
}

func loadMergedYAML(path, configDir string, active map[string]bool, depth int) (*yaml.Node, []string, error) {
	if depth > 16 {
		return nil, nil, fmt.Errorf("ejabberd config include depth exceeds 16")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = filepath.Clean(path)
	}
	if active[abs] {
		return newMappingNode(), nil, nil
	}
	active[abs] = true
	defer delete(active, abs)

	data, err := os.ReadFile(abs)
	if err != nil {
		return nil, nil, fmt.Errorf("read ejabberd config %s: %w", abs, err)
	}
	var document yaml.Node
	if err := yaml.Unmarshal(data, &document); err != nil {
		return nil, nil, fmt.Errorf("parse ejabberd config %s: %w", abs, err)
	}
	if len(document.Content) == 0 || document.Content[0].Kind != yaml.MappingNode {
		return nil, nil, fmt.Errorf("ejabberd config %s must contain a YAML mapping", abs)
	}
	current := document.Content[0]
	merged := newMappingNode()
	var included []string

	for _, spec := range parseIncludeSpecs(mapValue(current, "include_config_file")) {
		include := spec.Path
		if !filepath.IsAbs(include) {
			include = filepath.Join(configDir, include)
		}
		child, childIncludes, err := loadMergedYAML(include, configDir, active, depth+1)
		if err != nil {
			return nil, nil, err
		}
		child = filterIncludedMapping(child, spec)
		mergeMapping(merged, child)
		included = append(included, include)
		included = append(included, childIncludes...)
	}

	// The main/current file wins over included defaults.
	mergeMapping(merged, current)
	return merged, uniqueStrings(included), nil
}

func newMappingNode() *yaml.Node {
	return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
}

func mergeMapping(dst, src *yaml.Node) {
	if dst == nil || src == nil || dst.Kind != yaml.MappingNode || src.Kind != yaml.MappingNode {
		return
	}
	for i := 0; i+1 < len(src.Content); i += 2 {
		key := src.Content[i]
		value := src.Content[i+1]
		idx := mappingIndex(dst, key.Value)
		if idx < 0 {
			dst.Content = append(dst.Content, cloneNode(key), cloneNode(value))
			continue
		}
		// ejabberd treats a top-level option as one value. If the main file
		// defines it, an included file cannot partially redefine it.
		dst.Content[idx] = cloneNode(key)
		dst.Content[idx+1] = cloneNode(value)
	}
}

func mappingIndex(node *yaml.Node, key string) int {
	if node == nil || node.Kind != yaml.MappingNode {
		return -1
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return i
		}
	}
	return -1
}

func cloneNode(src *yaml.Node) *yaml.Node {
	if src == nil {
		return nil
	}
	dst := *src
	dst.Content = make([]*yaml.Node, len(src.Content))
	for i, child := range src.Content {
		dst.Content[i] = cloneNode(child)
	}
	return &dst
}

type includeSpec struct {
	Path      string
	AllowOnly map[string]bool
	Disallow  map[string]bool
}

func parseIncludeSpecs(node *yaml.Node) []includeSpec {
	if node == nil {
		return nil
	}
	switch node.Kind {
	case yaml.ScalarNode:
		path := strings.TrimSpace(node.Value)
		if path == "" {
			return nil
		}
		return []includeSpec{{Path: path}}
	case yaml.SequenceNode:
		var out []includeSpec
		for _, child := range node.Content {
			out = append(out, parseIncludeSpecs(child)...)
		}
		return out
	case yaml.MappingNode:
		var out []includeSpec
		for i := 0; i+1 < len(node.Content); i += 2 {
			path := strings.TrimSpace(node.Content[i].Value)
			if path == "" {
				continue
			}
			spec := includeSpec{Path: path}
			options := node.Content[i+1]
			if options != nil && options.Kind == yaml.MappingNode {
				spec.AllowOnly = stringSet(mapValue(options, "allow_only"))
				spec.Disallow = stringSet(mapValue(options, "disallow"))
			}
			out = append(out, spec)
		}
		return out
	default:
		return nil
	}
}

func stringSet(node *yaml.Node) map[string]bool {
	values := stringSequence(node)
	if len(values) == 0 {
		return nil
	}
	out := make(map[string]bool, len(values))
	for _, value := range values {
		out[value] = true
	}
	return out
}

func filterIncludedMapping(node *yaml.Node, spec includeSpec) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return node
	}
	if len(spec.AllowOnly) == 0 && len(spec.Disallow) == 0 {
		return node
	}
	out := newMappingNode()
	for i := 0; i+1 < len(node.Content); i += 2 {
		key := node.Content[i].Value
		if len(spec.AllowOnly) > 0 && !spec.AllowOnly[key] {
			continue
		}
		if spec.Disallow[key] {
			continue
		}
		out.Content = append(out.Content, cloneNode(node.Content[i]), cloneNode(node.Content[i+1]))
	}
	return out
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]bool, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = filepath.Clean(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}

func readModulesIntoSnapshot(modules *yaml.Node, out *EjabberdConfigSnapshot) {
	if modules == nil || modules.Kind != yaml.MappingNode {
		return
	}
	if modInvites := mapValue(modules, "mod_invites"); modInvites != nil {
		out.InvitesEnabled = true
		readInviteOptions(modInvites, out)
	}
	if modRegister := mapValue(modules, "mod_register"); modRegister != nil {
		out.RegisterEnabled = true
		if modRegister.Kind == yaml.MappingNode {
			out.RegisterAllowModules = stringSequence(mapValue(modRegister, "allow_modules"))
		}
	}
}

func applyModuleOverrides(modules *yaml.Node, out *EjabberdConfigSnapshot) {
	readModulesIntoSnapshot(modules, out)
}

func readInviteOptions(modInvites *yaml.Node, out *EjabberdConfigSnapshot) {
	if modInvites == nil || modInvites.Kind != yaml.MappingNode {
		return
	}
	if n := mapValue(modInvites, "access_create_account"); n != nil {
		out.InviteAccessRule = scalarOrDefault(n, out.InviteAccessRule)
	}
	if n := mapValue(modInvites, "max_invites"); n != nil {
		out.InviteMaxInvites = scalarOrDefault(n, out.InviteMaxInvites)
	}
	if n := mapValue(modInvites, "token_expire_seconds"); n != nil {
		if v, err := strconv.Atoi(scalarString(n)); err == nil {
			out.InviteTTLSeconds = v
		}
	}
	if n := mapValue(modInvites, "landing_page"); n != nil {
		out.InviteLandingPage = scalarOrDefault(n, out.InviteLandingPage)
	}
	if n := mapValue(modInvites, "templates_dir"); n != nil {
		out.InviteTemplatesDir = scalarString(n)
	}
	if n := mapValue(modInvites, "site_name"); n != nil {
		out.InviteSiteName = scalarOrDefault(n, out.InviteSiteName)
	}
	if n := mapValue(modInvites, "db_type"); n != nil {
		out.InviteDBType = scalarOrDefault(n, out.InviteDBType)
	}
	if n := mapValue(modInvites, "webchat_url"); n != nil {
		out.InviteWebchatURL = scalarOrDefault(n, out.InviteWebchatURL)
	}
}

func findHTTPAPI(doc *yaml.Node, out *EjabberdConfigSnapshot) {
	listeners := mapValue(doc, "listen")
	if listeners == nil || listeners.Kind != yaml.SequenceNode {
		return
	}
	for _, listener := range listeners.Content {
		if listener.Kind != yaml.MappingNode || scalarString(mapValue(listener, "module")) != "ejabberd_http" {
			continue
		}
		handlers := mapValue(listener, "request_handlers")
		if handlers == nil || handlers.Kind != yaml.MappingNode {
			continue
		}
		for i := 0; i+1 < len(handlers.Content); i += 2 {
			if scalarString(handlers.Content[i+1]) != "mod_http_api" {
				continue
			}
			apiPath := strings.TrimSpace(scalarString(handlers.Content[i]))
			if apiPath == "" {
				apiPath = "/api"
			}
			if !strings.HasPrefix(apiPath, "/") {
				apiPath = "/" + apiPath
			}
			ip := normalizeLocalHost(scalarString(mapValue(listener, "ip")))
			port := scalarOrDefault(mapValue(listener, "port"), "5280")
			scheme := "http"
			if strings.EqualFold(scalarString(mapValue(listener, "tls")), "true") {
				scheme = "https"
			}
			out.HTTPAPIURL = scheme + "://" + net.JoinHostPort(ip, port) + apiPath
			out.HTTPAPIListenerFound = true
			return
		}
	}
}

func mapValue(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	idx := mappingIndex(node, key)
	if idx < 0 {
		return nil
	}
	return node.Content[idx+1]
}

func stringSequence(node *yaml.Node) []string {
	if node == nil {
		return nil
	}
	if node.Kind == yaml.SequenceNode {
		out := make([]string, 0, len(node.Content))
		for _, item := range node.Content {
			if s := strings.TrimSpace(scalarString(item)); s != "" {
				out = append(out, s)
			}
		}
		return out
	}
	if s := strings.TrimSpace(scalarString(node)); s != "" {
		return []string{s}
	}
	return nil
}

func scalarString(node *yaml.Node) string {
	if node == nil {
		return ""
	}
	return strings.TrimSpace(node.Value)
}

func scalarOrDefault(node *yaml.Node, fallback string) string {
	if s := scalarString(node); s != "" {
		return s
	}
	return fallback
}

func normalizeLocalHost(ip string) string {
	ip = strings.TrimSpace(ip)
	switch ip {
	case "", "0.0.0.0":
		return "127.0.0.1"
	case "::":
		return "::1"
	default:
		return ip
	}
}

func resolveEjabberdConfigPath(explicit string) string {
	if p := strings.TrimSpace(explicit); p != "" {
		return p
	}
	candidates := []string{
		"/etc/ejabberd/ejabberd.yml",
		"/opt/ejabberd/conf/ejabberd.yml",
		"./ejabberd.yml",
	}
	for _, p := range candidates {
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p
		}
	}
	return "/etc/ejabberd/ejabberd.yml"
}
