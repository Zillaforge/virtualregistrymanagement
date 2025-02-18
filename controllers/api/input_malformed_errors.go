package api

import (
	cnt "VirtualRegistryManagement/constants"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
)

var (
	bindingErrorList map[string]string
	_binding         = "binding"
	_timeFormat      = "timeformat"
	format           string
	timeFormatsMap   = map[string]string{
		"ansic":       time.ANSIC,
		"unixdate":    time.UnixDate,
		"rubydate":    time.RubyDate,
		"rfc822":      time.RFC822,
		"rfc822Z":     time.RFC822Z,
		"rfc850":      time.RFC850,
		"rfc1123":     time.RFC1123,
		"rfc1123Z":    time.RFC1123Z,
		"rfc3339":     time.RFC3339,
		"rfc3339Nano": time.RFC3339Nano,
		"kitchen":     time.Kitchen,
		"stamp":       time.Stamp,
		"stampmilli":  time.StampMilli,
		"stampmicro":  time.StampMicro,
		"stampnano":   time.StampNano,
	}
)

func init() {
	bindingErrorList = map[string]string{
		"EOF": "invalid content",

		"TagID:required":  "tag-id is required",
		"UserID:required": "user-id is required",

		"Name:required":            "name is required",
		"OperatingSystem:required": "operatingSystem is required",
		"OperatingSystem:oneof":    "type valid values are linux, and windows",
		"Type:required":            "type is required",
		"DiskFormat:required":      "diskFormat is required",
		"ContainerFormat:required": "containerFormat is required",
		"Type:oneof":               "type valid values are common, and increase",
		"DiskFormat:oneof":         "diskFormat valid values are ami, ari, aki, vhd, vmdk, raw, qcow2, vdi, and iso",
		"ContainerFormat:oneof":    "containerFormat valid values are ami, ari, aki, bare, and ovf",

		"Filepath:required":                    "filepath is required",
		"RepositoryID:required_without_all":    "repositoryId or repository information is required",
		"RepositoryID:excluded_with":           "with repositoryId, no other repository information is needed",
		"Name:required_without_all":            "without repositoryId, name is required",
		"OperatingSystem:required_without_all": "without repositoryId, operationSystem is required",
		"TagID:required_without_all":           "tagId or tag information is required",
		"TagID:excluded_with":                  "with tagId, no other information is needed",
		"Version:required_without":             "without tagId, version is required",
		"Type:required_without":                "without tagId, type is required",

		// system
		"Filename:required": "filename is required",
		"From:required":     "from is required",
	}
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("timeformat", timeFormat, true)
		v.RegisterValidation("timeafter", timeAfter, true)
		v.RegisterValidation("timebefore", timeBefore, true)
	} else {
		fmt.Println("Validator Register FAILED.")
	}
}

func timeFormatExtractor(field reflect.StructField) string {
	tag := field.Tag.Get(_binding)
	bindingSlice := strings.Split(tag, ",")
	for _, s1 := range bindingSlice {
		if strings.Contains(s1, _timeFormat) {
			if timeFormatSlice := strings.Split(s1, "="); timeFormatSlice[0] == _timeFormat {
				if len(timeFormatSlice) <= 1 {
					break
				}
				format = timeFormatSlice[1]
				break
			}
		}
	}

	if v, ok := timeFormatsMap[format]; ok {
		return v
	} else {
		return ""
	}
}

func getTimeFormatAndField(fl validator.FieldLevel) (string, bool) {
	field, ok := reflect.TypeOf(fl.Parent().Interface()).FieldByName(fl.FieldName())
	if !ok {
		return "", false
	}

	if field.Type.Kind() == reflect.Ptr && reflect.ValueOf(fl.Parent().Interface()).FieldByName(fl.FieldName()).IsNil() {
		return "", true
	}

	format := timeFormatExtractor(field)
	if format == "" {
		return "", false
	}

	return format, true
}

func timeFormat(fl validator.FieldLevel) bool {
	format, valid := getTimeFormatAndField(fl)
	if !valid {
		return false
	}

	_, parseTimeErr := time.Parse(format, fl.Field().String())
	if parseTimeErr != nil {
		zap.L().With(
			zap.String(cnt.Controller, "time.Parse(...)"),
			zap.String("Param", fl.Field().String()),
		).Error(parseTimeErr.Error())
		return false
	}
	return true
}

func timeAfter(fl validator.FieldLevel) bool {
	format, valid := getTimeFormatAndField(fl)
	if !valid {
		return false
	}

	if fl.Param() != "" {
		if fl.Parent().FieldByName(fl.Param()).IsValid() {
			f1TimeFormat, _ := time.Parse(format, fl.Field().String())
			paramTimeFormat, _ := time.Parse(format, fl.Parent().FieldByName(fl.Param()).String())
			return paramTimeFormat.After(f1TimeFormat)
		} else {
			return false
		}
	}

	timeFormat, _ := time.Parse(format, fl.Field().String())
	return !time.Now().UTC().After(timeFormat)
}

func timeBefore(fl validator.FieldLevel) bool {
	format, valid := getTimeFormatAndField(fl)
	if !valid {
		return false
	}
	if fl.Param() != "" {
		if fl.Parent().FieldByName(fl.Param()).IsValid() {
			f1TimeFormat, _ := time.Parse(format, fl.Field().String())
			paramTimeFormat, _ := time.Parse(format, fl.Parent().FieldByName(fl.Param()).String())
			return paramTimeFormat.Before(f1TimeFormat)
		} else {
			return true
		}
	}
	timeFormat, _ := time.Parse(format, fl.Field().String())
	return time.Now().UTC().Before(timeFormat)
}

func Malformed(err error) error {
	if eMessage := shouldBindErrorParser(reflect.ValueOf(err)); eMessage != nil {
		var innerError error

		if val, ok := bindingErrorList[*eMessage]; ok {
			innerError = fmt.Errorf(val)
		} else {
			innerError = fmt.Errorf(*eMessage)
		}
		return tkErr.New(cnt.AdminAPISenderMalformedInputErr, innerError)
	}
	return tkErr.New(cnt.AdminAPISenderMalformedInputErr).WithInner(err)
}

func shouldBindErrorParser(v reflect.Value) *string {
	switch v.Kind() {
	case reflect.Ptr:
		if v.Type().Implements(reflect.TypeOf((*validator.FieldError)(nil)).Elem()) { // validation error
			field := v.MethodByName("Field").Call([]reflect.Value{})[0].String()
			tag := v.MethodByName("Tag").Call([]reflect.Value{})[0].String()
			key := fmt.Sprintf("%s:%s", field, tag)
			return &key
		}
		switch err := v.Interface().(type) {
		case *json.UnmarshalTypeError: // type error
			key := fmt.Sprintf("%s must be %s", err.Field, err.Type)
			return &key
		case *json.SyntaxError: // json format error
			key := "EOF"
			return &key
		case error:
			if errors.Is(err, io.EOF) { // empty message body
				key := "EOF"
				return &key
			}
		}
		return shouldBindErrorParser(v.Elem())
	case reflect.Slice:
		return shouldBindErrorParser(v.Index(0).Elem())
	default:
		return nil
	}
}
