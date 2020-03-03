package freshdesk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Kamva/mgm/v2"
	"github.com/nextlinktechnology/go-freshdesk/querybuilder"
)

type UserManager interface {
	All() (UserSlice, error)
	Create(*User) (*User, error)
	Search(querybuilder.Query) (UserResults, error)
}

type userManager struct {
	client *ApiClient
}

func newUserManager(client *ApiClient) userManager {
	return userManager{
		client,
	}
}

type UserResults struct {
	next    string
	Results UserSlice
	client  *ApiClient
}

type User struct {
	mgm.DefaultModel `bson:",inline"`
	ID               int64                  `bson:"id" json:"id,omitempty"`
	Name             string                 `bson:"name" json:"name,omitempty"`
	Active           string                 `bson:"active" json:"active,omitempty"`
	Email            string                 `bson:"email" json:"email,omitempty"`
	JobTitle         string                 `bson:"job_title" json:"job_title,omitempty"`
	Language         string                 `bson:"language" json:"language,omitempty"`
	LastLoginAt      *time.Time             `bson:"last_login_at" json:"last_login_at,omitempty"`
	Mobile           int                    `bson:"mobile" json:"mobile,omitempty"`
	Phone            int                    `bson:"phone" json:"phone,omitempty"`
	TimeZone         string                 `bson:"time_zone" json:"time_zone,omitempty"`
	CreatedAt        *time.Time             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt        *time.Time             `bson:"updated_at" json:"updated_at,omitempty"`
	Address          string                 `bson:"address" json:"address,omitempty"`
	Avatar           interface{}            `bson:"avatar" json:"avatar,omitempty"`
	CompanyID        int64                  `bson:"company_id" json:"company_id,omitempty"`
	ViewAllTickets   bool                   `bson:"view_all_tickets" json:"view_all_tickets,omitempty"`
	CustomFields     map[string]interface{} `bson:"custom_fields" json:"custom_fields,omitempty"`
	Deleted          bool                   `bson:"deleted" json:"deleted,omitempty"`
	Description      string                 `bson:"description" json:"description,omitempty"`
	OtherEmails      []string               `bson:"other_emails" json:"other_emails,omitempty"`
	Tags             []string               `bson:"tags" json:"tags,omitempty"`
	TwitterID        string                 `bson:"twitter_id" json:"twitter_id,omitempty"`
	UniqueExternalID string                 `bson:"unique_external_id" json:"unique_external_id,omitempty"`
	OtherCompanies   []string               `bson:"other_companies" json:"other_companies,omitempty"`
}

type UserSlice []User

func (s UserSlice) Len() int { return len(s) }

func (s UserSlice) Less(i, j int) bool { return s[i].ID < s[j].ID }

func (s UserSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s UserSlice) Print() {
	for _, user := range s {
		fmt.Println(user.Name)
	}
}

func (manager userManager) All() (UserSlice, error) {
	output := UserSlice{}
	headers, err := manager.client.get(endpoints.contacts.all, &output)
	if err != nil {
		return UserSlice{}, err
	}
	for {
		nextLink := manager.client.getNextLink(headers)
		if nextLink == "" {
			break
		}
		nextSlice := UserSlice{}
		headers, err = manager.client.get(nextLink, &nextSlice)
		if err != nil {
			return UserSlice{}, err
		}
		output = append(output, nextSlice...)
	}
	return output, nil
}

func (manager userManager) Search(query querybuilder.Query) (UserResults, error) {
	output := struct {
		Slice UserSlice `json:"results,omitempty"`
	}{}
	headers, err := manager.client.get(endpoints.contacts.search(query.URLSafe()), &output)
	if err != nil {
		return UserResults{}, err
	}
	return UserResults{
		Results: output.Slice,
		client:  manager.client,
		next:    manager.client.getNextLink(headers),
	}, nil
}

func (manager userManager) Create(user *User) (*User, error) {
	output := &User{}
	jsonb, err := json.Marshal(user)
	if err != nil {
		return output, err
	}
	err = manager.client.postJSON(endpoints.contacts.create, jsonb, &output, http.StatusCreated)
	if err != nil {
		return &User{}, err
	}
	return output, nil
}

func (manager userManager) Update(id int64, user *User) (*User, error) {
	output := &User{}
	jsonb, err := json.Marshal(user)
	if err != nil {
		return output, err
	}
	err = manager.client.put(endpoints.contacts.update(id), jsonb, &output, http.StatusOK)
	if err != nil {
		return &User{}, err
	}
	return output, nil
}
