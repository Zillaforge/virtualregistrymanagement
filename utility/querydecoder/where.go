package querydecoder

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const labelPrefix = "*label*"

var _ QueryInterface = (*Query)(nil)

type (
	QueryInterface interface {
		GetWhere(key string) []condition
		AddWhere(key, symbol string, value interface{})
		GetLabels() []string
		AddLabel(label string)
	}

	// Query 用來紀錄搜尋條件以及標籤
	//   使用時宣告 struct 並繼承 where
	//   Field Tag:
	//     where: 欲支援搜尋之欄位
	//     label: 設定欲搜尋標籤之欄位
	//     prefix: 當使用 join 語法時，需額外設定 table-name
	//
	//   type exampleWhere struct {
	//       id    *string `where:"id"`
	//       Query `label:"label"`      // 需要繼承，非定義該型態之變數
	//   }
	Query struct {
		// where: key: condition{symbol, value}
		where map[string][]condition

		// label: {"condition1", "condition2"}
		label []string
	}

	condition struct {
		Operator string
		Value    interface{}
	}
)

// Check 用於檢查是否有將參數存至 where 中
func Check(w QueryInterface) {
	nonPtrReflect := func(input interface{}) (reflect.Value, reflect.Type) {
		v, t := reflect.ValueOf(input), reflect.TypeOf(input)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		return v, t
	}

	v, t := nonPtrReflect(w)
	for i := 0; i < t.NumField(); i++ {
		if v.Field(i).IsZero() {
			continue
		}
		switch v.Field(i).Interface().(type) {
		case Query:
			continue
		}
		if key := strings.ToLower(t.Field(i).Tag.Get("where")); w.GetWhere(key) == nil {
			fV, fT := nonPtrReflect(v.Field(i).Interface())
			switch fT.Kind() {
			case reflect.Slice:
				for i := 0; i < fV.Len(); i++ {
					w.AddWhere(key, "=", fV.Index(i))
				}
			default:
				w.AddWhere(key, "=", fV.Interface())
			}
		}
	}
}

func (w *Query) GetWhere(key string) []condition {
	if v, exist := w.where[key]; exist {
		return v
	}
	return nil
}

func (w *Query) AddWhere(key, symbol string, value interface{}) {
	if w.where == nil {
		w.where = make(map[string][]condition)
	}

	switch v := value.(type) {
	case string:
		if b, err := strconv.ParseBool(v); err == nil {
			value = b
		}
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			value = f
		}
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			value = i
		}
		if u, err := strconv.ParseUint(v, 10, 64); err == nil {
			value = u
		}
	}

	w.where[key] = append(w.where[key], condition{Operator: symbol, Value: value})
}

func (w *Query) GetLabels() []string {
	return w.label
}

func (w *Query) AddLabel(label string) {
	w.label = append(w.label, label)
}

func WhereAppendLabels(where, labels []string) []string {
	for _, label := range labels {
		where = append(where, fmt.Sprintf("%s=%s", labelPrefix, label))
	}
	return where
}
