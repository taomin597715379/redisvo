package splitecommand

import (
	"errors"
	"fmt"
)

func AnalyzeCommand(commandLine string) ([]string, error) {
	var pos = 0
	var ret []string
	var length = len(commandLine)
	if commandLine == `` {
		return []string{}, nil
	}
	for {
		for pos < length && Isspace(string(commandLine[pos])) {
			pos++
		}
		if pos == length {
			break
		}
		var inQuotes = false
		var inSingQuotes = false
		var done = false
		var current string = ``
		for !done && pos != length {
			var c = string(commandLine[pos])
			if inQuotes {
				if c == "\\" && (pos+1) < length {
					pos += 1
					switch string(commandLine[pos]) {
					case "n":
						c = "\n"
						break
					case "r":
						c = "\r"
						break
					case "t":
						c = "\t"
						break
					case "b":
						c = "\b"
						break
					case "a":
						c = "\a"
						break
					default:
						c = string(commandLine[pos])
						break
					}
					current += c
				} else if c == `"` {
					if (pos+1) < length && !Isspace(string(commandLine[pos+1])) {
						return []string{}, errors.New(fmt.Sprintf("Expect \" followed by a space or nothing, got %s ", string(commandLine[pos+1])))
					}
					done = true
				} else if pos == length {
					return []string{}, errors.New(`Unterminated quotes... `)
				} else {
					current += c
				}
			} else if inSingQuotes {
				if c == string("\\") && string(commandLine[pos+1]) == `'` {
					pos += 1
					current += `'`
				} else if c == `'` {
					if (pos+1) < length && !Isspace(string(commandLine[pos+1])) {
						return []string{}, errors.New(fmt.Sprintf("Expect \\' followed by a space or nothing, got %s ", string(commandLine[pos+1])))
					}
					done = true
				} else if pos == length {
					return []string{}, errors.New(`Unterminated quotes... `)
				} else {
					current += c
				}
			} else {
				if pos == length {
					done = true
				} else {
					switch c {
					case " ", "\n", "\r", "\t":
						done = true
					case `"`:
						inQuotes = true
					case `'`:
						inSingQuotes = true
					default:
						current += c
					}
				}
			}
			if pos < length {
				pos++
			}
		}
		ret = append(ret, current)
	}
	return ret, nil
}

func Isspace(ch string) bool {
	return (ch == " ") || (ch == "\t") || (ch == "\n") || (ch == "\r") || (ch == "\v") || (ch == "\f")
}
