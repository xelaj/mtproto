package deeplinks

import (
	"net/url"
	"strings"
)

func stringListContains(l []string, s string) bool {
	for i := range l {
		if l[i] == s {
			return true
		}
	}
	return false
}

func stringStringMapKeys(l map[string]string) []string {
	res := make([]string, 0, len(l))
	for k := range l {
		res = append(res, k)
	}
	return res
}

func stringsTrimSuffixPrefix(prefix, str, suffix string) string {
	return strings.TrimSuffix(strings.TrimPrefix(str, prefix), suffix)
}

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
