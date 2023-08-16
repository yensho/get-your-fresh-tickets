package api

type AuthCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SpaceCreateRequest struct {
	Name  string `json:"name"`
	Areas []Area `json:"areas,omitempty"`
}

type Space struct {
	Name    string `db:"space_nm"`
	Section string `db:"space_section_nm"`
	Seats   []byte `db:"space_section_seats"`
}

type Area struct {
	Section string   `json:"section_name"`
	Seats   []string `json:"seats"`
}

type Customer struct {
	Name     string  `json:"name"`
	Age      int     `json:"age"`
	Gender   string  `json:"gender"`
	PhoneNum string  `json:"phone_nm"`
	Address  Address `json:"addr"`
	Email    string  `json:"email"`
}

type Address struct {
	Line1    string `json:"line1"`
	Line2    string `json:"line2"`
	City     string `json:"city"`
	State    string `json:"state"`
	PostalCd string `json:"zip"`
}

type Show struct {
	Name        string                 `json:"name"`
	Location    string                 `json:"location"`
	ShowTimes   []string               `json:"show_times"`
	Seats       []string               `json:"seats"`
	Prices      map[string]interface{} `json:"prices"`
	AgeRestrict bool                   `json:"age_restrict"`
}
