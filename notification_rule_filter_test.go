package dtrack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotificationRuleMarshalFilterExpression(t *testing.T) {
	rule := NotificationRule{
		Name:             "Policy violations",
		Scope:            NotificationRuleScopePortfolio,
		FilterExpression: `subject.policy_violation.condition.policy.violation_state == "FAIL"`,
	}

	encoded, err := json.Marshal(rule)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.Equal(t, `subject.policy_violation.condition.policy.violation_state == "FAIL"`, decoded["filterExpression"])
}

func TestNotificationRuleMarshalOmitsEmptyFilterExpression(t *testing.T) {
	rule := NotificationRule{
		Name:  "Policy violations",
		Scope: NotificationRuleScopePortfolio,
	}

	encoded, err := json.Marshal(rule)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.NotContains(t, decoded, "filterExpression")
}

func TestNotificationRuleUnmarshalFilterExpression(t *testing.T) {
	const body = `{
		"uuid": "fdcff7ae-a0c0-4f54-9bb6-7bdb6c56c9fd",
		"name": "Policy violations",
		"scope": "PORTFOLIO",
		"filterExpression": "group != Group.GROUP_POLICY_VIOLATION"
	}`

	var rule NotificationRule
	require.NoError(t, json.Unmarshal([]byte(body), &rule))
	require.Equal(t, "group != Group.GROUP_POLICY_VIOLATION", rule.FilterExpression)
}
