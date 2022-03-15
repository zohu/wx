package wx

import "github.com/hhcool/structs"

func StructToMap(data interface{}) map[string]interface{} {
	return structs.Map(data)
}
