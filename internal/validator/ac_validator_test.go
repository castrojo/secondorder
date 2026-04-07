package validator

import (
	"reflect"
	"testing"
)

func TestValidateAC(t *testing.T) {
	tests := []struct {
		name         string
		issueType    string
		description  string
		wantWarnings []string
	}{
		{
			name:        "api: missing AC section",
			issueType:   "api",
			description: "Some description without AC",
			wantWarnings: []string{
				"Issue is missing an 'Acceptance Criteria' or 'AC' section.",
			},
		},
		{
			name:        "backend: missing AC section",
			issueType:   "backend",
			description: "Some description without AC",
			wantWarnings: []string{
				"Issue is missing an 'Acceptance Criteria' or 'AC' section.",
			},
		},
		{
			name:         "task: missing AC section (no warnings expected)",
			issueType:    "task",
			description:  "Some description without AC",
			wantWarnings: nil,
		},
		{
			name:      "api: incomplete AC",
			issueType: "api",
			description: `## AC
- some stuff
`,
			wantWarnings: []string{
				"Type 'api' usually requires an endpoint path and method.",
				"Type 'api' should specify a request schema.",
				"Type 'api' should specify a response schema.",
				"Type 'api' should specify expected status codes.",
			},
		},
		{
			name:      "api: complete AC",
			issueType: "api",
			description: `## Acceptance Criteria
- path: /api/v1/test
- method: GET
- request body: none
- response schema: { "ok": true }
- status code: 200
`,
			wantWarnings: nil,
		},
		{
			name:      "backend: incomplete AC",
			issueType: "backend",
			description: `## AC
- just some text
`,
			wantWarnings: []string{
				"Type 'backend' should describe core business logic or rules.",
				"Type 'backend' should specify database or persistence changes.",
				"Type 'backend' should list external or internal dependencies.",
			},
		},
		{
			name:      "backend: complete AC",
			issueType: "backend",
			description: `## AC
- logic: implement user auth
- database: update users table
- service: depends on auth-service
`,
			wantWarnings: nil,
		},
		{
			name:      "feature: missing list in AC",
			issueType: "feature",
			description: `## AC
This is just a paragraph without any list items.
`,
			wantWarnings: []string{
				"Acceptance criteria should include at least one bullet point or numbered item.",
			},
		},
		{
			name:      "feature: valid AC with list",
			issueType: "feature",
			description: `## AC
- User can login
- User can logout
`,
			wantWarnings: nil,
		},
		{
			name:      "feature: valid AC with numbered list",
			issueType: "feature",
			description: `## Acceptance Criteria
1. First step
2. Second step
`,
			wantWarnings: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWarnings := ValidateAC(tt.issueType, tt.description)
			if !reflect.DeepEqual(gotWarnings, tt.wantWarnings) {
				t.Errorf("ValidateAC() = %v, want %v", gotWarnings, tt.wantWarnings)
			}
		})
	}
}
