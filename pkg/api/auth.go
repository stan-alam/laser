package api

// @todo: abstract Login & Token Request behavior?
type Auth struct {
	User  authStorage
	Token tokenStorage
}

// @todo: define login operation; eg. accept basic auth or bearer refresh token, return refresh & access token
// @todo: define token operation; eg. generate new access token from refresh token
