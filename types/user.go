package types

// User represents the authenticated user
type User struct {
	ID       string
	Username string
	Email    string
	Role     string
}
