package wx

import "github.com/hhcool/gtls/structs"

func StructToMap(data interface{}) map[string]interface{} {
	return structs.Map(data)
}
