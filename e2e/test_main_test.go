//go:build browser

package e2e

import e2etest "github.com/livetemplate/lvt/testing"

func cleanupChromeContainers() {
	e2etest.CleanupChromeContainers()
}

// cleanupAllTestContainers removes both Chrome and app test containers.
// This should be called at the start of tests to ensure a clean environment.
func cleanupAllTestContainers() {
	e2etest.CleanupAllTestContainers()
}
