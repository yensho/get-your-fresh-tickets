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
