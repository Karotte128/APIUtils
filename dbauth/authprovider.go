package dbauth

import (
	"errors"

	"github.com/karotte128/apiutils/v2/database"
	"github.com/karotte128/apiutils/v2/simpleauth"
)

// This is a database based permission provider for simpleauth.
func GetAuthProvider() (simpleauth.AuthProvider, error) {
	if database.ConnPool == nil {
		return simpleauth.AuthProvider{}, errors.New("No database connection!")
	}

	return simpleauth.AuthProvider{
		ReadAuthInfo: func(key string) (simpleauth.AuthInfo, error) {
			info, err := GetAuth(database.ConnPool, "authentication", key)
			if err != nil {
				return simpleauth.AuthInfo{}, err
			}

			return simpleauth.AuthInfo{
				ApiKey:      info.ApiKey,
				Permissions: info.Permissions,
				ValidUntil:  info.ValidUntil,
				Info:        info.Info,
			}, nil
		},

		WriteAuthInfo: func(authInfo simpleauth.AuthInfo) error {
			info := DbAuthInfo{
				ApiKey:      authInfo.ApiKey,
				Permissions: authInfo.Permissions,
				ValidUntil:  authInfo.ValidUntil,
				Info:        authInfo.Info,
			}

			return UpdateAuth(database.ConnPool, "authentication", info)
		},
	}, nil
}
