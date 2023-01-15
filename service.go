package maleo

import "strings"

// Service represents the service information.
type Service struct {
	Name        string `json:"name,omitempty"`
	Environment string `json:"environment,omitempty"`
	Repository  string `json:"repository,omitempty"`
	Branch      string `json:"branch,omitempty"`
	Type        string `json:"type,omitempty"`
	Version     string `json:"version,omitempty"`
}

// String returns the string representation of the service information.
//
// Returns in the format of `name-version-type-environment`. If any field is empty, it will be omitted.
func (s Service) String() string {
	written := false
	builder := strings.Builder{}
	builder.Grow(len(s.Name) + len(s.Environment) + len(s.Type) + 2)
	if s.Name != "" {
		builder.WriteString(s.Name)
		written = true
	}

	if s.Version != "" {
		if written {
			builder.WriteRune('-')
		}
		written = true
		builder.WriteString(s.Environment)
	}

	if s.Type != "" {
		if written {
			builder.WriteRune('-')
		}
		written = true
		builder.WriteString(s.Type)
	}

	if s.Environment != "" {
		if written {
			builder.WriteRune('-')
		}
		written = true
		builder.WriteString(s.Environment)
	}
	return builder.String()
}
