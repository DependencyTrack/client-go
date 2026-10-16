package dtrack

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

type NotificationService struct {
	client *Client
}

type NotificationPublisher struct {
	UUID             uuid.UUID `json:"uuid"`
	Name             string    `json:"name"`
	Description      string    `json:"description,omitempty"`
	PublisherClass   string    `json:"publisherClass"`
	Template         string    `json:"template,omitempty"`
	TemplateMIMEType string    `json:"templateMimeType"`
	DefaultPublisher bool      `json:"defaultPublisher"`
}

type NotificationRule struct {
	UUID                    uuid.UUID                   `json:"uuid"`
	Name                    string                      `json:"name"`
	Enabled                 bool                        `json:"enabled"`
	NotifyChildren          bool                        `json:"notifyChildren"`
	LogSuccessfulPublish    bool                        `json:"logSuccessfulPublish"`
	Scope                   NotificationRuleScope       `json:"scope"`
	NotificationLevel       NotificationRuleLevel       `json:"notificationLevel,omitempty"`
	NotifyOn                []NotificationRuleNotifyOn  `json:"notifyOn,omitempty"`
	TriggerType             NotificationRuleTriggerType `json:"triggerType"`
	Message                 string                      `json:"message,omitempty"`
	PublisherConfig         string                      `json:"publisherConfig,omitempty"`
	ScheduleLastTriggeredAt int64                       `json:"scheduleLastTriggeredAt,omitempty"`
	ScheduleNextTriggerAt   int64                       `json:"scheduleNextTriggerAt,omitempty"`
	ScheduleCron            string                      `json:"scheduleCron,omitempty"`
	ScheduleSkipUnchanged   bool                        `json:"scheduleSkipUnchanged,omitempty"`
	FilterExpression        string                      `json:"filterExpression,omitempty"`
	Publisher               NotificationPublisher       `json:"publisher,omitempty"`
	Projects                []Project                   `json:"projects,omitempty"`
	Tags                    []Tag                       `json:"tags,omitempty"`
	Teams                   []Team                      `json:"teams,omitempty"`
}

type CreateScheduledNotificationRuleRequest struct {
	Name              string                `json:"name"`
	Scope             NotificationRuleScope `json:"scope"`
	NotificationLevel NotificationRuleLevel `json:"notificationLevel"`
	Publisher         Publisher             `json:"publisher"`
}

type Publisher struct {
	UUID uuid.UUID `json:"uuid"`
}

type NotificationRuleScope string

const (
	NotificationRuleScopeSystem    NotificationRuleScope = "SYSTEM"
	NotificationRuleScopePortfolio NotificationRuleScope = "PORTFOLIO"
)

type NotificationRuleLevel string

const (
	NotificationRuleLevelInformational NotificationRuleLevel = "INFORMATIONAL"
	NotificationRuleLevelWarning       NotificationRuleLevel = "WARNING"
	NotificationRuleLevelError         NotificationRuleLevel = "ERROR"
)

type NotificationRuleNotifyOn string

type NotificationRuleTriggerType string

const (
	NotificationRuleTriggerTypeEvent    NotificationRuleTriggerType = "EVENT"
	NotificationRuleTriggerTypeSchedule NotificationRuleTriggerType = "SCHEDULE"
)

func (ns NotificationService) GetAllPublishers(ctx context.Context) (ps []NotificationPublisher, err error) {
	err = ns.client.assertServerVersionAtLeast("3.2.0")
	if err != nil {
		return
	}

	req, err := ns.client.newRequest(ctx, http.MethodGet, "api/v1/notification/publisher")
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &ps)
	return
}

func (ns NotificationService) CreatePublisher(ctx context.Context, publisher NotificationPublisher) (p NotificationPublisher, err error) {
	err = ns.client.assertServerVersionAtLeast("4.6.0")
	if err != nil {
		return
	}

	req, err := ns.client.newRequest(ctx, http.MethodPut, "api/v1/notification/publisher", withBody(publisher))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &p)
	return
}

func (ns NotificationService) UpdatePublisher(ctx context.Context, publisher NotificationPublisher) (p NotificationPublisher, err error) {
	err = ns.client.assertServerVersionAtLeast("4.6.0")
	if err != nil {
		return
	}

	req, err := ns.client.newRequest(ctx, http.MethodPost, "api/v1/notification/publisher", withBody(publisher))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &p)
	return
}

func (ns NotificationService) DeletePublisher(ctx context.Context, publisherUUID uuid.UUID) (err error) {
	err = ns.client.assertServerVersionAtLeast("4.6.0")
	if err != nil {
		return
	}

	req, err := ns.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("api/v1/notification/publisher/%s", publisherUUID.String()))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, nil)
	return
}

func (ns NotificationService) RestoreDefaultTemplates(ctx context.Context) (err error) {
	err = ns.client.assertServerVersionAtLeast("4.6.0")
	if err != nil {
		return
	}

	req, err := ns.client.newRequest(ctx, http.MethodPost, "api/v1/notification/publisher/restoreDefaultTemplates")
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, nil)
	return
}

func (ns NotificationService) TestRule(ctx context.Context, ruleUUID uuid.UUID) (err error) {
	err = ns.client.assertServerVersionAtLeast("4.12.0")
	if err != nil {
		return
	}

	req, err := ns.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("api/v1/notification/publisher/test/%s", ruleUUID.String()))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, nil)
	return
}

func (ns NotificationService) TestSMTP(ctx context.Context, destination string) (err error) {
	err = ns.client.assertServerVersionAtLeast("3.4.0")
	if err != nil {
		return
	}
	values := url.Values{}
	values.Set("destination", destination)

	req, err := ns.client.newRequest(ctx, http.MethodPost, "api/v1/notification/publisher/test/smtp", withBody(values))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, nil)
	return
}

func (ns NotificationService) AddProjectToRule(ctx context.Context, ruleUUID, projectUUID uuid.UUID) (r NotificationRule, err error) {
	err = ns.client.assertServerVersionAtLeast("3.2.0")
	if err != nil {
		return
	}

	req, err := ns.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("api/v1/notification/rule/%s/project/%s", ruleUUID.String(), projectUUID.String()))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &r)
	return
}

func (ns NotificationService) RemoveProjectFromRule(ctx context.Context, ruleUUID, projectUUID uuid.UUID) (r NotificationRule, err error) {
	err = ns.client.assertServerVersionAtLeast("3.2.0")
	if err != nil {
		return
	}

	req, err := ns.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("api/v1/notification/rule/%s/project/%s", ruleUUID.String(), projectUUID.String()))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &r)
	return
}

func (ns NotificationService) AddTeamToRule(ctx context.Context, ruleUUID, teamUUID uuid.UUID) (r NotificationRule, err error) {
	err = ns.client.assertServerVersionAtLeast("4.7.0")
	if err != nil {
		return
	}

	req, err := ns.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("api/v1/notification/rule/%s/team/%s", ruleUUID.String(), teamUUID.String()))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &r)
	return
}

func (ns NotificationService) RemoveTeamFromRule(ctx context.Context, ruleUUID, teamUUID uuid.UUID) (r NotificationRule, err error) {
	err = ns.client.assertServerVersionAtLeast("4.7.0")
	if err != nil {
		return
	}

	req, err := ns.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("api/v1/notification/rule/%s/team/%s", ruleUUID.String(), teamUUID.String()))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &r)
	return
}

type GetAllRulesFilterOptions struct {
	TriggerType NotificationRuleTriggerType
}

func withGetAllRulesFilterOptions(filterOptions GetAllRulesFilterOptions) requestOption {
	return func(req *http.Request) error {
		query := req.URL.Query()
		if len(filterOptions.TriggerType) > 0 {
			query.Set("triggerType", string(filterOptions.TriggerType))
		}
		req.URL.RawQuery = query.Encode()
		return nil
	}
}

func (ns NotificationService) GetAllRules(ctx context.Context, po PageOptions, so SortOptions, filterOptions GetAllRulesFilterOptions) (p Page[NotificationRule], err error) {
	err = ns.client.assertServerVersionAtLeast("3.2.0")
	if err != nil {
		return
	}

	req, err := ns.client.newRequest(ctx, http.MethodGet, "api/v1/notification/rule", withPageOptions(po), withSortOptions(so), withGetAllRulesFilterOptions(filterOptions))
	if err != nil {
		return
	}

	res, err := ns.client.doRequest(req, &p.Items)
	if err != nil {
		return
	}

	p.TotalCount = res.TotalCount
	return
}

func (ns NotificationService) CreateRule(ctx context.Context, ruleReq NotificationRule) (ruleRes NotificationRule, err error) {
	err = ns.client.assertServerVersionAtLeast("3.2.0")
	if err != nil {
		return
	}

	req, err := ns.client.newRequest(ctx, http.MethodPut, "api/v1/notification/rule", withBody(ruleReq))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &ruleRes)
	return
}

func (ns NotificationService) UpdateRule(ctx context.Context, ruleReq NotificationRule) (ruleRes NotificationRule, err error) {
	err = ns.client.assertServerVersionAtLeast("3.2.0")
	if err != nil {
		return
	}

	req, err := ns.client.newRequest(ctx, http.MethodPost, "api/v1/notification/rule", withBody(ruleReq))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &ruleRes)
	return
}

func (ns NotificationService) DeleteRule(ctx context.Context, rule NotificationRule) (err error) {
	err = ns.client.assertServerVersionAtLeast("3.2.0")
	if err != nil {
		return
	}

	req, err := ns.client.newRequest(ctx, http.MethodDelete, "api/v1/notification/rule", withBody(rule))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, nil)
	return
}

func (ns NotificationService) CreateScheduledRule(ctx context.Context, schedule CreateScheduledNotificationRuleRequest) (r NotificationRule, err error) {
	err = ns.client.assertServerVersionAtLeast("4.13.0")
	if err != nil {
		return
	}

	req, err := ns.client.newRequest(ctx, http.MethodPut, "api/v1/notification/rule/scheduled", withBody(schedule))
	if err != nil {
		return
	}

	_, err = ns.client.doRequest(req, &r)
	return
}
