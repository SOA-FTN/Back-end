package DtoObjects

type ProfileDto struct {
	Id       uint     `json:"Id"`
	Email    string   `json:"Email"`
	UserName string   `json:"Username"`
	Role     UserRole `json:"Role"`
	IsActive bool     `json:"IsActive"`
}

type UserRole int
