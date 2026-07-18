package dto

type SearchHit struct {
	ID          string `json:"id"`
	ModuleID    string `json:"module_id"`
	ModuleLabel string `json:"module_label"`
	APIName     string `json:"api_name"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle,omitempty"`
}

type SearchResponse struct {
	Results []SearchHit `json:"results"`
}
