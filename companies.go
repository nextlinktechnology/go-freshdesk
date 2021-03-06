package freshdesk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nextlinktechnology/mgm/v3"
)

type CompanyManager interface {
	All() (CompanySlice, error)
	Create(CreateCompany) (Company, error)
	Update(int64, CreateCompany) (Company, error)
}

type companyManager struct {
	client *ApiClient
}

func newCompanyManager(client *ApiClient) companyManager {
	return companyManager{
		client,
	}
}

type Company struct {
	mgm.DefaultModel `bson:",inline"`
	ID               int64                  `bson:"id" json:"id"`
	Name             string                 `bson:"name" json:"name"`
	Description      string                 `bson:"description" json:"description"`
	Domains          []string               `bson:"domains" json:"domains"`
	Note             string                 `bson:"note" json:"note"`
	HealthScore      string                 `bson:"health_score" json:"health_score"`
	AccountTier      string                 `bson:"account_tier" json:"account_tier"`
	RenewalDate      *time.Time             `bson:"renewal_date" json:"renewal_date"`
	Industry         string                 `bson:"industry" json:"industry"`
	CreatedAt        *time.Time             `bson:"created_at" json:"created_at"`
	UpdatedAt        *time.Time             `bson:"updated_at" json:"updated_at"`
	CustomFields     map[string]interface{} `bson:"custom_fields" json:"custom_fields"`
}

type CreateCompany struct {
	Name         string                 `json:"name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Domains      []string               `json:"domains,omitempty"`
	Note         string                 `json:"note,omitempty"`
	HealthScore  string                 `json:"health_score,omitempty"`
	AccountTier  string                 `json:"account_tier,omitempty"`
	RenewalDate  *time.Time             `json:"renewal_date,omitempty"`
	Industry     string                 `json:"industry,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

type CompanySlice []Company

func (c CompanySlice) Len() int {
	return len(c)
}

func (c CompanySlice) Less(i, j int) bool {
	return c[i].ID < c[j].ID
}

func (c CompanySlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c CompanySlice) Print() {
	for _, company := range c {
		fmt.Println(company.Name)
	}
}

func (manager companyManager) All() (CompanySlice, error) {
	output := CompanySlice{}
	headers, err := manager.client.get(endpoints.companies.all, &output)
	if err != nil {
		return CompanySlice{}, err
	}
	for {
		nextLink := manager.client.getNextLink(headers)
		if nextLink == "" {
			break
		}
		nextSlice := CompanySlice{}
		headers, err = manager.client.get(nextLink, &nextSlice)
		if err != nil {
			return CompanySlice{}, err
		}
		output = append(output, nextSlice...)
	}
	return output, nil
}

func (manager companyManager) Create(company CreateCompany) (Company, error) {
	output := Company{}
	jsonb, err := json.Marshal(company)
	if err != nil {
		return output, err
	}
	err = manager.client.postJSON(endpoints.companies.create, jsonb, &output, http.StatusCreated)
	if err != nil {
		return Company{}, err
	}
	return output, nil
}

func (manager companyManager) Update(id int64, company CreateCompany) (Company, error) {
	output := Company{}
	jsonb, err := json.Marshal(company)
	if err != nil {
		return output, err
	}
	err = manager.client.put(endpoints.companies.update(id), jsonb, &output, http.StatusOK)
	if err != nil {
		return Company{}, err
	}
	return output, nil
}
