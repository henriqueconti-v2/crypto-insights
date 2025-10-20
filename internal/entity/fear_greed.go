package entity

type FearGreedResponse struct {
	Data   FearGreedData   `json:"data"`
	Status FearGreedStatus `json:"status"`
}

type FearGreedData struct {
	Value               int    `json:"value"`
	UpdateTime          string `json:"update_time"`
	ValueClassification string `json:"value_classification"`
}

type FearGreedStatus struct {
	Timestamp    string `json:"timestamp"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Elapsed      int    `json:"elapsed"`
	CreditCount  int    `json:"credit_count"`
}

type FearGreedIndex struct {
	Value        int
	Classification string
}
