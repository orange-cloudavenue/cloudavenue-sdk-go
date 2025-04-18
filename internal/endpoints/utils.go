package endpoints

import "strings"

func InlineTemplate(template string, data map[string]string) string {
	for key, value := range data {
		template = strings.ReplaceAll(template, "{"+key+"}", value)
	}
	return template
}
