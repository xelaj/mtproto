package deeplinks

import (
	"fmt"
	"strings"
)

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
	if !strings.HasPrefix(tpl, "/") || !strings.HasPrefix(path, "/") {
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
		templateObject := stringsTrimSuffixPrefix("{", tplPathItem, "}")
		// TODO: decide, do we REALLY need check patterns?
		//if !stringsContainsOnlyFunc(templateObject, validIdentRunes) {
		//	panic("got invalid template") // panicing, cause we must check templates BEFORE it's usage
		//}
		res[templateObject] = pathItems[i]
	}

	return res, true
}

func fillTemplate(tpl string, data map[string]string) (string, error) {
	//? if template doesn't have pattern
	if !strings.ContainsAny(tpl, "{}") {
		if len(data) > 0 {
			return "", fmt.Errorf("unused keys: [%v]", strings.Join(stringStringMapKeys(data), ", "))
		}
		return tpl, nil
	}

	var isAbstract bool
	//? if template or path are not global filepath
	if !strings.HasPrefix(tpl, "/") {
		isAbstract = true
	}
	tplPathItems := strings.Split(tpl, "/")
	for i, tplPathItem := range tplPathItems {
		if !strings.HasPrefix(tplPathItem, "{") || !strings.HasSuffix(tplPathItem, "}") {
			continue
		}

		//? {chat_id} -> chat_id
		dataKey := stringsTrimSuffixPrefix("{", tplPathItem, "}")
		v, ok := data[dataKey]
		if !ok {
			return "", fmt.Errorf("key '%v' not found", dataKey)
		}
		tplPathItems[i] = v
	}

	res := strings.Join(tplPathItems, "/")
	if !isAbstract {
		res = "/" + res
	}

	return res, nil
}
