package laser

type User struct {
	Name     string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
