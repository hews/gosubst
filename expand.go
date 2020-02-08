package main

// NOTE: all of the below is copied, with a few changes, from
//       https://golang.org/src/os/env.go (go1.13). In essence, I want
//       to replicate os.ExpandEnv(), but only targeting variables of
//       the form:
//
//       ${VAR_NAME}, or ${varName}, etc.
//
//       ... excepting also the case of:
//
//       $${VAR_NAME} (which becomes ${VAR_NAME} in the output).

// Expand replaces ${var} in the string based on the mapping function.
// For example, Expand(s, os.Getenv) is (mostly) equivalent to
// os.ExpandEnv(s).
func Expand(s string, mapping func(string) string) string {
	var buf []byte
	// ${} is all ASCII, so bytes are fine for this operation.
	i := 0
	for j := 0; j < len(s); j++ {
		if j+1 < len(s) && s[j:j+2] == "${" {
			if buf == nil {
				buf = make([]byte, 0, 2*len(s))
			}
			if j-1 >= 0 && s[j-1] == '$' {
				buf = append(buf, s[i:j]...)
				j++
				i = j
				continue
			}
			buf = append(buf, s[i:j]...)
			name, w := getShellName(s[j+1:])
			if name == "" && w > 0 {
				// Encountered invalid syntax; eat the
				// characters.
			} else {
				buf = append(buf, mapping(name)...)
			}
			j += w
			i = j + 1
		}
	}
	if buf == nil {
		return s
	}
	return string(buf) + s[i:]
}

// isShellSpecialVar reports whether the character identifies a special
// shell variable such as $*.
func isShellSpecialVar(c uint8) bool {
	switch c {
	case '*', '#', '$', '@', '!', '?', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

// isAlphaNum reports whether the byte is an ASCII letter, number, or underscore
func isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

// getShellName returns the name that begins the string and the number of bytes
// consumed to extract it. Since the name is enclosed in {}, it's part of a ${}
// expansion and two more bytes are needed than the length of the name. If the
// internal syntax is un-env-iable (get it?), then just "eat" the variable.
func getShellName(s string) (string, int) {
	if len(s) > 2 && isShellSpecialVar(s[1]) && s[2] == '}' {
		return s[1:2], 3
	}
	// Scan to closing brace
	for i := 1; i < len(s); i++ {
		if s[i] == '}' {
			if i == 1 {
				return "", 2 // Bad syntax; eat "${}"
			}
			return s[1:i], i + 1
		}
	}
	return "", 1 // Bad syntax; eat "${"
}
