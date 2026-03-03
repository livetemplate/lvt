package commands

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/livetemplate/lvt/internal/evolution"
	"github.com/livetemplate/lvt/internal/evolution/knowledge"
	"github.com/livetemplate/lvt/internal/telemetry"
)

const defaultLookbackDays = 30

// Evolution is the main router for evolution system commands.
func Evolution(args []string) error {
	if len(args) == 0 {
		printEvolutionHelp()
		return nil
	}

	switch args[0] {
	case "status":
		return evolutionStatus(args[1:])
	case "metrics":
		return evolutionMetrics(args[1:])
	case "failures":
		return evolutionFailures(args[1:])
	case "patterns":
		return evolutionPatterns(args[1:])
	case "propose":
		return evolutionPropose(args[1:])
	case "apply":
		return evolutionApply(args[1:])
	case "components":
		return evolutionComponents(args[1:])
	case "upstream-status":
		return evolutionUpstreamStatus(args[1:])
	case "help", "--help", "-h":
		printEvolutionHelp()
		return nil
	default:
		return fmt.Errorf("unknown evolution subcommand: %s\nRun 'lvt evolution help' for usage", args[0])
	}
}

func printEvolutionHelp() {
	fmt.Println("Evolution System — Self-improving feedback loop for code generation")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  lvt evolution <command> [args...]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  status                          Show event counts and success rates")
	fmt.Println("  metrics                         Show per-command metrics")
	fmt.Println("  failures [--last N]             List recent failed events")
	fmt.Println("  patterns                        List all known patterns from knowledge base")
	fmt.Println("  components [--days N]            Show per-component health dashboard")
	fmt.Println("  propose <event-id>              Propose fixes for a specific event")
	fmt.Println("  apply <fix-id> [--dry-run]      Apply a proposed fix [coming soon]")
	fmt.Println("  upstream-status                 Show upstream pattern status")
}

func evolutionStatus(_ []string) error {
	collector := telemetry.NewCollector()
	defer collector.Close()

	if !collector.IsEnabled() {
		fmt.Println("Telemetry is disabled (set LVT_TELEMETRY=true to enable)")
		return nil
	}

	ctx := context.Background()
	since := time.Now().AddDate(0, 0, -defaultLookbackDays)
	total, successes, err := collector.Store().CountBySuccess(ctx, since)
	if err != nil {
		return fmt.Errorf("count events: %w", err)
	}

	// Load knowledge base
	patternsFile, err := knowledge.FindPatternsFile()
	patternCount := 0
	if err == nil {
		if kb, err := knowledge.New(patternsFile); err == nil {
			patternCount = kb.PatternCount()
		}
	}

	failures := total - successes
	successRate := float64(0)
	failRate := float64(0)
	if total > 0 {
		successRate = float64(successes) / float64(total) * 100
		failRate = float64(failures) / float64(total) * 100
	}

	fmt.Println("Evolution System Status")
	fmt.Println("=======================")
	fmt.Printf("Events (last %d days): %d\n", defaultLookbackDays, total)
	fmt.Printf("  Successes: %d (%.1f%%)\n", successes, successRate)
	fmt.Printf("  Failures:  %d (%.1f%%)\n", failures, failRate)
	fmt.Printf("Knowledge Base: %d patterns\n", patternCount)

	return nil
}

func evolutionMetrics(_ []string) error {
	collector := telemetry.NewCollector()
	defer collector.Close()

	if !collector.IsEnabled() {
		fmt.Println("Telemetry is disabled (set LVT_TELEMETRY=true to enable)")
		return nil
	}

	ctx := context.Background()
	since := time.Now().AddDate(0, 0, -defaultLookbackDays)
	events, err := collector.Store().List(ctx, telemetry.ListOptions{Since: since})
	if err != nil {
		return fmt.Errorf("list events: %w", err)
	}

	// Aggregate by command
	type stats struct {
		total   int
		success int
		totalMs int64
	}
	byCommand := make(map[string]*stats)
	for _, e := range events {
		s, ok := byCommand[e.Command]
		if !ok {
			s = &stats{}
			byCommand[e.Command] = s
		}
		s.total++
		if e.Success {
			s.success++
		}
		s.totalMs += e.DurationMs
	}

	fmt.Printf("Per-Command Metrics (last %d days)\n", defaultLookbackDays)
	fmt.Println("==================================")
	fmt.Printf("%-20s %8s %8s %8s %10s\n", "Command", "Total", "Success", "Rate", "Avg (ms)")
	fmt.Println(strings.Repeat("-", 60))

	cmds := make([]string, 0, len(byCommand))
	for cmd := range byCommand {
		cmds = append(cmds, cmd)
	}
	sort.Strings(cmds)

	for _, cmd := range cmds {
		s := byCommand[cmd]
		rate := float64(0)
		if s.total > 0 {
			rate = float64(s.success) / float64(s.total) * 100
		}
		avgMs := int64(0)
		if s.total > 0 {
			avgMs = s.totalMs / int64(s.total)
		}
		fmt.Printf("%-20s %8d %8d %7.1f%% %10d\n", cmd, s.total, s.success, rate, avgMs)
	}

	if len(byCommand) == 0 {
		fmt.Println("  No events recorded yet.")
	}

	return nil
}

func evolutionFailures(args []string) error {
	last := 10
	for i := 0; i < len(args); i++ {
		if args[i] == "--last" && i+1 < len(args) {
			if n, err := strconv.Atoi(args[i+1]); err == nil && n > 0 {
				last = n
			}
			i++
		}
	}

	collector := telemetry.NewCollector()
	defer collector.Close()

	if !collector.IsEnabled() {
		fmt.Println("Telemetry is disabled (set LVT_TELEMETRY=true to enable)")
		return nil
	}

	ctx := context.Background()
	successFalse := false
	events, err := collector.Store().List(ctx, telemetry.ListOptions{
		SuccessOnly: &successFalse,
		Limit:       last,
	})
	if err != nil {
		return fmt.Errorf("list failures: %w", err)
	}

	fmt.Printf("Recent Failures (last %d)\n", last)
	fmt.Println("========================")

	if len(events) == 0 {
		fmt.Println("  No failures recorded.")
		return nil
	}

	fmt.Printf("%-36s %-16s %-30s %s\n", "Event ID", "Command", "Error", "Timestamp")
	fmt.Println(strings.Repeat("-", 100))

	for _, e := range events {
		errSummary := ""
		if len(e.Errors) > 0 {
			errSummary = truncate(e.Errors[0].Message, 30)
		}
		ts := e.Timestamp.Format("2006-01-02 15:04")
		fmt.Printf("%-36s %-16s %-30s %s\n", e.ID, e.Command, errSummary, ts)
	}

	return nil
}

func evolutionPatterns(_ []string) error {
	patternsFile, err := knowledge.FindPatternsFile()
	if err != nil {
		return fmt.Errorf("find patterns file: %w", err)
	}

	kb, err := knowledge.New(patternsFile)
	if err != nil {
		return fmt.Errorf("load knowledge base: %w", err)
	}

	patterns := kb.ListPatterns()

	fmt.Printf("Knowledge Base Patterns (%d total)\n", len(patterns))
	fmt.Println("==================================")
	fmt.Printf("%-30s %-35s %10s %8s %8s\n", "ID", "Name", "Confidence", "Fixes", "Upstream")
	fmt.Println(strings.Repeat("-", 95))

	for _, p := range patterns {
		upstream := ""
		if p.UpstreamRepo != "" {
			upstream = "yes"
		}
		fmt.Printf("%-30s %-35s %10.2f %8d %8s\n",
			p.ID, truncate(p.Name, 33), p.Confidence, len(p.Fixes), upstream)
	}

	return nil
}

func evolutionPropose(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: lvt evolution propose <event-id>")
	}
	eventID := args[0]

	collector := telemetry.NewCollector()
	defer collector.Close()

	if !collector.IsEnabled() {
		fmt.Println("Telemetry is disabled (set LVT_TELEMETRY=true to enable)")
		return nil
	}

	ctx := context.Background()
	event, err := collector.Store().Get(ctx, eventID)
	if err != nil {
		return fmt.Errorf("get event %s: %w", eventID, err)
	}

	patternsFile, err := knowledge.FindPatternsFile()
	if err != nil {
		return fmt.Errorf("find patterns file: %w", err)
	}
	kb, err := knowledge.New(patternsFile)
	if err != nil {
		return fmt.Errorf("load knowledge base: %w", err)
	}

	proposer := evolution.NewProposer(kb)
	proposal, err := proposer.ProposeFor(event)
	if err != nil {
		return fmt.Errorf("propose: %w", err)
	}

	fmt.Printf("Proposals for event %s\n", eventID)
	fmt.Printf("Command: %s\n", event.Command)
	fmt.Printf("Errors: %d\n", len(event.Errors))
	fmt.Println()

	if len(proposal.Fixes) == 0 {
		fmt.Println("No fixes proposed — no matching patterns in knowledge base.")
		return nil
	}

	fmt.Printf("Proposed Fixes (%d):\n", len(proposal.Fixes))
	fmt.Println(strings.Repeat("-", 80))

	for i, fix := range proposal.Fixes {
		fmt.Printf("\n  Fix %d: %s\n", i+1, fix.ID)
		fmt.Printf("  Pattern: %s\n", fix.PatternID)
		fmt.Printf("  Target:  %s\n", fix.TargetFile)

		loc := evolution.ClassifyFix(fix)
		switch loc.Type {
		case "component":
			fmt.Printf("  Location: component (%s)\n", loc.Component)
		case "kit":
			fmt.Printf("  Location: kit template\n")
		case "generated":
			fmt.Printf("  Location: generated code\n")
		default:
			fmt.Printf("  Location: %s\n", loc.Type)
		}

		fmt.Printf("  Confidence: %.0f%%\n", fix.Confidence*100)
		fmt.Printf("  Find:    %s\n", fix.FindPattern)
		fmt.Printf("  Replace: %s\n", fix.Replace)
	}

	return nil
}

func evolutionComponents(args []string) error {
	days := defaultLookbackDays
	for i := 0; i < len(args); i++ {
		if args[i] == "--days" && i+1 < len(args) {
			if n, err := strconv.Atoi(args[i+1]); err == nil && n > 0 {
				days = n
			}
			i++
		}
	}

	collector := telemetry.NewCollector()
	defer collector.Close()

	if !collector.IsEnabled() {
		fmt.Println("Telemetry is disabled (set LVT_TELEMETRY=true to enable)")
		return nil
	}

	ctx := context.Background()
	since := time.Now().AddDate(0, 0, -days)
	events, err := collector.Store().List(ctx, telemetry.ListOptions{Since: since})
	if err != nil {
		return fmt.Errorf("list events: %w", err)
	}

	// Aggregate by component
	type compStats struct {
		usage   int
		success int
		errors  map[string]int // error message → count
	}
	byComponent := make(map[string]*compStats)

	for _, e := range events {
		for _, comp := range e.ComponentsUsed {
			s, ok := byComponent[comp]
			if !ok {
				s = &compStats{errors: make(map[string]int)}
				byComponent[comp] = s
			}
			s.usage++
			if e.Success {
				s.success++
			}
		}
		for _, ce := range e.ComponentErrors {
			s, ok := byComponent[ce.Component]
			if !ok {
				s = &compStats{errors: make(map[string]int)}
				byComponent[ce.Component] = s
			}
			s.errors[truncate(ce.Message, 60)]++
		}
	}

	fmt.Printf("Component Health Dashboard (last %d days)\n", days)
	fmt.Println("==========================================")

	if len(byComponent) == 0 {
		fmt.Println("  No component data recorded yet.")
		fmt.Println("  Components are tracked in 'gen resource' commands.")
		return nil
	}

	fmt.Printf("%-20s %8s %8s %10s %s\n", "Component", "Uses", "Success", "Rate", "Status")
	fmt.Println(strings.Repeat("-", 65))

	// Sort component names for stable output
	names := make([]string, 0, len(byComponent))
	for name := range byComponent {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		s := byComponent[name]
		rate := float64(0)
		if s.usage > 0 {
			rate = float64(s.success) / float64(s.usage) * 100
		}
		status := "OK"
		if rate < 90 && s.usage > 0 {
			status = "WARN"
		}
		fmt.Printf("%-20s %8d %8d %9.1f%% %s\n", name, s.usage, s.success, rate, status)
	}

	// Print top errors per component
	hasErrors := false
	for _, name := range names {
		s := byComponent[name]
		if len(s.errors) == 0 {
			continue
		}
		if !hasErrors {
			fmt.Println()
			fmt.Println("Top Errors by Component")
			fmt.Println("-----------------------")
			hasErrors = true
		}
		fmt.Printf("\n  %s:\n", name)

		// Sort errors by count descending
		type errEntry struct {
			msg   string
			count int
		}
		errs := make([]errEntry, 0, len(s.errors))
		for msg, count := range s.errors {
			errs = append(errs, errEntry{msg, count})
		}
		sort.Slice(errs, func(i, j int) bool {
			return errs[i].count > errs[j].count
		})

		limit := 3
		if len(errs) < limit {
			limit = len(errs)
		}
		for _, e := range errs[:limit] {
			fmt.Printf("    [%dx] %s\n", e.count, e.msg)
		}
	}

	return nil
}

func evolutionApply(_ []string) error {
	return fmt.Errorf("apply is not yet implemented; use 'lvt evolution propose <event-id>' to see proposed fixes")
}

func evolutionUpstreamStatus(_ []string) error {
	patternsFile, err := knowledge.FindPatternsFile()
	if err != nil {
		return fmt.Errorf("find patterns file: %w", err)
	}

	kb, err := knowledge.New(patternsFile)
	if err != nil {
		return fmt.Errorf("load knowledge base: %w", err)
	}

	fmt.Println("Upstream Pattern Status")
	fmt.Println("=======================")

	found := false
	for _, p := range kb.ListPatterns() {
		if p.UpstreamRepo == "" {
			continue
		}
		found = true
		fmt.Printf("\n  Pattern:    %s\n", p.ID)
		fmt.Printf("  Name:       %s\n", p.Name)
		fmt.Printf("  Upstream:   %s\n", p.UpstreamRepo)
		fmt.Printf("  Confidence: %.0f%%\n", p.Confidence*100)
		fmt.Printf("  Fixes:      %d\n", len(p.Fixes))
	}

	if !found {
		fmt.Println("  No upstream patterns found.")
	}

	return nil
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-2]) + ".."
}
