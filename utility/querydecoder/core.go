package querydecoder

import (
	"VirtualRegistryManagement/utility"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// where: <key><symbol><value>
//
//	key: 只允許英文字母，數字以及 -
//	symbol: 只允許 >=, <=, >, <, =, ==, <>
//	value: 不限制
const (
	whereRegexStr = `^(?P<key>[a-zA-Z0-9-]*)(?P<symbol>(>=|<=|==|<>|!=)|(=|<|>))(?P<value>.*)$`
)

var whereRegex = utility.MustCompile(whereRegexStr)

type (
	query struct {
		Key    string `regroup:"key,required"`
		Symbol string `regroup:"symbol,required"`
		Value  string `regroup:"value,required"`
	}

	RegexError struct {
		err error // key from the source map.
	}
)

func (e RegexError) Error() string {
	return fmt.Sprintf("schema: invalid regex %q", e.err.Error())
}

// ShouldBindWhereSlice 檢查 slice 中的條件是否合法
//
//	除 label-type 外，需定義於 obj 中
func ShouldBindWhereSlice(obj QueryInterface, slice []string) error {
	urlV := url.Values{}
	errors := MultiError{}
	for _, f := range slice {
		// store label if string is label-type
		if strings.HasPrefix(f, labelPrefix) {
			s := strings.SplitN(f, "=", 2)
			obj.AddLabel(s[len(s)-1])
			continue
		}
		q := &query{}
		if err := whereRegex.MatchToTarget(f, q); err != nil {
			errors[q.Key] = RegexError{err: err}
			continue
		}
		key := strings.ToLower(q.Key)
		urlV.Add(key, q.Value)
		obj.AddWhere(key, q.Symbol, q.Value)
	}
	if len(errors) > 0 {
		return errors
	}
	return ShouldBindWhere(obj, urlV)
}

func ShouldBindWhere(obj QueryInterface, src url.Values) error {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		for i := 0; i < t.NumField(); i++ {
			// if tag is case:lower
			// trans the value to lower case
			if t.Field(i).Tag.Get("case") == "lower" {
				if v := src.Get(strings.ToLower(t.Field(i).Name)); v != "" {
					src.Set(strings.ToLower(t.Field(i).Name), strings.ToLower(v))
				}
			}
		}
	}
	decoder := NewDecoder()
	decoder.SetAliasTag("where")
	return decoder.Decode(obj, src)
}
