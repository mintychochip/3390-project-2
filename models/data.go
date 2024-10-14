package models

type Users struct {
	ID       string `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}
type Prices struct {
	CMSCertificationNum string  `db:"cms_certification_num" json:"cms_certification_num"`
	Payer               string  `db:"payer" json:"payer"`
	Code                string  `db:"code" json:"code"`
	InternalRevenueCode string  `db:"internal_revenue_code" json:"internal_revenue_code"`
	Units               string  `db:"units" json:"units"`
	Description         string  `db:"description" json:"description"`
	InpatientOutpatient string  `db:"inpatient_outpatient" json:"inpatient_outpatient"`
	Price               float64 `db:"price" json:"price"`
	CodeDisambiguator   string  `db:"code_disambiguator" json:"code_disambiguator"`
}

type Hospital struct {
	ID                        string `db:"id" json:"id"`
	Ein                       string `db:"ein" json:"ein"`
	Name                      string `db:"name" json:"name"`
	AltName                   string `db:"alt_name" json:"alt_name"`
	SystemName                string `db:"system_name" json:"system_name"`
	Address                   string `db:"addr" json:"addr"`
	City                      string `db:"city" json:"city"`
	State                     string `db:"state" json:"state"`
	Zip                       string `db:"zip" json:"zip"`
	Phone                     string `db:"phone" json:"phone"`
	UrbanRural                string `db:"urban_rural" json:"urban_rural"`
	Category                  string `db:"category" json:"category"`
	ControlType               string `db:"control_type" json:"control_type"`
	MedicareTerminationStatus string `db:"medicare_termination_status" json:"medicare_termination_status"`
	LastUpdated               string `db:"last_updated" json:"last_updated"`
	FileName                  string `db:"file_name" json:"file_name"`
	MrfURL                    string `db:"mrf_url" json:"mrf_url"`
	Permalink                 string `db:"permalink" json:"permalink"`
	TransparencyPage          string `db:"transparency_page" json:"transparency_page"`
	AdditionalNotes           string `db:"additional_notes" json:"additional_notes"`
}
