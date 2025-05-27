package entitie

type UsersDocument struct {
	ID            string         `json:"id"`
	Name          string         `json:"name" `
	Email         string         `json:"email" `
	Phone         string         `json:"phone"`
	DateOfBirth   string         `json:"date_of_birth"`
	Age           string         `json:"age"`
	Address       string         `json:"address"`
	City          string         `json:"city"`
	State         string         `json:"state"`
	Direction     string         `json:"direction"`
	Country       string         `json:"country"`
	PostalCode    string         `json:"postal_code"`
	CreatedAt     int64          `json:"created_at"`
	UpdatedAt     int64          `json:"updated_at,omitempty"`
	DeletedAt     int64          `json:"deleted_at,omitempty"`
	Formatted     map[string]any `json:"_formatted,omitempty"`
	MatchPosition map[string]any `json:"_matchesPosition,omitempty"`
}
