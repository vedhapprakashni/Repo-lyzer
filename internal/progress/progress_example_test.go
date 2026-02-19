package progress

import (
	"testing"
	"time"
)

// TestOverallProgress demonstrates the overall progress tracking functionality
func TestOverallProgress(t *testing.T) {
	// Create a new progress tracker with 5 steps
	progress := NewOverallProgress(5)

	// Simulate different steps
	steps := []string{
		"Fetching repository information",
		"Analyzing commits",
		"Processing languages",
		"Calculating metrics",
		"Generating report",
	}

	for i, step := range steps {
		progress.StartStep(step)

		// Simulate work
		time.Sleep(100 * time.Millisecond)

		progress.CompleteStep(step)

		// Verify progress
		completed, total, percentage := progress.GetProgress()
		if completed != i+1 || total != 5 {
			t.Errorf("Step %d: expected %d/%d, got %d/%d", i, i+1, 5, completed, total)
		}

		expectedPercentage := ((i + 1) * 100) / 5
		if percentage != expectedPercentage {
			t.Errorf("Step %d: expected %d%% progress, got %d%%", i, expectedPercentage, percentage)
		}
	}

	progress.Finish()

	// Verify completion
	completed, total, percentage := progress.GetProgress()
	if completed != 5 || total != 5 || percentage != 100 {
		t.Errorf("Expected 5/5 (100%%), got %d/%d (%d%%)", completed, total, percentage)
	}
}

// TestOverallProgressUpdate tests step message updates
func TestOverallProgressUpdate(t *testing.T) {
	progress := NewOverallProgress(3)

	progress.StartStep("Step 1: Starting")
	progress.UpdateStep("Step 1: Processing items (25%)")
	progress.UpdateStep("Step 1: Processing items (50%)")
	progress.UpdateStep("Step 1: Processing items (75%)")
	progress.CompleteStep("Step 1: Complete")

	progress.StartStep("Step 2: Starting")
	progress.CompleteStep("Step 2: Complete")

	progress.StartStep("Step 3: Starting")
	progress.CompleteStep("Step 3: Complete")

	progress.Finish()

	completed, _, percentage := progress.GetProgress()
	if completed != 3 || percentage != 100 {
		t.Errorf("Expected 100%% progress, got %d%%", percentage)
	}
}
