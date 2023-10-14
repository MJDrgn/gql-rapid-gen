package util

import (
	"strings"
	"unicode"
)

func split(in string) (ret []string) {
	if len(in) == 0 {
		return nil
	}
	if strings.ContainsRune(in, '_') {
		ret = strings.Split(in, "_")
	} else if strings.ContainsRune(in, '-') {
		ret = strings.Split(in, "-")
	} else if strings.ContainsRune(in, ' ') {
		ret = strings.Split(in, " ")
	} else if strings.ToUpper(in) == in {
		// If the whole thing is uppercase, assume ENUM or something similar like ID that should be a single word.
		ret = []string{in}
	} else {
		// Assume its case separated (TitleCase or camelCase)
		containsUpper := false
		indexes := make([]int, 0)
		for idx, r := range in {
			if unicode.IsUpper(r) {
				indexes = append(indexes, idx)
				containsUpper = true
			}
		}

		// If lowercase only, just return as is.
		if !containsUpper {
			return []string{in}
		}

		prevIdx := 0
		for _, idx := range indexes {
			if idx == 0 {
				continue
			}
			ret = append(ret, in[prevIdx:idx])
			prevIdx = idx
		}
		ret = append(ret, in[prevIdx:])
	}
	for idx := range ret {
		ret[idx] = strings.ToLower(ret[idx])
	}
	return ret
}

func TitleCase(in string) string {
	var b strings.Builder
	b.Grow(len(in))
	sp := split(in)
	for _, w := range sp {
		rs := []rune(w)
		b.WriteRune(unicode.ToUpper(rs[0]))
		if len(rs) > 1 {
			b.WriteString(string(rs[1:]))
		}
	}
	return b.String()
}

func CamelCase(in string) string {
	var b strings.Builder
	b.Grow(len(in))
	sp := split(in)
	b.WriteString(sp[0])
	if len(sp) > 1 {
		for _, w := range sp[1:] {
			rs := []rune(w)
			b.WriteRune(unicode.ToUpper(rs[0]))
			if len(rs) > 1 {
				b.WriteString(string(rs[1:]))
			}
		}
	}
	return b.String()
}

func DashCase(in string) string {
	sp := split(in)
	return strings.Join(sp, "-")
}

func UnderCase(in string) string {
	sp := split(in)
	return strings.Join(sp, "_")
}

func LabelCase(in string) string {
	sp := split(in)

	// Uppercase first letter of each segment
	for i := range sp {
		if sp[i] == "id" {
			// ID is an exception
			sp[i] = "ID"
			continue
		}
		sp[i] = string(unicode.ToUpper([]rune(sp[i])[0])) + sp[i][1:]
	}
	return strings.Join(sp, " ")
}
