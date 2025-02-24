package model

type Certificate struct {
	CertificateType      CertificateType `json:"certificateType,omitempty"`
	CertificateNo        string          `json:"certificateNo,omitempty"`
	HolderName           *UserName       `json:"holderName,omitempty"`
	FileKeys             []string        `json:"fileKeys,omitempty"`
	CertificateAuthority string          `json:"certificateAuthority,omitempty"`
	GrantType            string          `json:"grantType,omitempty"`
}

type CertificateType string

const (
	//ENTERPRISE_REGISTRATION,
	//LICENSE_INFO,
	//ID_CARD,
	//PASSPORT,
	//DRIVING_LICENSE,
	//CPF,
	//CNPJ,

	CertificateType_ENTERPRISE_REGISTRATION CertificateType = "ENTERPRISE_REGISTRATION"
	CertificateType_LICENSE_INFO            CertificateType = "LICENSE_INFO"
	CertificateType_ID_CARD                 CertificateType = "ID_CARD"
	CertificateType_PASSPORT                CertificateType = "PASSPORT"
	CertificateType_DRIVING_LICENSE         CertificateType = "DRIVING_LICENSE"
	CertificateType_CPF                     CertificateType = "CPF"
	CertificateType_CNPJ                    CertificateType = "CNPJ"
)
