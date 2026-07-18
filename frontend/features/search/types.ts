export interface SearchHit {
  id: string;
  module_id: string;
  module_label: string;
  api_name: string;
  title: string;
  subtitle?: string;
}

export interface SearchResponse {
  results: SearchHit[];
}
