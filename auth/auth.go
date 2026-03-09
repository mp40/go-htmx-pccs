package auth

// temp in memory till need to impliment
type Auth struct {
	signedIn bool
}

func NewAuthService() *Auth {
	return &Auth{}
}

func (a *Auth) SignIn(email string, password string) error {
	a.signedIn = true
	return nil
}

func (a *Auth) SignUp(email string, password string) error {
	a.signedIn = true
	return nil
}
