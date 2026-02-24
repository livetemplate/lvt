package knowledge

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/livetemplate/lvt/internal/telemetry"
)

// KnowledgeBase loads and queries evolution patterns.
type KnowledgeBase struct {
	patternsFile string
	patterns     []*Pattern
	byID         map[string]*Pattern
	mu           sync.RWMutex
}

// New creates a KnowledgeBase by parsing the patterns file at the given path.
func New(patternsFilePath string) (*KnowledgeBase, error) {
	kb := &KnowledgeBase{patternsFile: patternsFilePath}
	if err := kb.Reload(); err != nil {
		return nil, err
	}
	return kb, nil
}

// Reload re-parses the patterns file.
func (kb *KnowledgeBase) Reload() error {
	patterns, err := Parse(kb.patternsFile)
	if err != nil {
		return fmt.Errorf("load knowledge base: %w", err)
	}

	byID := make(map[string]*Pattern, len(patterns))
	for _, p := range patterns {
		byID[p.ID] = p
	}

	kb.mu.Lock()
	kb.patterns = patterns
	kb.byID = byID
	kb.mu.Unlock()

	return nil
}

// LookupFixes returns all fixes from patterns matching the given error.
func (kb *KnowledgeBase) LookupFixes(err telemetry.GenerationError) []FixTemplate {
	kb.mu.RLock()
	defer kb.mu.RUnlock()

	var fixes []FixTemplate
	for _, p := range kb.patterns {
		if p.Match(err) {
			fixes = append(fixes, p.Fixes...)
		}
	}
	return fixes
}

// ListPatterns returns all loaded patterns.
func (kb *KnowledgeBase) ListPatterns() []*Pattern {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	out := make([]*Pattern, len(kb.patterns))
	copy(out, kb.patterns)
	return out
}

// GetPattern returns a pattern by ID.
func (kb *KnowledgeBase) GetPattern(id string) (*Pattern, bool) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	p, ok := kb.byID[id]
	return p, ok
}

// NewFromPatterns creates a KnowledgeBase from pre-built patterns (for testing).
func NewFromPatterns(patterns []*Pattern) *KnowledgeBase {
	byID := make(map[string]*Pattern, len(patterns))
	for _, p := range patterns {
		byID[p.ID] = p
	}
	return &KnowledgeBase{patterns: patterns, byID: byID}
}

// PatternCount returns the number of loaded patterns.
func (kb *KnowledgeBase) PatternCount() int {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	return len(kb.patterns)
}

// FindPatternsFile locates the evolution/patterns.md file by searching:
// 1. Relative to the lvt binary (for installed lvt)
// 2. Current working directory (for development)
// 3. User config directory ~/.config/lvt/patterns.md (user override)
func FindPatternsFile() (string, error) {
	// 1. Relative to binary
	exe, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exe)
		candidates := []string{
			filepath.Join(exeDir, "evolution", "patterns.md"),
			filepath.Join(exeDir, "..", "evolution", "patterns.md"),
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				return c, nil
			}
		}
	}

	// 2. Current working directory and ancestors (walk up to find repo root)
	cwd, err := os.Getwd()
	if err == nil {
		dir := cwd
		for {
			p := filepath.Join(dir, "evolution", "patterns.md")
			if _, err := os.Stat(p); err == nil {
				return p, nil
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	// 3. User config directory
	home, err := os.UserHomeDir()
	if err == nil {
		p := filepath.Join(home, ".config", "lvt", "patterns.md")
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("evolution/patterns.md not found (checked binary dir, cwd, ~/.config/lvt/)")
}
