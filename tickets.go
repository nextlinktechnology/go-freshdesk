package freshdesk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nextlinktechnology/go-freshdesk/querybuilder"
)

type TicketManager interface {
	All() (TicketResults, error)
	Create(CreateTicket) (Ticket, error)
	View(int64) (Ticket, error)
	Search(querybuilder.Query) (TicketResults, error)
	Reply(int64, CreateReply) (Reply, error)
	Conversations(int64) (ConversationSlice, error)
	UpdatedSinceAll(string) (TicketResults, error)
}

type ticketManager struct {
	client *ApiClient
}

type TicketResults struct {
	next    string
	Results TicketSlice
	client  *ApiClient
}

func newTicketManager(client *ApiClient) ticketManager {
	return ticketManager{
		client,
	}
}

type Ticket struct {
	ID                     int64                  `xorm:"pk 'id'" json:"id"`
	Subject                string                 `xorm:"Text" json:"subject"`
	Type                   string                 `json:"type"`
	Description            string                 `xorm:"Text" json:"description"`
	Attachments            []interface{}          `json:"attachments"`
	CCEmails               []string               `xorm:"'cc_emails'" json:"cc_emails"`
	CompanyID              int64                  `xorm:"'company_id'" json:"company_id"`
	Deleted                bool                   `json:"deleted"`
	DescriptionText        string                 `xorm:"Text" json:"description_text"`
	DueBy                  *time.Time             `json:"due_by"`
	Email                  string                 `json:"email"`
	EmailConfigID          int64                  `xorm:"'email_config_id'" json:"email_config_id"`
	FacebookID             string                 `xorm:"'facebook_id'" json:"facebook_id"`
	FirstResponseDueBy     *time.Time             `xorm:"'fr_due_by'" json:"fr_due_by"`
	FirstResponseEscalated bool                   `xorm:"'fr_escalated'" json:"fr_escalated"`
	FwdEmails              []string               `json:"fwd_emails"`
	GroupID                int64                  `xorm:"'group_id'" json:"group_id"`
	IsEscalated            bool                   `json:"is_escalated"`
	Name                   string                 `json:"name"`
	Phone                  string                 `json:"phone"`
	Priority               int                    `json:"priority"`
	ProductID              int64                  `xorm:"'product_id'" json:"product_id"`
	ReplyCCEmails          []string               `xorm:"'reply_cc_emails'" json:"reply_cc_emails"`
	RequesterID            int64                  `xorm:"'requester_id'" json:"requester_id"`
	ResponderID            int64                  `xorm:"'responder_id'" json:"responder_id"`
	Source                 int                    `json:"source"`
	Spam                   bool                   `json:"spam"`
	Status                 int                    `json:"status"`
	Tags                   []string               `json:"tags"`
	ToEmails               []string               `json:"to_emails"`
	TwitterID              string                 `xorm:"'twitter_id'" json:"twitter_id"`
	CreatedAt              *time.Time             `json:"created_at"`
	UpdatedAt              *time.Time             `json:"updated_at"`
	CustomFields           map[string]interface{} `json:"custom_fields"`
}

type CreateTicket struct {
	Name               string                 `json:"name,omitempty"`
	RequesterID        int                    `json:"requester_id,omitempty"`
	Email              string                 `json:"email,omitempty"`
	FacebookID         string                 `json:"facebook_id,omitempty"`
	Phone              string                 `json:"phone,omitempty"`
	TwitterID          string                 `json:"twitter_id,omitempty"`
	UniqueExternalID   string                 `json:"unique_external_id,omitempty"`
	Subject            string                 `json:"subject,omitempty"`
	Type               string                 `json:"type,omitempty"`
	Status             int                    `json:"status,omitempty"`
	Priority           int                    `json:"priority,omitempty"`
	Description        string                 `json:"description,omitempty"`
	ResponderID        int                    `json:"responder_id,omitempty"`
	Attachments        []interface{}          `json:"attachments,omitempty"`
	CCEmails           []string               `json:"cc_emails,omitempty"`
	CustomFields       map[string]interface{} `json:"custom_fields,omitempty"`
	DueBy              *time.Time             `json:"due_by,omitempty"`
	EmailConfigID      int                    `json:"email_config_id,omitempty"`
	FirstResponseDueBy *time.Time             `json:"fr_due_by,omitempty"`
	GroupID            int                    `json:"group_id,omitempty"`
	ProductID          int                    `json:"product_id,omitempty"`
	Source             int                    `json:"source,omitempty"`
	Tags               []string               `json:"tags,omitempty"`
	CompanyID          int                    `json:"company_id,omitempty"`
}

type Conversation struct {
	BodyText     string   `xorm:"Text" json:"body_text"`
	Body         string   `xorm:"Text" json:"body"`
	ID           int64    `xorm:"pk 'id'" json:"id"`
	Incoming     bool     `json:"incoming"`
	ToEmails     []string `json:"to_emails"`
	Private      bool     `json:"private"`
	Source       int      `json:"source"`
	SupportEmail string   `json:"support_email"`
	TicketID     int64    `xorm:"pk 'ticket_id'" json:"ticket_id"`

	UserID    int64      `xorm:"'user_id'" json:"user_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	FromEmail string     `json:"from_email"`
	CCEmails  []string   `xorm:"'cc_emails'" json:"cc_emails"`
	BCCEmails []string   `xorm:"'bcc_emails'" json:"bcc_emails"`

	Attachments []interface{} `json:"attachments"`
}

type Reply struct {
	BodyText    string        `json:"body_text"`
	Body        string        `json:"body"`
	ID          int           `json:"id"`
	UserID      int           `json:"user_id"`
	FromEmail   string        `json:"from_email"`
	CCEmails    []string      `json:"cc_emails"`
	BCCEmails   []string      `json:"bcc_emails"`
	ToEmails    []string      `json:"to_emails"`
	TicketID    int           `json:"ticket_id"`
	RepliedTo   []string      `json:"replied_to"`
	Attachments []interface{} `json:"attachments"`
	CreatedAt   *time.Time    `json:"created_at"`
	UpdatedAt   *time.Time    `json:"updated_at"`
}

type CreateReply struct {
	Body        string        `json:"body,omitempty"`
	FromEmail   string        `json:"from_email,omitempty"`
	Attachments []interface{} `json:"attachments,omitempty"`
	UserID      int           `json:"user_id,omitempty"`
	CCEmails    []string      `json:"cc_emails,omitempty"`
	BCCEmails   []string      `json:"bcc_emails,omitempty"`
}

type Source int
type Status int
type Priority int

const (
	SourceEmail Source = 1 + iota
	SourcePortal
	SourcePhone
	_
	_
	_
	SourceChat
	SourceMobihelp
	SourceFeedbackWidget
	SourceOutboundEmail
)

const (
	StatusOpen Status = 2 + iota
	StatusPending
	StatusResolved
	StatusClosed
)

const (
	PriorityLow Priority = 1 + iota
	PriorityMedium
	PriorityHigh
	PriorityUrgent
)

func (s Source) Value() int {
	return int(s)
}

func (s Status) Value() int {
	return int(s)
}

func (p Priority) Value() int {
	return int(p)
}

func (t Ticket) Print() {
	jsonb, _ := json.MarshalIndent(t, "", "    ")
	fmt.Println(string(jsonb))
}
func (r Reply) Print() {
	jsonb, _ := json.MarshalIndent(r, "", "    ")
	fmt.Println(string(jsonb))
}

type TicketSlice []Ticket

func (s TicketSlice) Len() int { return len(s) }

func (s TicketSlice) Less(i, j int) bool { return s[i].ID < s[j].ID }

func (s TicketSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s TicketSlice) Print() {
	for _, ticket := range s {
		fmt.Println(ticket.Subject)
	}
}

type ConversationSlice []Conversation

func (s ConversationSlice) Len() int { return len(s) }

func (s ConversationSlice) Less(i, j int) bool { return s[i].CreatedAt.Unix() > s[j].CreatedAt.Unix() }

func (s ConversationSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ConversationSlice) Print() {
	for _, ticket := range s {
		fmt.Println(ticket.BodyText)
	}
}

func (manager ticketManager) All() (TicketResults, error) {
	output := TicketSlice{}
	headers, err := manager.client.get(endpoints.tickets.all, &output)
	if err != nil {
		return TicketResults{}, err
	}
	return TicketResults{
		Results: output,
		client:  manager.client,
		next:    manager.client.getNextLink(headers),
	}, nil
}

func (manager ticketManager) UpdatedSinceAll(timeString string) (TicketResults, error) {
	output := TicketSlice{}
	headers, err := manager.client.get(endpoints.tickets.updatedSinceAll(timeString), &output)
	if err != nil {
		return TicketResults{}, err
	}
	return TicketResults{
		Results: output,
		client:  manager.client,
		next:    manager.client.getNextLink(headers),
	}, nil
}

func (manager ticketManager) Create(ticket CreateTicket) (Ticket, error) {
	output := Ticket{}
	jsonb, err := json.Marshal(ticket)
	if err != nil {
		return output, err
	}
	err = manager.client.postJSON(endpoints.tickets.create, jsonb, &output, http.StatusCreated)
	if err != nil {
		return Ticket{}, err
	}
	return output, nil
}

func (manager ticketManager) View(id int64) (Ticket, error) {
	output := Ticket{}
	_, err := manager.client.get(endpoints.tickets.view(id), &output)
	if err != nil {
		return Ticket{}, err
	}

	return output, nil
}

func (manager ticketManager) Conversations(id int64) (ConversationSlice, error) {
	output := ConversationSlice{}
	_, err := manager.client.get(endpoints.tickets.conversations(id), &output)
	if err != nil {
		return ConversationSlice{}, err
	}
	return output, nil
}

func (manager ticketManager) Reply(id int64, reply CreateReply) (Reply, error) {
	output := Reply{}
	jsonb, err := json.Marshal(reply)
	if err != nil {
		return output, err
	}
	err = manager.client.postJSON(endpoints.tickets.reply(id), jsonb, &output, http.StatusCreated)
	if err != nil {
		return Reply{}, err
	}
	return output, nil
}

func (manager ticketManager) Search(query querybuilder.Query) (TicketResults, error) {
	output := struct {
		Slice TicketSlice `json:"results"`
		Total int         `json:"total"`
	}{}
	_, err := manager.client.get(endpoints.tickets.search(query.URLSafe()), &output)
	if err != nil {
		return TicketResults{}, err
	}

	page := 1
	for {
		if len(output.Slice) >= output.Total || page == 10 {
			break
		}
		page++
		nextSlice := struct {
			Slice TicketSlice `json:"results"`
			Total int         `json:"total"`
		}{}
		_, err := manager.client.get(
			fmt.Sprintf("%s&page=%d", endpoints.tickets.search(query.URLSafe()), page),
			&nextSlice,
		)
		if err != nil {
			break
		}

		output.Slice = append(output.Slice, nextSlice.Slice...)
		output.Total = nextSlice.Total
	}

	return TicketResults{
		Results: output.Slice,
		client:  manager.client,
	}, nil
}

func (results TicketResults) Next() (TicketResults, error) {
	if results.next == "" {
		return TicketResults{}, errors.New("no more tickets")
	}
	output := TicketSlice{}
	headers, err := results.client.get(results.next, &output)
	if err != nil {
		return TicketResults{}, err
	}
	return TicketResults{
		Results: output,
		client:  results.client,
		next:    results.client.getNextLink(headers),
	}, nil
}

func (results *TicketResults) FilterTags(tags ...string) *TicketResults {
	filtered := TicketSlice{}
	for _, ticket := range results.Results {
		_filterFlag := false
		for _, ticketTag := range ticket.Tags {
			for _, filterTag := range tags {
				if ticketTag == filterTag {
					_filterFlag = true
					break
				}
			}
		}
		if _filterFlag {
			continue
		}
		filtered = append(filtered, ticket)
	}
	results.Results = filtered
	return results
}

func (results *TicketResults) FilterTypes(filterTypes ...string) *TicketResults {
	filtered := TicketSlice{}
	for _, ticket := range results.Results {
		_filterFlag := false
		for _, filterType := range filterTypes {
			if ticket.Type == filterType {
				_filterFlag = true
				break
			}
		}
		if _filterFlag {
			continue
		}
		filtered = append(filtered, ticket)
	}
	results.Results = filtered
	return results
}

func (results *TicketResults) FilterGroups(filterGroups ...string) *TicketResults {
	groups, _ := results.client.Groups.All()
	filterIDs := []int64{}
	for _, group := range groups {
		for _, filterGroup := range filterGroups {
			if group.Name == filterGroup {
				filterIDs = append(filterIDs, group.ID)
			}
		}
	}
	return results.FilterGroupsID(filterIDs...)
}

func (results *TicketResults) FilterGroupsID(filterIDs ...int64) *TicketResults {
	filtered := TicketSlice{}
	for _, ticket := range results.Results {
		_filterFlag := false
		for _, filterID := range filterIDs {
			if ticket.GroupID == filterID {
				_filterFlag = true
				break
			}
		}
		if _filterFlag {
			continue
		}
		filtered = append(filtered, ticket)
	}
	results.Results = filtered
	return results
}
