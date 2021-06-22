package deeplinks

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

var decoder = schema.NewDecoder()

func Resolve(link string) (Deeplink, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, errors.Wrap(err, "not a uri")
	}

	switch u.Scheme {
	case "", "http", "https":
		return resolveHttpLink(u)
	case "tg":
		return resolveTgLink(u)
	default:
		return nil, fmt.Errorf("'%v': invalid uri scheme (only tg://, http://, https:// are valid)", u.Scheme)
	}
}

func resolveHttpLink(u *url.URL) (Deeplink, error) {
	//? url host without schema parsing as path, see func description
	fixURLHost(u)

	//? Hostname(), cause MAYBE someone want add port for unknown reason. Logically it's better than u.Host
	if !stringListContains(ReservedHosts(), u.Hostname()) {
		return nil, fmt.Errorf("'%v' hostname is not owned by telegram", u.Hostname())
	}

	type pathVariablesConverter = func(path map[string]string, query url.Values) (Deeplink, error)

	for tpl, f := range map[string]pathVariablesConverter{
		"/joinchat/{token}": func(path map[string]string, query url.Values) (Deeplink, error) {
			token, ok := path["token"]
			if !ok || token == "" {
				return nil, errors.New("invite token required")
			}
			return &JoinParameters{
				Invite: token,
			}, nil
		},
		"/{username}": func(path map[string]string, query url.Values) (Deeplink, error) {
			username, ok := path["username"]
			if !ok || username == "" {
				return nil, errors.New("username required")
			}
			return &ResolveParameters{
				Domain: strings.ToLower(username),
			}, nil
		},
	} {
		vars, ok := matchPath(tpl, u.Path)
		if ok {
			return f(vars, u.Query())
		}
	}
	return nil, fmt.Errorf("'%v': this path does not look valid", u.Path)
}

func resolveTgLink(u *url.URL) (Deeplink, error) {
	return nil, errors.New("not implemented")
}
