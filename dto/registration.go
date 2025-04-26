package dto

type RegistrationRequest struct {
	Name        string `json:"nama" validate:"required"`
	NIK         string `json:"nik" validate:"required"`
	PhoneNumber string `json:"no_hp" validate:"required"`
}

type RegistrationResponse struct {
	AccountNumber string `json:"no_rekening,omitempty"`
	Remark        string `json:"remark,omitempty"`
}
