package notification

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseNotification(t *testing.T) {
	t.Run("BomConsumed", func(t *testing.T) {
		notification := parseFromFile(t, "./testdata/bom-consumed.json")

		require.Equal(t, LevelInformational, notification.Level)
		require.Equal(t, ScopePortfolio, notification.Scope)
		require.Equal(t, GroupBOMConsumed, notification.Group)
		require.Equal(t, 2019, notification.Timestamp.Year())
		require.Equal(t, time.August, notification.Timestamp.Month())
		require.Equal(t, 23, notification.Timestamp.Day())
		require.Equal(t, 21, notification.Timestamp.Hour())
		require.Equal(t, 57, notification.Timestamp.Minute())
		require.Equal(t, 57, notification.Timestamp.Second())

		require.IsType(t, &BOMSubject{}, notification.Subject)
		subject := notification.Subject.(*BOMSubject)

		require.Equal(t, "6fb1820f-5280-4577-ac51-40124aabe307", subject.Project.UUID.String())
		require.Equal(t, "<base64 encoded bom>", subject.BOM.Content)
	})

	t.Run("BomProcessed", func(t *testing.T) {
		notification := parseFromFile(t, "./testdata/bom-processed.json")

		require.Equal(t, LevelInformational, notification.Level)
		require.Equal(t, ScopePortfolio, notification.Scope)
		require.Equal(t, GroupBOMProcessed, notification.Group)
		require.Equal(t, 2019, notification.Timestamp.Year())
		require.Equal(t, time.August, notification.Timestamp.Month())
		require.Equal(t, 23, notification.Timestamp.Day())
		require.Equal(t, 21, notification.Timestamp.Hour())
		require.Equal(t, 57, notification.Timestamp.Minute())
		require.Equal(t, 57, notification.Timestamp.Second())

		require.IsType(t, &BOMSubject{}, notification.Subject)
		subject := notification.Subject.(*BOMSubject)

		require.Equal(t, "6fb1820f-5280-4577-ac51-40124aabe307", subject.Project.UUID.String())
		require.Equal(t, "<base64 encoded bom>", subject.BOM.Content)
	})

	t.Run("NewVulnerableDependency", func(t *testing.T) {
		notification := parseFromFile(t, "./testdata/new-vulnerable-dependency.json")

		require.Equal(t, LevelInformational, notification.Level)
		require.Equal(t, ScopePortfolio, notification.Scope)
		require.Equal(t, GroupNewVulnerableDependency, notification.Group)
		require.Equal(t, 2018, notification.Timestamp.Year())
		require.Equal(t, time.August, notification.Timestamp.Month())
		require.Equal(t, 27, notification.Timestamp.Day())
		require.Equal(t, 23, notification.Timestamp.Hour())
		require.Equal(t, 23, notification.Timestamp.Minute())
		require.Equal(t, 0, notification.Timestamp.Second())

		require.IsType(t, &NewVulnerableDependencySubject{}, notification.Subject)
		subject := notification.Subject.(*NewVulnerableDependencySubject)

		require.Equal(t, "6fb1820f-5280-4577-ac51-40124aabe307", subject.Project.UUID.String())
		require.Equal(t, "4d5cd8df-cff7-4212-a038-91ae4ab79396", subject.Component.UUID.String())
		require.Len(t, subject.Vulnerabilities, 2)
		require.Equal(t, "941a93f5-e06b-4304-84de-4d788eeb4969", subject.Vulnerabilities[0].UUID.String())
		require.Equal(t, "ca318ca7-616f-4af0-9c6b-15b8e208c586", subject.Vulnerabilities[1].UUID.String())
	})

	t.Run("NewVulnerability", func(t *testing.T) {
		notification := parseFromFile(t, "./testdata/new-vulnerability.json")

		require.Equal(t, LevelInformational, notification.Level)
		require.Equal(t, ScopePortfolio, notification.Scope)
		require.Equal(t, GroupNewVulnerability, notification.Group)
		require.Equal(t, 2018, notification.Timestamp.Year())
		require.Equal(t, time.August, notification.Timestamp.Month())
		require.Equal(t, 27, notification.Timestamp.Day())
		require.Equal(t, 23, notification.Timestamp.Hour())
		require.Equal(t, 26, notification.Timestamp.Minute())
		require.Equal(t, 22, notification.Timestamp.Second())

		require.IsType(t, &NewVulnerabilitySubject{}, notification.Subject)
		subject := notification.Subject.(*NewVulnerabilitySubject)

		require.Equal(t, "4d5cd8df-cff7-4212-a038-91ae4ab79396", subject.Component.UUID.String())
		require.Equal(t, "941a93f5-e06b-4304-84de-4d788eeb4969", subject.Vulnerability.UUID.String())
		require.Len(t, subject.AffectedProjects, 1)
		require.Equal(t, "6fb1820f-5280-4577-ac51-40124aabe307", subject.AffectedProjects[0].UUID.String())
	})

	t.Run("PolicyViolation", func(t *testing.T) {
		notification := parseFromFile(t, "./testdata/policy-violation.json")

		require.Equal(t, LevelInformational, notification.Level)
		require.Equal(t, ScopePortfolio, notification.Scope)
		require.Equal(t, GroupPolicyViolation, notification.Group)
		require.Equal(t, 2022, notification.Timestamp.Year())
		require.Equal(t, time.May, notification.Timestamp.Month())
		require.Equal(t, 12, notification.Timestamp.Day())
		require.Equal(t, 23, notification.Timestamp.Hour())
		require.Equal(t, 7, notification.Timestamp.Minute())
		require.Equal(t, 59, notification.Timestamp.Second())

		require.IsType(t, &PolicyViolationSubject{}, notification.Subject)
		subject := notification.Subject.(*PolicyViolationSubject)

		require.Equal(t, "4e04c695-9acd-46fc-9bf6-ed23d7eb551e", subject.Component.UUID.String())
		require.Equal(t, "7a36e5c0-9f09-42dd-b401-360da56c2abe", subject.Project.UUID.String())
		require.Equal(t, "c82fcb50-029a-4636-a657-96242b20680e", subject.PolicyViolation.UUID.String())
		require.Equal(t, "8e5c0a5b-71fb-45c5-afac-6c6a99742cbe", subject.PolicyViolation.PolicyCondition.UUID.String())
		require.Equal(t, "6d4c7398-689a-4ec7-b5c5-9abb6b5393e9", subject.PolicyViolation.PolicyCondition.Policy.UUID.String())
	})
}

func parseFromFile(t *testing.T, filePath string) (n Notification) {
	t.Helper()
	file, err := os.Open(filePath)
	require.NoError(t, err)
	defer file.Close()

	n, err = Parse(file)
	require.NoError(t, err)

	return
}
