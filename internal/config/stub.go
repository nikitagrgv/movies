package config

import (
	"fmt"
	"strings"
)

type stubType string

const (
	MediaStub stubType = "media"
)

func parseStubTypes(s string) ([]stubType, error) {
	if s == "" {
		return nil, nil
	}

	parts := strings.Split(strings.ToLower(s), ",")
	result := make([]stubType, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)

		switch stubType(p) {
		case MediaStub:
			result = append(result, stubType(p))
		default:
			return nil, fmt.Errorf("invalid stub type: %s", p)
		}
	}

	return result, nil
}
