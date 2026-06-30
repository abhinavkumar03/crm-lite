package dto

import (
	contactDto "github.com/abhinavkumar03/crm-lite/backend/internal/contact/dto"
	leadDto "github.com/abhinavkumar03/crm-lite/backend/internal/lead/dto"
	taskDto "github.com/abhinavkumar03/crm-lite/backend/internal/task/dto"
)

type LeadResult struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Company string `json:"company"`
	Status  string `json:"status"`
}

type ContactResult struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type TaskResult struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type SearchResponse struct {
	Leads    []leadDto.LeadResponse       `json:"leads"`
	Contacts []contactDto.ContactResponse `json:"contacts"`
	Tasks    []taskDto.TaskResponse       `json:"tasks"`
}
