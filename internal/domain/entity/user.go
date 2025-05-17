package entity

const (
	UserStateInit = "init"
)

type User struct {
	UserID      int64
	Username    string
	FirstName   string
	LastName    string
	ContactInfo string
	State       string
	Context     map[string]any
	Admin       bool
}
