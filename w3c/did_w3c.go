// Package w3c is a set of tools to work with Decentralized Identifiers (DIDs)
// as described in the DID spec https://w3c.github.io/did-core/
// Got from https://github.com/build-trust/did
package w3c

import (
	"fmt"
	"strings"
)

// Param represents a parsed DID param,
// which contains a name and value. A generic param is defined
// as a param name and value separated by a colon.
// generic-param-name:param-value
// A param may also be method specific, which
// requires the method name to prefix the param name separated by a colon
// method-name:param-name.
// param = param-name [ "=" param-value ]
// https://w3c.github.io/did-core/#generic-did-parameter-names
// https://w3c.github.io/did-core/#method-specific-did-parameter-names
type Param struct {
	// param-name = 1*param-char
	// Name may include a method name and param name separated by a colon
	Name string
	// param-value = *param-char
	Value string
}

// String encodes a Param struct into a valid Param string.
// Name is required by the grammar. Value is optional
func (p *Param) String() string {
	if p.Name == "" {
		return ""
	}

	if 0 < len(p.Value) {
		return p.Name + "=" + p.Value
	}

	return p.Name
}

// A DID represents a parsed DID or a DID URL
type DID struct {
	// DID Method
	// https://w3c.github.io/did-core/#method-specific-syntax
	Method string

	// The method-specific-id component of a DID
	// method-specific-id = *idchar *( ":" *idchar )
	ID string

	// method-specific-id may be composed of multiple `:` separated idstrings
	IDStrings []string

	// DID URL
	// did-url = did *( ";" param ) path-abempty [ "?" query ] [ "#" fragment ]
	// did-url may contain multiple params, a path, query, and fragment
	Params []Param

	// DID Path, the portion of a DID reference that follows the first forward slash character.
	// https://w3c.github.io/did-core/#path
	Path string

	// Path may be composed of multiple `/` separated segments
	// path-abempty  = *( "/" segment )
	PathSegments []string

	// DID Query
	// https://w3c.github.io/did-core/#query
	// query = *( pchar / "/" / "?" )
	Query string

	// DID Fragment, the portion of a DID reference that follows the first hash sign character ("#")
	// https://w3c.github.io/did-core/#fragment
	Fragment string
}

// the parsers internal state
type parser struct {
	input        string // input to the parser
	currentIndex int    // index in the input which the parser is currently processing
	out          *DID   // the output DID that the parser will assemble as it steps through its state machine
	err          error  // an error in the parser state machine
}

// a step in the parser state machine that returns the next step
type parserStep func() parserStep

// IsURL returns true if a DID has a Path, a Query or a Fragment
// https://w3c-ccg.github.io/did-spec/#dfn-did-reference
func (d *DID) IsURL() bool {
	return (len(d.Params) > 0 || d.Path != "" || len(d.PathSegments) > 0 || d.Query != "" || d.Fragment != "")
}

// String encodes a DID struct into a valid DID string.
// nolint: gocyclo
func (d *DID) String() string {
	var buf strings.Builder

	// write the did: prefix
	buf.WriteString("did:") // nolint, returned error is always nil

	if d.Method != "" {
		// write method followed by a `:`
		buf.WriteString(d.Method) // nolint, returned error is always nil
		buf.WriteByte(':')        // nolint, returned error is always nil
	} else {
		// if there is no Method, return an empty string
		return ""
	}

	if d.ID != "" {
		buf.WriteString(d.ID) // nolint, returned error is always nil
	} else if len(d.IDStrings) > 0 {
		// join IDStrings with a colon to make the ID
		buf.WriteString(strings.Join(d.IDStrings[:], ":")) // nolint, returned error is always nil
	} else {
		// if there is no ID, return an empty string
		return ""
	}

	if len(d.Params) > 0 {
		// write a leading ; for each param
		for _, p := range d.Params {
			// get a string that represents the param
			param := p.String()
			if param != "" {
				// params must start with a ;
				buf.WriteByte(';')     // nolint, returned error is always nil
				buf.WriteString(param) // nolint, returned error is always nil
			} else {
				// if a param exists but is empty, return an empty string
				return ""
			}
		}
	}

	if d.Path != "" {
		// write a leading / and then Path
		buf.WriteByte('/')      // nolint, returned error is always nil
		buf.WriteString(d.Path) // nolint, returned error is always nil
	} else if len(d.PathSegments) > 0 {
		// write a leading / and then PathSegments joined with / between them
		buf.WriteByte('/')                                    // nolint, returned error is always nil
		buf.WriteString(strings.Join(d.PathSegments[:], "/")) // nolint, returned error is always nil
	}

	if d.Query != "" {
		// write a leading ? and then Query
		buf.WriteByte('?')       // nolint, returned error is always nil
		buf.WriteString(d.Query) // nolint, returned error is always nil
	}

	if d.Fragment != "" {
		// add fragment only when there is no path
		buf.WriteByte('#')          // nolint, returned error is always nil
		buf.WriteString(d.Fragment) // nolint, returned error is always nil
	}

	return buf.String()
}

// ParseDID parses the input string into a DID structure.
func ParseDID(input string) (*DID, error) {
	// initialize the parser state
	p := &parser{input: input, out: &DID{}}

	// the parser state machine is implemented as a loop over parser steps
	// steps increment p.currentIndex as they consume the input, each step returns the next step to run
	// the state machine halts when one of the steps returns nil
	//
	// This design is based on this talk from Rob Pike, although the talk focuses on lexical scanning,
	// the DID grammar is simple enough for us to combine lexing and parsing into one lexerless parse
	// http://www.youtube.com/watch?v=HxaD_trXwRE
	parserState := p.checkLength
	for parserState != nil {
		parserState = parserState()
	}

	// If one of the steps added an err to the parser state, exit. Return nil and the error.
	err := p.err
	if err != nil {
		return nil, err
	}

	// join IDStrings with : to make up ID
	p.out.ID = strings.Join(p.out.IDStrings[:], ":")

	// join PathSegments with / to make up Path
	p.out.Path = strings.Join(p.out.PathSegments[:], "/")

	return p.out, nil
}

// checkLength is a parserStep that checks if the input length is atleast 7
// the grammar requires
//
//	`did:` prefix (4 chars)
//	+ atleast one methodchar (1 char)
//	+ `:` (1 char)
//	+ atleast one idchar (1 char)
//
// i.e. at least 7 chars
// The current specification does not take a position on maximum length of a DID.
// https://w3c-ccg.github.io/did-spec/#upper-limits-on-did-character-length
func (p *parser) checkLength() parserStep {
	inputLength := len(p.input)

	if inputLength < 7 {
		return p.errorf(inputLength, "input length is less than 7")
	}

	return p.parseScheme
}

// parseScheme is a parserStep that validates that the input begins with 'did:'
func (p *parser) parseScheme() parserStep {

	currentIndex := 3 // 4 bytes in 'did:', i.e index 3

	// the grammar requires `did:` prefix
	if p.input[:currentIndex+1] != "did:" {
		return p.errorf(currentIndex, "input does not begin with 'did:' prefix")
	}

	p.currentIndex = currentIndex
	return p.parseMethod
}

// parseMethod is a parserStep that extracts the DID Method
// from the grammar:
//
//	did        = "did:" method ":" specific-idstring
//	method     = 1*methodchar
//	methodchar = %x61-7A / DIGIT ; 61-7A is a-z in US-ASCII
func (p *parser) parseMethod() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1
	startIndex := currentIndex

	// parse method name
	// loop over every byte following the ':' in 'did:' unlil the second ':'
	// method is the string between the two ':'s
	for {
		if currentIndex == inputLength {
			// we got to the end of the input and didn't find a second ':'
			return p.errorf(currentIndex, "input does not have a second `:` marking end of method name")
		}

		// read the input character at currentIndex
		char := input[currentIndex]

		if char == ':' {
			// we've found the second : in the input that marks the end of the method
			if currentIndex == startIndex {
				// return error is method is empty, ex- did::1234
				return p.errorf(currentIndex, "method is empty")
			}
			break
		}

		// as per the grammar method can only be made of digits 0-9 or small letters a-z
		if isNotDigit(char) && isNotSmallLetter(char) {
			return p.errorf(currentIndex, "character is not a-z OR 0-9")
		}

		// move to the next char
		currentIndex = currentIndex + 1
	}

	// set parser state
	p.currentIndex = currentIndex
	p.out.Method = input[startIndex:currentIndex]

	// method is followed by specific-idstring, parse that next
	return p.parseID
}

// parseID is a parserStep that extracts : separated idstrings that are part of a specific-idstring
// and adds them to p.out.IDStrings
// from the grammar:
//
//	specific-idstring = idstring *( ":" idstring )
//	idstring          = 1*idchar
//	idchar            = ALPHA / DIGIT / "." / "-"
//
// p.out.IDStrings is later concatented by the ParseDID function before it returns.
func (p *parser) parseID() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1
	startIndex := currentIndex

	var next parserStep

	for {
		if currentIndex == inputLength {
			// we've reached end of input, no next state
			next = nil
			break
		}

		char := input[currentIndex]

		if char == ':' {
			// encountered : input may have another idstring, parse ID again
			next = p.parseID
			break
		}

		if char == ';' {
			// encountered ; input may have a parameter, parse that next
			next = p.parseParamName
			break
		}

		if char == '/' {
			// encountered / input may have a path following specific-idstring, parse that next
			next = p.parsePath
			break
		}

		if char == '?' {
			// encountered ? input may have a query following specific-idstring, parse that next
			next = p.parseQuery
			break
		}

		if char == '#' {
			// encountered # input may have a fragment following specific-idstring, parse that next
			next = p.parseFragment
			break
		}

		// make sure current char is a valid idchar
		// idchar = ALPHA / DIGIT / "." / "-"
		if isNotValidIDChar(char) {
			return p.errorf(currentIndex, "byte is not ALPHA OR DIGIT OR '.' OR '-'")
		}

		// move to the next char
		currentIndex = currentIndex + 1
	}

	if currentIndex == startIndex {
		// idstring length is zero
		// from the grammar:
		//   idstring = 1*idchar
		// return error because idstring is empty, ex- did:a::123:456
		return p.errorf(currentIndex, "idstring must be atleast one char long")
	}

	// set parser state
	p.currentIndex = currentIndex
	p.out.IDStrings = append(p.out.IDStrings, input[startIndex:currentIndex])

	// return the next parser step
	return next
}

// parseParamName is a parserStep that extracts a did-url param-name.
// A Param struct is created for each param name that is encountered.
// from the grammar:
//
//	param              = param-name [ "=" param-value ]
//	param-name         = 1*param-char
//	param-char         = ALPHA / DIGIT / "." / "-" / "_" / ":" / pct-encoded
func (p *parser) parseParamName() parserStep {
	input := p.input
	startIndex := p.currentIndex + 1
	next := p.paramTransition()
	currentIndex := p.currentIndex

	if currentIndex == startIndex {
		// param-name length is zero
		// from the grammar:
		//   1*param-char
		// return error because param-name is empty, ex- did:a::123:456;param-name
		return p.errorf(currentIndex, "Param name must be at least one char long")
	}

	// Create a new param with the name
	p.out.Params = append(p.out.Params, Param{Name: input[startIndex:currentIndex], Value: ""})

	// return the next parser step
	return next
}

// parseParamValue is a parserStep that extracts a did-url param-value.
// A parsed Param value requires that a Param was previously created when parsing a param-name.
// from the grammar:
//
//	param              = param-name [ "=" param-value ]
//	param-value         = 1*param-char
//	param-char         = ALPHA / DIGIT / "." / "-" / "_" / ":" / pct-encoded
func (p *parser) parseParamValue() parserStep {
	input := p.input
	startIndex := p.currentIndex + 1
	next := p.paramTransition()
	currentIndex := p.currentIndex

	// Get the last Param in the DID and append the value
	// values may be empty according to the grammar- *param-char
	p.out.Params[len(p.out.Params)-1].Value = input[startIndex:currentIndex]

	// return the next parser step
	return next
}

// paramTransition is a parserStep that extracts and transitions a param-name or
// param-value.
// nolint: gocyclo
func (p *parser) paramTransition() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1

	var indexIncrement int
	var next parserStep
	var percentEncoded bool

	for {
		if currentIndex == inputLength {
			// we've reached end of input, no next state
			next = nil
			break
		}

		char := input[currentIndex]

		if char == ';' {
			// encountered : input may have another param, parse paramName again
			next = p.parseParamName
			break
		}

		// Separate steps for name and value?
		if char == '=' {
			// parse param value
			next = p.parseParamValue
			break
		}

		if char == '/' {
			// encountered / input may have a path following current param, parse that next
			next = p.parsePath
			break
		}

		if char == '?' {
			// encountered ? input may have a query following current param, parse that next
			next = p.parseQuery
			break
		}

		if char == '#' {
			// encountered # input may have a fragment following current param, parse that next
			next = p.parseFragment
			break
		}

		if char == '%' {
			// a % must be followed by 2 hex digits
			if (currentIndex+2 >= inputLength) ||
				isNotHexDigit(input[currentIndex+1]) ||
				isNotHexDigit(input[currentIndex+2]) {
				return p.errorf(currentIndex, "%% is not followed by 2 hex digits")
			}
			// if we got here, we're dealing with percent encoded char, jump three chars
			percentEncoded = true
			indexIncrement = 3
		} else {
			// not percent encoded
			percentEncoded = false
			indexIncrement = 1
		}

		// make sure current char is a valid param-char
		// idchar = ALPHA / DIGIT / "." / "-"
		if !percentEncoded && isNotValidParamChar(char) {
			return p.errorf(currentIndex, "character is not allowed in param - %c", char)
		}

		// move to the next char
		currentIndex = currentIndex + indexIncrement
	}

	// set parser state
	p.currentIndex = currentIndex

	return next
}

// parsePath is a parserStep that extracts a DID Path from a DID Reference
// from the grammar:
//
//	did-path      = segment-nz *( "/" segment )
//	segment       = *pchar
//	segment-nz    = 1*pchar
//	pchar         = unreserved / pct-encoded / sub-delims / ":" / "@"
//	unreserved    = ALPHA / DIGIT / "-" / "." / "_" / "~"
//	pct-encoded   = "%" HEXDIG HEXDIG
//	sub-delims    = "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
//
// nolint: gocyclo
func (p *parser) parsePath() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1
	startIndex := currentIndex

	var indexIncrement int
	var next parserStep
	var percentEncoded bool

	for {
		if currentIndex == inputLength {
			next = nil
			break
		}

		char := input[currentIndex]

		if char == '/' {
			// encountered / input may have another path segment, try to parse that next
			next = p.parsePath
			break
		}

		if char == '?' {
			// encountered ? input may have a query following path, parse that next
			next = p.parseQuery
			break
		}

		if char == '%' {
			// a % must be followed by 2 hex digits
			if (currentIndex+2 >= inputLength) ||
				isNotHexDigit(input[currentIndex+1]) ||
				isNotHexDigit(input[currentIndex+2]) {
				return p.errorf(currentIndex, "%% is not followed by 2 hex digits")
			}
			// if we got here, we're dealing with percent encoded char, jump three chars
			percentEncoded = true
			indexIncrement = 3
		} else {
			// not pecent encoded
			percentEncoded = false
			indexIncrement = 1
		}

		// pchar = unreserved / pct-encoded / sub-delims / ":" / "@"
		if !percentEncoded && isNotValidPathChar(char) {
			return p.errorf(currentIndex, "character is not allowed in path")
		}

		// move to the next char
		currentIndex = currentIndex + indexIncrement
	}

	if currentIndex == startIndex && len(p.out.PathSegments) == 0 {
		// path segment length is zero
		// first path segment must have atleast one character
		// from the grammar
		//   did-path = segment-nz *( "/" segment )
		return p.errorf(currentIndex, "first path segment must have atleast one character")
	}

	// update parser state
	p.currentIndex = currentIndex
	p.out.PathSegments = append(p.out.PathSegments, input[startIndex:currentIndex])

	return next
}

// parseQuery is a parserStep that extracts a DID Query from a DID Reference
// from the grammar:
//
//	did-query     = *( pchar / "/" / "?" )
//	pchar         = unreserved / pct-encoded / sub-delims / ":" / "@"
//	unreserved    = ALPHA / DIGIT / "-" / "." / "_" / "~"
//	pct-encoded   = "%" HEXDIG HEXDIG
//	sub-delims    = "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
func (p *parser) parseQuery() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1
	startIndex := currentIndex

	var indexIncrement int
	var next parserStep
	var percentEncoded bool

	for {
		if currentIndex == inputLength {
			// we've reached the end of input
			// it's ok for query to be empty, so we don't need a check for that
			// did-query     = *( pchar / "/" / "?" )
			break
		}

		char := input[currentIndex]

		if char == '#' {
			// encountered # input may have a fragment following the query, parse that next
			next = p.parseFragment
			break
		}

		if char == '%' {
			// a % must be followed by 2 hex digits
			if (currentIndex+2 >= inputLength) ||
				isNotHexDigit(input[currentIndex+1]) ||
				isNotHexDigit(input[currentIndex+2]) {
				return p.errorf(currentIndex, "%% is not followed by 2 hex digits")
			}
			// if we got here, we're dealing with percent encoded char, jump three chars
			percentEncoded = true
			indexIncrement = 3
		} else {
			// not pecent encoded
			percentEncoded = false
			indexIncrement = 1
		}

		// did-query = *( pchar / "/" / "?" )
		// pchar = unreserved / pct-encoded / sub-delims / ":" / "@"
		// isNotValidQueryOrFragmentChar checks for all the valid chars except pct-encoded
		if !percentEncoded && isNotValidQueryOrFragmentChar(char) {
			return p.errorf(currentIndex, "character is not allowed in query - %c", char)
		}

		// move to the next char
		currentIndex = currentIndex + indexIncrement
	}

	// update parser state
	p.currentIndex = currentIndex
	p.out.Query = input[startIndex:currentIndex]

	return next
}

// parseFragment is a parserStep that extracts a DID Fragment from a DID Reference
// from the grammar:
//
//	did-fragment  = *( pchar / "/" / "?" )
//	pchar         = unreserved / pct-encoded / sub-delims / ":" / "@"
//	unreserved    = ALPHA / DIGIT / "-" / "." / "_" / "~"
//	pct-encoded   = "%" HEXDIG HEXDIG
//	sub-delims    = "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
func (p *parser) parseFragment() parserStep {
	input := p.input
	inputLength := len(input)
	currentIndex := p.currentIndex + 1
	startIndex := currentIndex

	var indexIncrement int
	var percentEncoded bool

	for {
		if currentIndex == inputLength {
			// we've reached the end of input
			// it's ok for reference to be empty, so we don't need a check for that
			// did-fragment = *( pchar / "/" / "?" )
			break
		}

		char := input[currentIndex]

		if char == '%' {
			// a % must be followed by 2 hex digits
			if (currentIndex+2 >= inputLength) ||
				isNotHexDigit(input[currentIndex+1]) ||
				isNotHexDigit(input[currentIndex+2]) {
				return p.errorf(currentIndex, "%% is not followed by 2 hex digits")
			}
			// if we got here, we're dealing with percent encoded char, jump three chars
			percentEncoded = true
			indexIncrement = 3
		} else {
			// not pecent encoded
			percentEncoded = false
			indexIncrement = 1
		}

		// did-fragment = *( pchar / "/" / "?" )
		// pchar = unreserved / pct-encoded / sub-delims / ":" / "@"
		// isNotValidQueryOrFragmentChar checks for all the valid chars except pct-encoded
		if !percentEncoded && isNotValidQueryOrFragmentChar(char) {
			return p.errorf(currentIndex, "character is not allowed in fragment - %c", char)
		}

		// move to the next char
		currentIndex = currentIndex + indexIncrement
	}

	// update parser state
	p.currentIndex = currentIndex
	p.out.Fragment = input[startIndex:currentIndex]

	// no more parsing needed after a fragment,
	// cause the state machine to exit by returning nil
	return nil
}

// errorf is a parserStep that returns nil to cause the state machine to exit
// before returning it sets the currentIndex and err field in parser state
// other parser steps use this function to exit the state machine with an error
func (p *parser) errorf(index int, format string, args ...interface{}) parserStep {
	p.currentIndex = index
	p.err = fmt.Errorf(format, args...)
	return nil
}

// INLINABLE
// Calls to all functions below this point should be inlined by the go compiler
// See output of `go build -gcflags -m` to confirm

// isNotValidIDChar returns true if a byte is not allowed in a ID
// from the grammar:
//
//	idchar = ALPHA / DIGIT / "." / "-"
func isNotValidIDChar(char byte) bool {
	return isNotAlpha(char) && isNotDigit(char) && char != '.' && char != '-'
}

// isNotValidParamChar returns true if a byte is not allowed in a param-name
// or param-value from the grammar:
//
//	idchar = ALPHA / DIGIT / "." / "-" / "_" / ":"
func isNotValidParamChar(char byte) bool {
	return isNotAlpha(char) && isNotDigit(char) &&
		char != '.' && char != '-' && char != '_' && char != ':'
}

// isNotValidQueryOrFragmentChar returns true if a byte is not allowed in a Fragment
// from the grammar:
//
//	did-fragment = *( pchar / "/" / "?" )
//	pchar        = unreserved / pct-encoded / sub-delims / ":" / "@"
//
// pct-encoded is not checked in this function
func isNotValidQueryOrFragmentChar(char byte) bool {
	return isNotValidPathChar(char) && char != '/' && char != '?'
}

// isNotValidPathChar returns true if a byte is not allowed in Path
//
//	did-path    = segment-nz *( "/" segment )
//	segment     = *pchar
//	segment-nz  = 1*pchar
//	pchar       = unreserved / pct-encoded / sub-delims / ":" / "@"
//
// pct-encoded is not checked in this function
func isNotValidPathChar(char byte) bool {
	return isNotUnreservedOrSubdelim(char) && char != ':' && char != '@'
}

// isNotUnreservedOrSubdelim returns true if a byte is not unreserved or sub-delims
// from the grammar:
//
//	unreserved = ALPHA / DIGIT / "-" / "." / "_" / "~"
//	sub-delims = "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
//
// https://tools.ietf.org/html/rfc3986#appendix-A
func isNotUnreservedOrSubdelim(char byte) bool {
	switch char {
	case '-', '.', '_', '~', '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=':
		return false
	default:
		if isNotAlpha(char) && isNotDigit(char) {
			return true
		}
		return false
	}
}

// isNotHexDigit returns true if a byte is not a digit between 0-9 or A-F or a-f
// in US-ASCII http://www.columbia.edu/kermit/ascii.html
// https://tools.ietf.org/html/rfc5234#appendix-B.1
func isNotHexDigit(char byte) bool {
	// '\x41' is A, '\x46' is F
	// '\x61' is a, '\x66' is f
	return isNotDigit(char) && (char < '\x41' || char > '\x46') && (char < '\x61' || char > '\x66')
}

// isNotDigit returns true if a byte is not a digit between 0-9
// in US-ASCII http://www.columbia.edu/kermit/ascii.html
// https://tools.ietf.org/html/rfc5234#appendix-B.1
func isNotDigit(char byte) bool {
	// '\x30' is digit 0, '\x39' is digit 9
	return (char < '\x30' || char > '\x39')
}

// isNotAlpha returns true if a byte is not a big letter between A-Z or small letter between a-z
// https://tools.ietf.org/html/rfc5234#appendix-B.1
func isNotAlpha(char byte) bool {
	return isNotSmallLetter(char) && isNotBigLetter(char)
}

// isNotBigLetter returns true if a byte is not a big letter between A-Z
// in US-ASCII http://www.columbia.edu/kermit/ascii.html
// https://tools.ietf.org/html/rfc5234#appendix-B.1
func isNotBigLetter(char byte) bool {
	// '\x41' is big letter A, '\x5A' small letter Z
	return (char < '\x41' || char > '\x5A')
}

// isNotSmallLetter returns true if a byte is not a small letter between a-z
// in US-ASCII http://www.columbia.edu/kermit/ascii.html
// https://tools.ietf.org/html/rfc5234#appendix-B.1
func isNotSmallLetter(char byte) bool {
	// '\x61' is small letter a, '\x7A' small letter z
	return (char < '\x61' || char > '\x7A')
}
