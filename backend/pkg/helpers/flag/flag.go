package flag

import (
	"math/rand"
	"strings"

	"github.com/DragonSecSI/instancer/backend/pkg/database/models"
)

type Flag struct {
	Process    func(string, models.ChallengeFlagType) string
	Suffix     func(string) string
	Leetify    func(string) string
	Capitalize func(string) string
}

func NewFlagHelper() Flag {
	return Flag{
		Process:    process,
		Suffix:     addSuffix,
		Leetify:    leetify,
		Capitalize: capitalize,
	}
}

func process(flag string, flagType models.ChallengeFlagType) string {
	if flagType == models.ChallengeFlagTypeStatic {
		return flag
	}

	if flagType&models.ChallengeFlagTypeCapitalize != 0 {
		flag = capitalize(flag)
	}
	if flagType&models.ChallengeFlagTypeLeetify != 0 {
		flag = leetify(flag)
	}
	if flagType&models.ChallengeFlagTypeSuffix != 0 {
		flag = addSuffix(flag)
	}

	return flag
}

func addSuffix(input string) string {
	randsuffix := "_" + randomString(6)

	if strings.HasSuffix(input, "}") {
		input = strings.TrimSuffix(input, "}")
		input += randsuffix + "}"
	} else {
		input += randsuffix
	}

	return input
}

func leetify(input string) string {
	replacements := map[rune]rune{
		'a': '4',
		'e': '3',
		'i': '1',
		'o': '0',
		's': '5',
		't': '7',
		'g': '9',
		'A': '4',
		'E': '3',
		'I': '1',
		'O': '0',
		'S': '5',
		'T': '7',
		'G': '9',
	}

	bracket := strings.Index(input, "{")
	if bracket == -1 {
		bracket = 0
	}

	var result strings.Builder
	for i, char := range input {
		if i <= bracket {
			result.WriteRune(char)
			continue
		}
		if replacement, exists := replacements[char]; exists {
			if rand.Float32() < 0.4 {
				result.WriteRune(replacement)
			} else {
				result.WriteRune(char)
			}
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

func capitalize(input string) string {
	bracket := strings.Index(input, "{")
	if bracket == -1 {
		bracket = 0
	}

	var result strings.Builder
	for i, char := range input {
		if i <= bracket {
			result.WriteRune(char)
			continue
		}

		if rand.Float32() < 0.4 {
			if char >= 'a' && char <= 'z' {
				result.WriteRune(char - 32) // Convert to uppercase
			} else if char >= 'A' && char <= 'Z' {
				result.WriteRune(char + 32) // Convert to lowercase
			} else {
				result.WriteRune(char)
			}
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var result strings.Builder
	for range length {
		result.WriteByte(charset[rand.Intn(len(charset))])
	}

	return result.String()
}
