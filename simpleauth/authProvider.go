package simpleauth

type AuthProvider struct {
	ReadAuthInfo  func(string) (AuthInfo, error)
	WriteAuthInfo func(AuthInfo) error
}

var authProvider AuthProvider

// This sets up simpleauth using the auth provider.
func Setup(provider AuthProvider) {
	authProvider = provider
}
