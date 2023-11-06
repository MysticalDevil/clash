package auth

import (
	"github.com/MysticalDevil/clash/component/auth"
)

var authenticator auth.Authenticator

func Authenticator() auth.Authenticator {
	return authenticator
}

func SetAuthenticator(au auth.Authenticator) {
	authenticator = au
}
