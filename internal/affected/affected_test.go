package affected

import (
	"testing"

	"github.com/amido/mrbuild/internal/config"
	"github.com/amido/mrbuild/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// This test file contains comprehensive tests for the project ordering functionality.
// The tests cover:
// - Positive order values (ascending sort)
// - Negative order values (should come before positive)
// - Duplicate order values (stable sort maintains original order)
// - Default order value (0)
// - Extreme values (int32 min/max)
// - Partial project matches (only affected projects included)
// - Order field preservation from config to spawn
// - Sort stability verification

// TestGetProjectsOrdering tests that projects are sorted correctly by their Order field
func TestGetProjectsOrdering(t *testing.T) {
	// Create a logger for testing
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Suppress logs during tests

	tests := []struct {
		name          string
		projects      []config.Project
		changedFiles  string
		expectedOrder []string
		description   string
	}{
		{
			name: "Ascending order - positive numbers",
			projects: []config.Project{
				{
					Name:     "project-high",
					Folder:   "src/high",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo high"},
					Order:    10,
				},
				{
					Name:     "project-low",
					Folder:   "src/low",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo low"},
					Order:    1,
				},
				{
					Name:     "project-medium",
					Folder:   "src/medium",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo medium"},
					Order:    5,
				},
			},
			changedFiles:  "src/high/main.go\nsrc/low/main.go\nsrc/medium/main.go",
			expectedOrder: []string{"project-low", "project-medium", "project-high"},
			description:   "Projects with positive order values should be sorted in ascending order",
		},
		{
			name: "Negative and positive order numbers",
			projects: []config.Project{
				{
					Name:     "project-positive",
					Folder:   "src/positive",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo positive"},
					Order:    5,
				},
				{
					Name:     "project-negative",
					Folder:   "src/negative",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo negative"},
					Order:    -10,
				},
				{
					Name:     "project-zero",
					Folder:   "src/zero",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo zero"},
					Order:    0,
				},
			},
			changedFiles:  "src/positive/main.go\nsrc/negative/main.go\nsrc/zero/main.go",
			expectedOrder: []string{"project-negative", "project-zero", "project-positive"},
			description:   "Negative order values should come before zero and positive values",
		},
		{
			name: "Duplicate order numbers",
			projects: []config.Project{
				{
					Name:     "project-first",
					Folder:   "src/first",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo first"},
					Order:    5,
				},
				{
					Name:     "project-second",
					Folder:   "src/second",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo second"},
					Order:    5,
				},
				{
					Name:     "project-third",
					Folder:   "src/third",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo third"},
					Order:    5,
				},
			},
			changedFiles:  "src/first/main.go\nsrc/second/main.go\nsrc/third/main.go",
			expectedOrder: []string{"project-first", "project-second", "project-third"},
			description:   "Projects with duplicate order values maintain their original relative order (stable sort)",
		},
		{
			name: "Default order (zero) mixed with explicit values",
			projects: []config.Project{
				{
					Name:     "project-explicit-high",
					Folder:   "src/explicit-high",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo explicit-high"},
					Order:    10,
				},
				{
					Name:     "project-default",
					Folder:   "src/default",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo default"},
					Order:    0, // default value
				},
				{
					Name:     "project-explicit-low",
					Folder:   "src/explicit-low",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo explicit-low"},
					Order:    -5,
				},
			},
			changedFiles:  "src/explicit-high/main.go\nsrc/default/main.go\nsrc/explicit-low/main.go",
			expectedOrder: []string{"project-explicit-low", "project-default", "project-explicit-high"},
			description:   "Default order (0) should be sorted between negative and positive values",
		},
		{
			name: "Large order numbers",
			projects: []config.Project{
				{
					Name:     "project-max",
					Folder:   "src/max",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo max"},
					Order:    2147483647, // max int32
				},
				{
					Name:     "project-min",
					Folder:   "src/min",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo min"},
					Order:    -2147483648, // min int32
				},
				{
					Name:     "project-mid",
					Folder:   "src/mid",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo mid"},
					Order:    0,
				},
			},
			changedFiles:  "src/max/main.go\nsrc/min/main.go\nsrc/mid/main.go",
			expectedOrder: []string{"project-min", "project-mid", "project-max"},
			description:   "Extremely large and small order values should be handled correctly",
		},
		{
			name: "Only some projects affected",
			projects: []config.Project{
				{
					Name:     "project-a",
					Folder:   "src/a",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo a"},
					Order:    3,
				},
				{
					Name:     "project-b",
					Folder:   "src/b",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo b"},
					Order:    1,
				},
				{
					Name:     "project-c",
					Folder:   "src/c",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo c"},
					Order:    2,
				},
			},
			changedFiles:  "src/a/main.go\nsrc/c/main.go",
			expectedOrder: []string{"project-c", "project-a"},
			description:   "Only affected projects should be included and sorted",
		},
		{
			name: "Empty spawns list",
			projects: []config.Project{
				{
					Name:     "project-unaffected",
					Folder:   "src/unaffected",
					Patterns: []string{".*\\.go"},
					Build:    config.Build{Cmd: "echo unaffected"},
					Order:    1,
				},
			},
			changedFiles:  "other/file.txt",
			expectedOrder: []string{},
			description:   "No projects affected should return empty list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the configuration
			cfg := &config.Config{
				Input: config.InputConfig{
					Projects: tt.projects,
					Options:  config.Options{},
				},
			}

			// Create the App
			app := &models.App{
				Logger: logger,
			}

			// Create the Affected instance
			affected := New(app, cfg, logger)

			// Call getProjects
			spawns := affected.getProjects(tt.changedFiles)

			// Verify the number of spawns
			assert.Equal(t, len(tt.expectedOrder), len(spawns),
				"Number of affected projects should match expected count")

			// Verify the order
			actualOrder := make([]string, len(spawns))
			for i, spawn := range spawns {
				actualOrder[i] = spawn.Name
			}

			assert.Equal(t, tt.expectedOrder, actualOrder, tt.description)

			// Additional verification: ensure spawns are in ascending order by Order field
			for i := 1; i < len(spawns); i++ {
				assert.LessOrEqual(t, spawns[i-1].Order, spawns[i].Order,
					"Spawns should be sorted in ascending order by Order field")
			}
		})
	}
}

// TestGetProjectsOrderFieldPreservation tests that the Order field is correctly preserved
func TestGetProjectsOrderFieldPreservation(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	projects := []config.Project{
		{
			Name:     "project-1",
			Folder:   "src/p1",
			Patterns: []string{".*\\.go"},
			Build:    config.Build{Cmd: "echo p1"},
			Order:    42,
		},
		{
			Name:     "project-2",
			Folder:   "src/p2",
			Patterns: []string{".*\\.go"},
			Build:    config.Build{Cmd: "echo p2"},
			Order:    -7,
		},
	}

	cfg := &config.Config{
		Input: config.InputConfig{
			Projects: projects,
		},
	}

	app := &models.App{
		Logger: logger,
	}

	affected := New(app, cfg, logger)

	changedFiles := "src/p1/main.go\nsrc/p2/main.go"
	spawns := affected.getProjects(changedFiles)

	assert.Equal(t, 2, len(spawns), "Should have 2 spawns")

	// Find each spawn and verify Order is preserved
	for _, spawn := range spawns {
		for _, project := range projects {
			if spawn.Name == project.Name {
				assert.Equal(t, project.Order, spawn.Order,
					"Order field should be preserved from project to spawn for %s", spawn.Name)
			}
		}
	}
}

// TestGetProjectsSortStability tests that sort is stable (maintains relative order for equal values)
func TestGetProjectsSortStability(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Create projects with same order value in a specific sequence
	projects := []config.Project{
		{Name: "alpha", Folder: "src/alpha", Patterns: []string{".*"}, Build: config.Build{Cmd: "echo"}, Order: 1},
		{Name: "beta", Folder: "src/beta", Patterns: []string{".*"}, Build: config.Build{Cmd: "echo"}, Order: 1},
		{Name: "gamma", Folder: "src/gamma", Patterns: []string{".*"}, Build: config.Build{Cmd: "echo"}, Order: 1},
		{Name: "delta", Folder: "src/delta", Patterns: []string{".*"}, Build: config.Build{Cmd: "echo"}, Order: 1},
	}

	cfg := &config.Config{
		Input: config.InputConfig{
			Projects: projects,
		},
	}

	app := &models.App{
		Logger: logger,
	}

	affected := New(app, cfg, logger)

	// All projects affected
	changedFiles := "src/alpha/f\nsrc/beta/f\nsrc/gamma/f\nsrc/delta/f"
	spawns := affected.getProjects(changedFiles)

	// The order should be preserved as Go's sort.Slice is stable
	expectedOrder := []string{"alpha", "beta", "gamma", "delta"}
	actualOrder := make([]string, len(spawns))
	for i, spawn := range spawns {
		actualOrder[i] = spawn.Name
	}

	assert.Equal(t, expectedOrder, actualOrder,
		"Sort should be stable - projects with same order value should maintain their original relative order")
}
