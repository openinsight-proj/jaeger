package tracestore

import (
	"regexp"
	"strings"

	"github.com/jaegertracing/jaeger/internal/storage/elasticsearch/dbmodel"
)

const unrouted = "unrouted"

func parseDynamicKeys(template string) []string {
	re := regexp.MustCompile(`\{\{(.*?)\}\}`)
	matches := re.FindAllStringSubmatch(template, -1)

	keyMap := make(map[string]struct{})
	for _, match := range matches {
		if len(match) > 1 {
			keyMap[strings.TrimSpace(match[1])] = struct{}{}
		}
	}
	var keys []string
	for k := range keyMap {
		keys = append(keys, k)
	}

	return keys
}

func buildIndexSuffix(span *dbmodel.Span, dynamicKeys []string, template string) string {
	if len(dynamicKeys) == 0 {
		return template
	}

	data := make(map[string]string, len(dynamicKeys))
	for _, target := range dynamicKeys {
		value := unrouted
		for _, tag := range span.Process.Tags {
			if tag.Key == target && tag.Type == dbmodel.StringType {
				value = tag.Value.(string)
				break
			}
		}

		if value == unrouted {
			for _, tag := range span.Tags {
				if tag.Key == target && tag.Type == dbmodel.StringType {
					value = tag.Value.(string)
					break
				}
			}
		}

		data[target] = value
	}

	return applyTemplate(template, data)
}

func applyTemplate(template string, data map[string]string) string {
	for k, v := range data {
		template = regexp.MustCompile(`\{\{`+regexp.QuoteMeta(k)+`\}\}`).ReplaceAllString(template, v)
	}

	return template
}
