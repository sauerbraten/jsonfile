package jsonconf

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
)

type configFileFilter struct {
	configFile io.ByteReader
	pos        int  // current position in the file
	inString   bool // true if current position points inside a string; used to only strip whitespace outside of strings
}

func newConfigFileFilter(fileName string) (cff *configFileFilter, err error) {
	var file *os.File
	file, err = os.Open(fileName)

	cff = &configFileFilter{configFile: bufio.NewReader(file)}

	return
}

// Reads from the config file and strips whitespace outside of strings as well as comments. With this method, configFileTilter implements io.Reader.
func (cff *configFileFilter) Read(p []byte) (n int, err error) {
	// temporarily save current position
	i := cff.pos
	var b, c byte

	for i < len(p) {
		b, err = cff.configFile.ReadByte()
		if err != nil {
			return
		}

		if cff.inString {
			// use byte as-is
			p[n] = b
			n++

			// check if this is the end of the string
			if rune(b) == '"' {
				cff.inString = !cff.inString
			}
		} else {
			switch rune(b) {
			case '/':
				// this is a comment, next byte has to be '/' as well, else it's invalid JSON
				c, err = cff.configFile.ReadByte()
				if err != nil {
					return
				}

				if c == '/' {
					// skip until new line
					for c, err = cff.configFile.ReadByte(); err == nil && rune(c) != '\n'; c, err = cff.configFile.ReadByte() {
						i++
					}
					if err != nil {
						return
					}
				} else {
					// '/' is an illegal character in JSON outside of a string
					err = errors.New("illegal character '/' outside of string")
					return
				}

			case ' ', '\t', '\n':
				// skip whitespace

			case '"':
				// use byte
				p[n] = b
				n++
				// entering a string
				cff.inString = !cff.inString

			default:
				// use byte as-is
				p[n] = b
				n++
			}
		}

		i++
	}

	// advance position in file stream
	cff.pos += i

	// check if we're at the end of the file
	if i < len(p) {
		err = io.EOF
	}

	return
}

// Parses a config file at fileName into the provided interface, which must be of a pointer type.
func ParseFile(fileName string, v interface{}) (err error) {
	cff, err := newConfigFileFilter(fileName)
	if err != nil {
		return
	}

	// read filtered JSON and unmarshal it into the provided interface
	err = json.NewDecoder(cff).Decode(v)
	return
}
