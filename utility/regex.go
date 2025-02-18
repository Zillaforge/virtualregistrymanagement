package utility

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	RegexStructTag      = "regroup"
	RegexStructRequired = "required"

	RegexCompileErrMsg                 errorMessage = "compilation error"
	RegexNoMatchFoundErrMsg            errorMessage = "no match found for given string"
	RegexNotStructPtrErrMsg            errorMessage = "expected struct pointer"
	RegexNoSubexpressionsFoundErrMsg   errorMessage = "no subexpressions found in regex"
	RegexUnknownGroupErrMsg            errorMessage = "group \"%s\" haven't found in regex"
	RegexTypeNotParsableErrMsg         errorMessage = "type \"%v\" is not parsable"
	RegexParseErrMsg                   errorMessage = "error parsing group \"%s\""
	RegexRequiredGroupIsEmptyErrMsg    errorMessage = "required regroup \"%s\" is empty for field \"%s\""
	RegexNilPointerInStructFieldErrMsg errorMessage = "can't set value to nil pointer in struct field: %s"
	RegexNilPointerInFieldErrMsg       errorMessage = "can't set value to nil pointer in field: %s"
)

type errorMessage string

type RegexError struct {
	typ    errorMessage
	params []interface{}
	err    error
}

func (r *RegexError) Error() string {
	msg := fmt.Sprintf(string(r.typ), r.params...)
	if r.err != nil {
		return fmt.Sprintf(fmt.Sprintf("%s: %s", msg, "%v"), r.err)
	}
	return msg
}

func (r *RegexError) Type() errorMessage {
	return r.typ
}

// IsRegexError error is from Regex
func IsRegexError(err error) (*RegexError, bool) {
	switch err := err.(type) {
	case *RegexError:
		return err, true
	}
	return nil, false
}

type Regex struct {
	*regexp.Regexp
}

// Compile compiles given expression as regex and return new ReGroup with this expression as matching engine.
// If the expression can't be compiled as regex, a CompileError will be returned
func Compile(str string) (*Regex, error) {
	regex, err := regexp.Compile(str)
	if err != nil {
		return nil, &RegexError{typ: RegexCompileErrMsg, err: err}
	}

	return &Regex{regex}, nil
}

// MustCompile calls Compile and panicked if it returned an error
func MustCompile(s string) *Regex {
	regex, _ := Compile(s)
	return regex
}

// Groups returns a map contains each group name as a key and the group's matched value as value
func (regex *Regex) Groups(s string) (map[string]string, error) {
	subexp := regex.SubexpNames()
	if len(subexp) == 1 {
		return nil, &RegexError{typ: RegexNoSubexpressionsFoundErrMsg}
	}

	match := regex.FindStringSubmatch(s)
	if match == nil {
		return nil, &RegexError{typ: RegexNoMatchFoundErrMsg}
	}

	results := make(map[string]string)
	for i, name := range subexp {
		if i != 0 && name != "" {
			results[name] = match[i]
		}
	}
	return results, nil
}

// MatchToTarget matches a regex expression to string s and parse it into `target` argument.
// If no matches found, a &NoMatchFoundError error will be returned
func (regex *Regex) MatchToTarget(s string, target interface{}) error {
	matched, err := regex.Groups(s)
	if err != nil {
		return err
	}

	targetRef, err := regex.validateTarget(target)
	if err != nil {
		return err
	}
	return regex.fillTarget(matched, targetRef)
}

// validateTarget checks that given interface is a pointer of struct
func (regex *Regex) validateTarget(target interface{}) (reflect.Value, error) {
	targetPtr := reflect.ValueOf(target)
	if targetPtr.Kind() != reflect.Ptr {
		return reflect.Value{}, &RegexError{typ: RegexNotStructPtrErrMsg}
	}
	return targetPtr.Elem(), nil
}

func (regex *Regex) fillTarget(matchGroup map[string]string, targetRef reflect.Value) error {
	targetType := targetRef.Type()
	for i := 0; i < targetType.NumField(); i++ {
		fieldRef := targetRef.Field(i)
		if !fieldRef.CanSet() {
			continue
		}

		if err := regex.setField(targetType.Field(i), fieldRef, matchGroup); err != nil {
			return err
		}
	}

	return nil
}

// setField getting a single struct field and matching groups map and set the field value to its matching group value tag
// after parsing it to match the field type
func (regex *Regex) setField(fieldType reflect.StructField, fieldRef reflect.Value, matchGroup map[string]string) error {
	fieldRefType := fieldType.Type
	ptr := false
	if fieldRefType.Kind() == reflect.Ptr {
		ptr = true
		fieldRefType = fieldType.Type.Elem()
	}

	if fieldRefType.Kind() == reflect.Struct {
		if ptr {
			if fieldRef.IsNil() {
				return &RegexError{typ: RegexNilPointerInStructFieldErrMsg, params: []interface{}{fieldType.Name}}
			}
			fieldRef = fieldRef.Elem()
		}
		return regex.fillTarget(matchGroup, fieldRef)
	}

	if ptr {
		if fieldRef.IsNil() {
			return &RegexError{typ: RegexNilPointerInFieldErrMsg, params: []interface{}{fieldType.Name}}
		}
		fieldRef = fieldRef.Elem()
	}

	regroupKey, regroupOption := regex.groupAndOption(fieldType)
	if regroupKey == "" {
		return nil
	}

	matchedVal, ok := matchGroup[regroupKey]
	if !ok {
		return &RegexError{typ: RegexUnknownGroupErrMsg, params: []interface{}{regroupKey}}
	}

	if matchedVal == "" {
		if RegexStructRequired == regroupOption {
			return &RegexError{typ: RegexRequiredGroupIsEmptyErrMsg, params: []interface{}{regroupKey, fieldType.Name}}
		}
		return nil
	}

	parsedFunc := getParsingFunc(fieldRefType)
	if parsedFunc == nil {
		return &RegexError{typ: RegexTypeNotParsableErrMsg, params: []interface{}{fieldRefType}}
	}

	parsed, err := parsedFunc(matchedVal, fieldRefType)
	if err != nil {
		return &RegexError{typ: RegexParseErrMsg, params: []interface{}{regroupKey}, err: err}
	}

	fieldRef.Set(parsed)

	return nil
}

func getParsingFunc(typ reflect.Type) func(src string, typ reflect.Type) (reflect.Value, error) {
	if typ == reflect.TypeOf(time.Second) {
		return func(src string, _ reflect.Type) (reflect.Value, error) {
			d, err := time.ParseDuration(src)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(d), nil
		}
	}

	switch typ.Kind() {
	case reflect.Bool:
		return func(src string, _ reflect.Type) (reflect.Value, error) {
			b, err := strconv.ParseBool(src)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(b), nil
		}

	case reflect.String:
		return func(src string, typ reflect.Type) (reflect.Value, error) {
			return reflect.ValueOf(src).Convert(typ), nil
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(src string, typ reflect.Type) (reflect.Value, error) {
			n, err := strconv.ParseInt(src, 10, 64)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(n).Convert(typ), nil
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(src string, typ reflect.Type) (reflect.Value, error) {
			n, err := strconv.ParseUint(src, 10, 64)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(n).Convert(typ), nil
		}

	case reflect.Float32, reflect.Float64:
		return func(src string, typ reflect.Type) (reflect.Value, error) {
			n, err := strconv.ParseFloat(src, 64)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(n).Convert(typ), nil
		}
	}
	return nil
}

// groupAndOption return the requested regroup and it's option splitted by ','
func (r *Regex) groupAndOption(fieldType reflect.StructField) (group, option string) {
	regroupKey := fieldType.Tag.Get(RegexStructTag)
	if regroupKey == "" {
		return "", ""
	}
	split := strings.Split(regroupKey, ",")
	if len(split) == 1 {
		return strings.TrimSpace(split[0]), ""
	}
	return strings.TrimSpace(split[0]), strings.TrimSpace(strings.ToLower(split[1]))
}
