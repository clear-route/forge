package coding

import (
	"fmt"
	"strings"
)

// GenerateUnifiedDiff creates a unified diff between original and modified content.
func GenerateUnifiedDiff(original, modified, filename string) string {
	originalLines := strings.Split(original, "\n")
	modifiedLines := strings.Split(modified, "\n")

	var diff strings.Builder
	diff.WriteString(fmt.Sprintf("--- %s\n", filename))
	diff.WriteString(fmt.Sprintf("+++ %s\n", filename))

	changes := findChanges(originalLines, modifiedLines)
	if len(changes) == 0 {
		return "No changes"
	}

	for _, change := range changes {
		diff.WriteString(fmt.Sprintf("@@ -%d,%d +%d,%d @@\n",
			change.originalStart+1, change.originalCount,
			change.modifiedStart+1, change.modifiedCount))

		for _, line := range change.lines {
			diff.WriteString(line)
			diff.WriteString("\n")
		}
	}

	return diff.String()
}

type diffChange struct {
	originalStart int
	originalCount int
	modifiedStart int
	modifiedCount int
	lines         []string
}

func findChanges(original, modified []string) []diffChange {
	var changes []diffChange
	var currentChange *diffChange

	maxLen := len(original)
	if len(modified) > maxLen {
		maxLen = len(modified)
	}

	for i := 0; i < maxLen; i++ {
		origLine := ""
		modLine := ""

		if i < len(original) {
			origLine = original[i]
		}
		if i < len(modified) {
			modLine = modified[i]
		}

		if origLine != modLine {
			if currentChange == nil {
				currentChange = &diffChange{
					originalStart: i,
					modifiedStart: i,
					lines:         []string{},
				}
			}

			if i < len(original) {
				currentChange.lines = append(currentChange.lines, "-"+origLine)
				currentChange.originalCount++
			}

			if i < len(modified) {
				currentChange.lines = append(currentChange.lines, "+"+modLine)
				currentChange.modifiedCount++
			}
		} else {
			if currentChange != nil {
				changes = append(changes, *currentChange)
				currentChange = nil
			}
		}
	}

	if currentChange != nil {
		changes = append(changes, *currentChange)
	}

	return changes
}
