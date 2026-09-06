package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
	"gopkg.in/yaml.v3"
)

const ConfigFileName = ".mdlinker.yaml"

// KnownRuleIDs は v0.1 で公開する Rule。MDL005 は欠番。
var KnownRuleIDs = []string{
	"MDL001",
	"MDL002",
	"MDL003",
	"MDL004",
	"MDL006",
	"MDL007",
}

// Config は include / exclude / Rule Severity の最終設定。
type Config struct {
	Include []string                   `yaml:"include"`
	Exclude []string                   `yaml:"exclude"`
	Rules   map[string]engine.Severity `yaml:"rules"`
}

type fileConfig struct {
	Include []string          `yaml:"include"`
	Exclude []string          `yaml:"exclude"`
	Rules   map[string]string `yaml:"rules"`
}

// Default は設定ファイルが無いときの既定値を返す。
func Default() Config {
	return Config{
		Include: []string{"**/*.md", "**/*.mdx"},
		Exclude: []string{".git/**", "node_modules/**", "vendor/**"},
		Rules: map[string]engine.Severity{
			"MDL001": engine.SeverityError,
			"MDL002": engine.SeverityError,
			"MDL003": engine.SeverityWarning,
			"MDL004": engine.SeverityError,
			"MDL006": engine.SeverityWarning,
			"MDL007": engine.SeverityError,
		},
	}
}

// LoadBytes は YAML バイト列を既定設定の上に合成する。
func LoadBytes(data []byte) (Config, error) {
	cfg := Default()
	var parsed fileConfig
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return Config{}, fmt.Errorf("config.LoadBytes: parse yaml: %w", err)
	}

	if len(parsed.Include) > 0 {
		cfg.Include = parsed.Include
	}
	if len(parsed.Exclude) > 0 {
		cfg.Exclude = parsed.Exclude
	}

	known := map[string]struct{}{}
	for _, id := range KnownRuleIDs {
		known[id] = struct{}{}
	}

	for id, raw := range parsed.Rules {
		if _, ok := known[id]; !ok {
			return Config{}, fmt.Errorf("config.LoadBytes: unknown rule id %q", id)
		}
		severity, err := parseSeverity(raw)
		if err != nil {
			return Config{}, fmt.Errorf("config.LoadBytes: rule %s: %w", id, err)
		}
		cfg.Rules[id] = severity
	}

	return cfg, nil
}

// Load は指定パスの設定ファイルを読み込む。
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("config.Load: path=%q: %w", path, err)
	}
	return LoadBytes(data)
}

// LoadOrDefault は root 配下の .mdlinker.yaml を読む。無ければ既定値。
func LoadOrDefault(root string) (Config, error) {
	path := filepath.Join(root, ConfigFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return Config{}, fmt.Errorf("config.LoadOrDefault: path=%q: %w", path, err)
	}
	return LoadBytes(data)
}

func parseSeverity(raw string) (engine.Severity, error) {
	switch engine.Severity(raw) {
	case engine.SeverityError, engine.SeverityWarning, engine.SeverityOff:
		return engine.Severity(raw), nil
	default:
		return "", fmt.Errorf("invalid severity %q (allowed: error, warning, off)", raw)
	}
}
