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
	//? url host wthout schema parsing as path, see func description
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

// matchPath extracting path variables with template.
// it returns nil, false, if path doesn't match template
// got example: matchPath("/joinchat/{chat_id}", "/joinchat/abcdefg") returns {"chat_id":"abcdefg"}, false
// spiced up implementaition from https://git.io/Jtcv0 (cuz why not?)
func matchPath(tpl, path string) (map[string]string, bool) {
	//? if template doesn't have pattern
	if !strings.ContainsAny(tpl, "{}") {
		if tpl == path {
			return map[string]string{}, true
		}
		return nil, false
	}

	//? if template or path are not global filepath
	if !strings.HasPrefix(tpl, "/") || !strings.HasPrefix(tpl, "/") {
		return nil, false
	}
	tplPathItems := strings.Split(tpl, "/")
	pathItems := strings.Split(path, "/")
	if len(tplPathItems) != len(pathItems) {
		return nil, false
	}

	res := make(map[string]string)
	for i, tplPathItem := range tplPathItems {
		//? if this item not a variable, we just need to check it but don't extract
		if !strings.HasPrefix(tplPathItem, "{") || !strings.HasSuffix(tplPathItem, "}") {
			if tplPathItem != pathItems[i] {
				return nil, false
			}
			continue
		}

		//? {chat_id} -> chat_id
		templateObject := strings.TrimSuffix(strings.TrimPrefix(tplPathItem, "{"), "}")
		// TODO: decide, do we REALLY need check patterns?
		//if !stringsContainsOnlyFunc(templateObject, validIdentRunes) {
		//	panic("got invalid template") // panicing, cause we must check templates BEFORE it's usage
		//}
		res[templateObject] = pathItems[i]
	}

	return res, true
}

// func validIdentRunes(c rune) bool {
// 	return (c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'z' && c <= 'z') || c == '_'
// }
//
// func stringsContainsOnlyFunc(i string, f func(rune) bool) bool {
// 	for _, c := range []rune(i) {
// 		if !f(c) {
// 			return false
// 		}
// 	}
// 	return true
// }

// https://stackoverflow.com/questions/62083272/parsing-url-with-port-and-without-scheme
func fixURLHost(u *url.URL) {
	if u.Host != "" {
		return
	}
	if strings.HasPrefix(u.Path, "/") || u.Path == "" {
		// it's valid, something like '/t.me/somepath' is EXACTLY a path. But 't.me/somepath' has host
		return
	}
	i := strings.IndexRune(u.Path, '/')
	u.Host = u.Path[:i]
	u.Path = u.Path[i:]
}

func stringListContains(l []string, s string) bool {
	for i := range l {
		if l[i] == s {
			return true
		}
	}
	return false
}
