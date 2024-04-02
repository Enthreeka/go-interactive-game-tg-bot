package entity

type Answer struct {
	ID             int    `json:"id"`
	Answer         string `json:"answer"`
	CostOfResponse int    `json:"cost_of_response"`
}
