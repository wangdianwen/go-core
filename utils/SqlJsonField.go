package utils

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
)

type SqlJsonField map[string]interface{}

func (e *SqlJsonField) String() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func (e *SqlJsonField) FieldType() int {
	return orm.TypeTextField
}

func (e *SqlJsonField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case string:
		return json.Unmarshal([]byte(d), e)
	default:
		return fmt.Errorf("<JSONField.SetRaw> unknown value `%v`", value)
	}
}

func (e *SqlJsonField) RawValue() interface{} {
	return e.String()
}

var _ orm.Fielder = new(SqlJsonField)
