package openssllib

// Openssl - structure of openssl config
type Openssl struct {
	CaPassword  string `json:"CaPassword"`
	CaEncrypted string `json:"CaEncrypted"`

	//Folders
	EasyRSAFolder    string `json:"EasyRSAFolder"`
	RequestFolder    string `json:"RequestFolder"`
	KeyFolder        string `json:"KeyFolder"`
	CertFolder       string `json:"CertFolder"`
	ConfigsFolder    string `json:"ConfigsFolder"`
	CertSerialFolder string `json:"CertSerialFolder"`

	// Filenames
	UserKeyFilename         string `json:"UserKeyFilename"`
	RequestFilename         string `json:"RequestFilename"`
	UserCertificateFilename string `json:"UserCertificateFilename"`
	SerialFile              string `json:"SerialFile"`
	SerialOldFile           string `json:"SerialOldFile"`
	OutputConfigFileName    string `json:"OutputConfigFileName"`
	CaCertFileName          string `json:"CaCertFileName"`
	CaKeyFileName           string `json:"CaKeyFileName"`
	TaFileName              string `json:"TaFileName"`
	TemplateFilename        string `json:"TemplateFilename"`

	NumberOfYearsValidity int `json:"NumberOfYearsValidity"`
}

type CertStruct struct {
	KeyBytes  []byte
	CertBytes []byte
	CABytes   []byte
	TABytes   []byte
}
