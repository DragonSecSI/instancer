package strings

import "fmt"

func fuzzyTime(input int) string {
	suffix := ""
	if input < 0 {
		suffix = " ago"
		input = -input
	}

	if input == 0 {
		return "just now" + suffix
	}
	if input < 60 {
		return fmt.Sprintf("%d second%s%s", input, fuzzyTimePlural(input), suffix)
	}
	if input < 3600 {
		minutes := input / 60
		return fmt.Sprintf("%d minute%s%s", minutes, fuzzyTimePlural(minutes), suffix)
	}
	if input < 86400 {
		hours := input / 3600
		return fmt.Sprintf("%d hours%s%s", hours, fuzzyTimePlural(hours), suffix)
	}

	return fmt.Sprintf("%d days%s", input/86400, suffix)
}

func fuzzyTimePlural(input int) string {
	if input == 1 {
		return ""
	} else {
		return "s"
	}
}
