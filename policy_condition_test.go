package dtrack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPolicyConditionMarshalViolationType(t *testing.T) {
	condition := PolicyCondition{
		Operator:      PolicyConditionOperatorMatches,
		Subject:       PolicyConditionSubjectExpression,
		Value:         `component.is_internal == false`,
		ViolationType: PolicyConditionViolationTypeLicense,
	}

	encoded, err := json.Marshal(condition)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.Equal(t, "EXPRESSION", decoded["subject"])
	require.Equal(t, "LICENSE", decoded["violationType"])
}

func TestPolicyConditionMarshalOmitsEmptyViolationType(t *testing.T) {
	condition := PolicyCondition{
		Operator: PolicyConditionOperatorIs,
		Subject:  PolicyConditionSubjectSeverity,
		Value:    "CRITICAL",
	}

	encoded, err := json.Marshal(condition)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.NotContains(t, decoded, "violationType")
}

func TestPolicyConditionUnmarshalViolationType(t *testing.T) {
	const body = `{
		"uuid": "81b6f3d1-d9ae-49d2-b687-fd1c7b783736",
		"operator": "MATCHES",
		"subject": "EXPRESSION",
		"value": "component.is_internal == false",
		"violationType": "OPERATIONAL"
	}`

	var condition PolicyCondition
	require.NoError(t, json.Unmarshal([]byte(body), &condition))
	require.Equal(t, PolicyConditionSubjectExpression, condition.Subject)
	require.Equal(t, PolicyConditionViolationTypeOperational, condition.ViolationType)
}
