package main

// splitCommandLine splits a Windows command line into its arguments, by
// the rules the C runtime and CommandLineToArgvW use, which are not a
// plain split on spaces:
//
//   - space and tab separate arguments;
//   - a quote opens or closes a quoted run, in which separators are
//     literal;
//   - two quotes inside a quoted run are one literal quote;
//   - backslashes before a quote are literal, except that every pair
//     escapes one backslash and an odd one escapes the quote.
//
// It returns nil for an empty command line.
func splitCommandLine(line string) []string {
	var (
		args  []string
		cur   []rune
		arg   bool // an argument is being built, even if empty
		quote bool
	)
	runes := []rune(line)
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		switch {
		case c == '\\':
			// Count the run of backslashes: they are literal unless they
			// escape a quote.
			n := 0
			for i < len(runes) && runes[i] == '\\' {
				n++
				i++
			}
			if i < len(runes) && runes[i] == '"' {
				cur = append(cur, repeat('\\', n/2)...)
				if n%2 == 0 {
					quote = !quote
					arg = true
				} else {
					cur = append(cur, '"')
					arg = true
				}
				continue
			}
			cur = append(cur, repeat('\\', n)...)
			arg = true
			i--
		case c == '"':
			// A doubled quote inside a quoted run is one literal quote.
			if quote && i+1 < len(runes) && runes[i+1] == '"' {
				cur = append(cur, '"')
				i++
			} else {
				quote = !quote
			}
			arg = true
		case (c == ' ' || c == '\t') && !quote:
			if arg {
				args = append(args, string(cur))
				cur, arg = cur[:0], false
			}
		default:
			cur = append(cur, c)
			arg = true
		}
	}
	if arg {
		args = append(args, string(cur))
	}
	return args
}

func repeat(r rune, n int) []rune {
	if n <= 0 {
		return nil
	}
	out := make([]rune, n)
	for i := range out {
		out[i] = r
	}
	return out
}
