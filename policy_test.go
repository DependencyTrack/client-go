package dtrack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func b(v bool) *bool { return &v }

func TestPolicyJSONEncoding_Bools(t *testing.T) {
	cases := []struct {
		name string
		in   Policy
		want string
	}{
		{
			name: "omit both",
			in:   Policy{Name: "x"},
			want: `{"uuid":"00000000-0000-0000-0000-000000000000","name":"x","operator":"","violationState":""}`,
		},
		{
			name: "send false global",
			in:   Policy{Name: "x", Global: b(false)},
			want: `{"uuid":"00000000-0000-0000-0000-000000000000","name":"x","operator":"","violationState":"","global":false}`,
		},
		{
			name: "send true global",
			in:   Policy{Name: "x", Global: b(true)},
			want: `{"uuid":"00000000-0000-0000-0000-000000000000","name":"x","operator":"","violationState":"","global":true}`,
		},
		{
			name: "send false includeChildren",
			in:   Policy{Name: "x", IncludeChildren: b(false)},
			want: `{"uuid":"00000000-0000-0000-0000-000000000000","name":"x","operator":"","violationState":"","includeChildren":false}`,
		},
		{
			name: "send true includeChildren",
			in:   Policy{Name: "x", IncludeChildren: b(true)},
			want: `{"uuid":"00000000-0000-0000-0000-000000000000","name":"x","operator":"","violationState":"","includeChildren":true}`,
		},
		{
			name: "both flags",
			in:   Policy{Name: "x", Global: b(false), IncludeChildren: b(true)},
			want: `{"uuid":"00000000-0000-0000-0000-000000000000","name":"x","operator":"","violationState":"","includeChildren":true,"global":false}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal the Policy struct
			got, err := json.Marshal(tc.in)
			require.NoError(t, err)

			// Compare by unmarshaling both to ensure key order insensitivity
			var gotMap, wantMap map[string]interface{}
			err = json.Unmarshal(got, &gotMap)
			require.NoError(t, err)
			err = json.Unmarshal([]byte(tc.want), &wantMap)
			require.NoError(t, err)

			require.Equal(t, wantMap, gotMap)
		})
	}
}

func TestPolicyJSONDecoding_Bools(t *testing.T) {
	cases := []struct {
		name           string
		json           string
		wantGlobal     *bool
		wantIncludeCh  *bool
	}{
		{
			name:           "omit both",
			json:           `{"name":"x","operator":"","violationState":""}`,
			wantGlobal:     nil,
			wantIncludeCh:  nil,
		},
		{
			name:           "receive false global",
			json:           `{"name":"x","operator":"","violationState":"","global":false}`,
			wantGlobal:     b(false),
			wantIncludeCh:  nil,
		},
		{
			name:           "receive true global",
			json:           `{"name":"x","operator":"","violationState":"","global":true}`,
			wantGlobal:     b(true),
			wantIncludeCh:  nil,
		},
		{
			name:           "receive false includeChildren",
			json:           `{"name":"x","operator":"","violationState":"","includeChildren":false}`,
			wantGlobal:     nil,
			wantIncludeCh:  b(false),
		},
		{
			name:           "receive true includeChildren",
			json:           `{"name":"x","operator":"","violationState":"","includeChildren":true}`,
			wantGlobal:     nil,
			wantIncludeCh:  b(true),
		},
		{
			name:           "both flags",
			json:           `{"name":"x","operator":"","violationState":"","includeChildren":true,"global":false}`,
			wantGlobal:     b(false),
			wantIncludeCh:  b(true),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var p Policy
			err := json.Unmarshal([]byte(tc.json), &p)
			require.NoError(t, err)

			// Check Global field
			if tc.wantGlobal == nil {
				require.Nil(t, p.Global)
			} else {
				require.NotNil(t, p.Global)
				require.Equal(t, *tc.wantGlobal, *p.Global)
			}

			// Check IncludeChildren field
			if tc.wantIncludeCh == nil {
				require.Nil(t, p.IncludeChildren)
			} else {
				require.NotNil(t, p.IncludeChildren)
				require.Equal(t, *tc.wantIncludeCh, *p.IncludeChildren)
			}
		})
	}
}

func TestPolicyJSONEncoding_BoolsWithOptionalBoolOf(t *testing.T) {
	cases := []struct {
		name string
		in   Policy
		want string
	}{
		{
			name: "OptionalBoolOf true",
			in:   Policy{Name: "x", Global: OptionalBoolOf(true)},
			want: `{"uuid":"00000000-0000-0000-0000-000000000000","name":"x","operator":"","violationState":"","global":true}`,
		},
		{
			name: "OptionalBoolOf false",
			in:   Policy{Name: "x", Global: OptionalBoolOf(false)},
			want: `{"uuid":"00000000-0000-0000-0000-000000000000","name":"x","operator":"","violationState":"","global":false}`,
		},
		{
			name: "both with OptionalBoolOf",
			in:   Policy{Name: "x", Global: OptionalBoolOf(true), IncludeChildren: OptionalBoolOf(false)},
			want: `{"uuid":"00000000-0000-0000-0000-000000000000","name":"x","operator":"","violationState":"","includeChildren":false,"global":true}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal the Policy struct
			got, err := json.Marshal(tc.in)
			require.NoError(t, err)

			// Compare by unmarshaling both to ensure key order insensitivity
			var gotMap, wantMap map[string]interface{}
			err = json.Unmarshal(got, &gotMap)
			require.NoError(t, err)
			err = json.Unmarshal([]byte(tc.want), &wantMap)
			require.NoError(t, err)

			require.Equal(t, wantMap, gotMap)
		})
	}
}

func TestPolicyJSONRoundTrip_Bools(t *testing.T) {
	cases := []struct {
		name string
		in   Policy
	}{
		{
			name: "nil both",
			in:   Policy{Name: "x", Operator: PolicyOperatorAll, ViolationState: PolicyViolationStateInfo},
		},
		{
			name: "false global",
			in:   Policy{Name: "x", Operator: PolicyOperatorAll, ViolationState: PolicyViolationStateInfo, Global: b(false)},
		},
		{
			name: "true global",
			in:   Policy{Name: "x", Operator: PolicyOperatorAll, ViolationState: PolicyViolationStateInfo, Global: b(true)},
		},
		{
			name: "false includeChildren",
			in:   Policy{Name: "x", Operator: PolicyOperatorAny, ViolationState: PolicyViolationStateWarn, IncludeChildren: b(false)},
		},
		{
			name: "true includeChildren",
			in:   Policy{Name: "x", Operator: PolicyOperatorAny, ViolationState: PolicyViolationStateWarn, IncludeChildren: b(true)},
		},
		{
			name: "both flags set",
			in:   Policy{Name: "x", Operator: PolicyOperatorAll, ViolationState: PolicyViolationStateFail, Global: b(true), IncludeChildren: b(false)},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal to JSON
			jsonBytes, err := json.Marshal(tc.in)
			require.NoError(t, err)

			// Unmarshal back to struct
			var out Policy
			err = json.Unmarshal(jsonBytes, &out)
			require.NoError(t, err)

			// Compare relevant fields
			require.Equal(t, tc.in.Name, out.Name)
			require.Equal(t, tc.in.Operator, out.Operator)
			require.Equal(t, tc.in.ViolationState, out.ViolationState)

			// Check Global field
			if tc.in.Global == nil {
				require.Nil(t, out.Global)
			} else {
				require.NotNil(t, out.Global)
				require.Equal(t, *tc.in.Global, *out.Global)
			}

			// Check IncludeChildren field
			if tc.in.IncludeChildren == nil {
				require.Nil(t, out.IncludeChildren)
			} else {
				require.NotNil(t, out.IncludeChildren)
				require.Equal(t, *tc.in.IncludeChildren, *out.IncludeChildren)
			}
		})
	}
}
