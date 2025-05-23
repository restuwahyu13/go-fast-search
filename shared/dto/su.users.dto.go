package dto

type (
	CreateUsersDTO struct {
		Name        string `json:"name" validate:"required"`
		Email       string `json:"email" validate:"required,email"`
		Phone       string `json:"phone" validate:"required,e164"`
		DateOfBirth string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
		Age         string `json:"age" validate:"required,min=2,max=3"`
		Address     string `json:"address" validate:"required"`
		City        string `json:"city" validate:"required"`
		State       string `json:"state" validate:"required"`
		Direction   string `json:"direction" validate:"required"`
		Country     string `json:"country" validate:"required"`
		PostalCode  string `json:"postal_code" validate:"required,,len=4"`
	}

	UpdateUsersDTO struct {
		ID          string `json:"id,omitempty" validate:"omitempty,uuidv4"`
		Name        string `json:"name,omitempty" validate:"omitempty"`
		Email       string `json:"email,omitempty" validate:"omitempty,email"`
		Phone       string `json:"phone,omitempty" validate:"omitempty,e164"`
		DateOfBirth string `json:"date_of_birth,omitempty" omitempty:"omitempty,datetime=2006-01-02"`
		Age         string `json:"age,omitempty" validate:"omitempty,min=2,max=3"`
		Address     string `json:"address,omitempty" validate:"omitempty"`
		City        string `json:"city,omitempty" validate:"omitempty"`
		State       string `json:"state,omitempty" validate:"omitempty"`
		Direction   string `json:"direction,omitempty" validate:"omitempty"`
		Country     string `json:"country,omitempty" validate:"omitempty"`
		PostalCode  string `json:"postal_code,omitempty" validate:"omitempty,len=4"`
	}
)
