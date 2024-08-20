package api

type ResponseContractCode struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Result  []ContractCode `json:"result"`
}

type ContractCode struct {
	SourceCode string `json:"SourceCode"`
	ABI        string `json:"ABI"`
	Proxy      string `json:"Proxy"`
}

type IsScam struct {
	IsHoneypot *bool   `json:"IsHoneypot"`
	Error      *string `json:"Error"`
}

type StaySafu struct {
	Result struct {
		IsToken  *bool `json:"isToken"`
		Warnings []*struct {
			Message string `json:"message"`
		} `json:"warnings"`
	} `json:"result"`
}
