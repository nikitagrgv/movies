package config

import (
	"fmt"
	"strings"
)

type stubType string

const (
	SearchStub stubType = "search"
	GetStub    stubType = "get"
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
		case SearchStub, GetStub:
			result = append(result, stubType(p))
		default:
			return nil, fmt.Errorf("invalid stub type: %s", p)
		}
	}

	return result, nil
}
