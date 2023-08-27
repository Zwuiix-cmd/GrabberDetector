package Injection

import (
	"regexp"
)

type Donerium struct{}

func (d Donerium) Flag(content string) (string, bool) {
	keywordFind := false

	doeneriumPattern := regexp.MustCompile(`doenerium`)
	if doeneriumPattern.MatchString(content) {
		keywordFind = true
	}

	webhookPattern := regexp.MustCompile(`webhook: "([^"]+)"`)
	match := webhookPattern.FindStringSubmatch(content)
	if len(match) > 1 {
		return match[1], keywordFind
	}

	return "", keywordFind
}
