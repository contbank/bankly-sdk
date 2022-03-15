package models

// CreateTicketRequest ...
type CreateTicketRequest struct {
	GroupID     int      `json:"group_id"`
	ProductID   int      `json:"product_id"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Subject     string   `json:"subject"`
	Email       string   `json:"email"`
	CopyEmails  []string `json:"copy_emails"`
	Priority    int      `json:"priority"`
	Status      int      `json:"status"`
	Attachments []string `json:"attachments"`
	Request     string   `json:"cf_solicitacao"`
	Request2    string   `json:"cf_solicitacao2"`
	CompanyKey  string   `json:"cf_companykey"`
	RazaoSocial string   `json:"cf_razao_social"`
	CNPJ        string   `json:"cf_cnpj_empresa"`
	Cellphone   int      `json:"cf_celular"`
	Last4Digits int      `json:"cf_comprastransaesestornos_4_ltimos_dgitos_do_carto"`
}

// CreateTicketResponse ...
type CreateTicketResponse struct {
	ID int `json:"id"`
}

// GetTicketResponse ...
type GetTicketResponse struct {
	ID            int    `json:"id"`
	Status        int    `json:"status"`
	Reason        string `json:"cf_sdkycpjbankly"`
	AccountBranch string `json:"cf_sdkycpjbanklyagencia"`
	AccountNumber string `json:"cf_sdkycpjbanklyconta"`
}

// FilterTicketsRequest ...
type FilterTicketsRequest struct {
	Status int
}

// FilterTicketsResponse ...
type FilterTicketsResponse struct {
	Total   int                  `json:"total"`
	Results []*GetTicketResponse `json:"results"`
}
