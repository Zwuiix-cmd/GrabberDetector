package Injection

import (
	"regexp"
)

type OtherGrabber struct{}

func (d OtherGrabber) Flag(content string) ([]string, []string) {
	webhookLinkPattern := `https?://[^\s"'` + `"]*webhook[^\s"'` + `"]*`
	linkPattern := `https?://[^\s'"]+`
	reWebhook := regexp.MustCompile(webhookLinkPattern)
	webhookLinks := reWebhook.FindAllString(content, -1)
	reAllLinks := regexp.MustCompile(linkPattern)
	allLinks := reAllLinks.FindAllString(content, -1)
	allLinks = append(webhookLinks, allLinks...)

	return webhookLinks, allLinks
}
