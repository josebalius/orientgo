package obinary

import (
	"fmt"
	"github.com/istreamdata/orientgo/oschema"
	"math"
	"strconv"
	"strings"
)

const (
	string_LINK                 = '#'
	string_EMBEDDED_BEGIN       = '('
	string_EMBEDDED_END         = ')'
	string_LIST_BEGIN           = '['
	string_LIST_END             = ']'
	string_SET_BEGIN            = '<'
	string_SET_END              = '>'
	string_MAP_BEGIN            = '{'
	string_MAP_END              = '}'
	string_BAG_BEGIN            = '%'
	string_BAG_END              = ';'
	string_BINARY_BEGINEND      = '_'
	string_CUSTOM_TYPE          = '^'
	string_ENTRY_SEPARATOR      = ':'
	string_PARAMETER_NAMED      = ':'
	string_PARAMETER_POSITIONAL = '?'

	string_DECIMAL_SEPARATOR = '.'
)

var (
	string_MaxInt = strconv.Itoa(math.MaxInt32)
)

type stringRecordFormatAbs struct{}

func (stringRecordFormatAbs) GetType(s string) oschema.OType {
	if s == "" {
		return oschema.UNKNOWN
	}
	rs := []rune(s)
	firstChar := rs[0]
	switch firstChar {
	case string_LINK: // RID
		return oschema.LINK
	case '\'', '"':
		return oschema.STRING
	case string_BINARY_BEGINEND:
		return oschema.BINARY
	case string_EMBEDDED_BEGIN:
		return oschema.EMBEDDED
	case string_LIST_BEGIN:
		return oschema.EMBEDDEDLIST
	case string_SET_BEGIN:
		return oschema.EMBEDDEDSET
	case string_MAP_BEGIN:
		return oschema.EMBEDDEDMAP
	case string_CUSTOM_TYPE:
		return oschema.CUSTOM
	}

	// BOOLEAN?
	if ls := strings.ToLower(s); ls == "true" || ls == "false" {
		return oschema.BOOLEAN
	}

	// NUMBER OR STRING?
	integer := true
	for i, c := range rs {
		if c >= '0' && c <= '9' {
			continue
		} else if i == 0 && (c == '+' || c == '0') {
			continue
		} else if c == string_DECIMAL_SEPARATOR {
			integer = false // maybe float, seek for other string char to be sure
		} else {
			if i == 0 {
				return oschema.STRING
			}
			if !integer && (c == 'E' || c == 'e') {
				// CHECK FOR SCIENTIFIC NOTATION
				if i+1 < len(rs) {
					if rs[i+1] == '-' {
						// JUMP THE DASH IF ANY (NOT MANDATORY)
						i++
					}
					continue
				}
			} else {
				switch c {
				case 'f':
					return oschema.FLOAT
				case 'c':
					return oschema.DECIMAL
				case 'l':
					return oschema.LONG
				case 'd':
					return oschema.DOUBLE
				case 'b':
					return oschema.BYTE
				case 'a':
					return oschema.DATE
				case 't':
					return oschema.DATETIME
				case 's':
					return oschema.SHORT
				}
			}
			return oschema.STRING
		}
	}

	if integer {
		// AUTO CONVERT TO LONG IF THE INTEGER IS TOO BIG
		if n, mn := len(rs), len(string_MaxInt); n > mn || (n == mn && s > string_MaxInt) {
			return oschema.LONG
		}
		return oschema.INTEGER
	}

	if _, err := strconv.ParseFloat(s, 32); err == nil {
		return oschema.FLOAT
	} else if _, err = strconv.ParseFloat(s, 64); err == nil {
		return oschema.DOUBLE
	} else {
		return oschema.DECIMAL
	}
}
func (f stringRecordFormatAbs) FieldTypeFromStream(tp oschema.OType, s string) interface{} {
	if s == "" {
		return nil
	} else if tp == oschema.UNKNOWN {
		tp = oschema.EMBEDDED
	}

	switch tp {
	case oschema.STRING:
		return s // TODO: implement in a right way
	case oschema.INTEGER:
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			panic(err)
		}
		return int32(v)
	case oschema.LONG:
		v, err := strconv.ParseInt(strings.TrimSuffix(s, "l"), 10, 64)
		if err != nil {
			panic(err)
		}
		return int64(v)
	case oschema.BOOLEAN:
		switch strings.ToLower(s) {
		case "true":
			return true
		case "false":
			return false
		default:
			panic(fmt.Errorf("unknown val for bool: '%s'", s))
		}
	default: // TODO: more types
		panic(fmt.Errorf("unsupported type for stringRecordFormatAbs: %s", tp))
	}
}
