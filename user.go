package laser

type User struct {
	Name     string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"email"`
}
