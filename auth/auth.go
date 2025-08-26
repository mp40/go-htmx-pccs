package auth

// temp in memory till need to impliment
type Auth struct {
	signedIn bool
}

func NewAuthService() *Auth {
	return &Auth{}
}

func (a *Auth) SignIn() error {
	a.signedIn = true
	return nil
}
