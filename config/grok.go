package config

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/qiangyt/jog/static/grok_patterns"
	"github.com/qiangyt/jog/util"
	"github.com/vjeantet/grok"
)

// GrokT ...
type GrokT struct {

// DefaultGrokPatternsDir ...
func DefaultGrokPatternsDir() string {
	return JogHomeDir("grok-patterns")
}

// Grok ...
type Grok = *GrokT

// SaveDefaultGrokPatternFile ...
func SaveDefaultGrokPatternFile(patternFileName string, patternFileContent string) {
	dir := DefaultGrokPatternsDir()
	util.ReplaceFile(filepath.Join(dir, patternFileName), []byte(patternFileContent))
}

// ResetDefaultGrokPatternsDir ...
func ResetDefaultGrokPatternsDir() {
	dir := DefaultGrokPatternsDir()
	util.RemoveDir(dir)

	InitDefaultGrokPatternsDir()
}

// InitDefaultGrokPatternsDir ...
func InitDefaultGrokPatternsDir() {
	jogHomeDir := JogHomeDir()

	licensePath := filepath.Join(jogHomeDir, "vjeantet-grok.LICENSE")
	util.WriteFileIfNotFound(licensePath, []byte(grok_patterns.LICENSE))

	readmePath := filepath.Join(jogHomeDir, "vjeantet-grok.README.md")
	util.WriteFileIfNotFound(readmePath, []byte(grok_patterns.README_md))

	dir := DefaultGrokPatternsDir()
	if util.DirExists(dir) {
		return
	}
	util.MkdirAll(dir)

	SaveDefaultGrokPatternFile("aws", grok_patterns.Aws)
	SaveDefaultGrokPatternFile("bro", grok_patterns.Bro)
	SaveDefaultGrokPatternFile("firewalls", grok_patterns.Firewalls)
	SaveDefaultGrokPatternFile("haproxy", grok_patterns.Haproxy)
	SaveDefaultGrokPatternFile("junos", grok_patterns.Junos)
	SaveDefaultGrokPatternFile("linux-syslog", grok_patterns.Linux_syslog)
	SaveDefaultGrokPatternFile("mcollective-patterns", grok_patterns.Mcollective_patterns)
	SaveDefaultGrokPatternFile("nagios", grok_patterns.Nagios)
	SaveDefaultGrokPatternFile("rails", grok_patterns.Rails)
	SaveDefaultGrokPatternFile("redis", grok_patterns.Redis)
	SaveDefaultGrokPatternFile("bacula", grok_patterns.Bacula)
	SaveDefaultGrokPatternFile("exim", grok_patterns.Exim)
	SaveDefaultGrokPatternFile("grok-patterns", grok_patterns.Grok_patterns)
	SaveDefaultGrokPatternFile("java", grok_patterns.Java)
	SaveDefaultGrokPatternFile("mcollective", grok_patterns.Mcollective)
	SaveDefaultGrokPatternFile("mongodb", grok_patterns.Mongodb)
	SaveDefaultGrokPatternFile("postgresql", grok_patterns.Postgresql)
	SaveDefaultGrokPatternFile("ruby", grok_patterns.Ruby)
}

// Init ...
func (i Grok) Init(cfg Configuration) {

	InitDefaultGrokPatternsDir()

	i.grok, _ = grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})


// UnmarshalYAML ...
func (i Grok) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return UnmarshalYAML(i, unmarshal)
}

// MarshalYAML ...
func (i Grok) MarshalYAML() (interface{}, error) {
	return MarshalYAML(i)
}

// Reset ...
func (i Grok) Reset() {
	i.Uses = make([]string, 0)
	i.Patterns = make(map[string]string)
}

// FromMap ...
func (i Grok) FromMap(m map[string]interface{}) error {
	i.Patterns = util.ExtractFromMap(m, "patterns").(map[string]string)

	i.Uses = util.ExtractFromMap(m, "uses").([]string)
	for _, usedPatternName := range i.Uses {
		if _, found := i.Patterns[usedPatternName]; !found {
			return fmt.Errorf("using pattern '%s' but not defined in available patterns", usedPatternName)
		}
	}

	return nil
}

// ToMap ...
func (i Grok) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	r["uses"] = i.Uses

	patterns := make(map[string]string)
	for patternName, patternExpr := range i.Patterns {
		patterns[patternName] = patternExpr
	}
	r["patterns"] = patterns

	return r
}
