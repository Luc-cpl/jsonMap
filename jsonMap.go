package jsonMap

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

func MakeJSON(jsonMap map[string]string, key string) (s string) {
	value := jsonMap[key]
	if arrayJSON := strings.TrimPrefix(value, "json:array "); arrayJSON != value {
		array := strings.Split(arrayJSON, " ")
		arraysN := ""
		for i := 0; i < len(array); i++ {
			newKey := array[i]
			if key != "" {
				newKey = key + " : " + array[i]
			}
			t := MakeJSON(jsonMap, newKey)
			if i > 0 {
				arraysN = arraysN + "," + t
			} else {
				arraysN = t
			}
		}
		s = "[" + arraysN + "]"
	} else if objectJSON := strings.TrimPrefix(value, "json:object "); objectJSON != value {
		array := strings.Split(objectJSON, " ")
		arraysN := ""
		for i := 0; i < len(array); i++ {
			newKey := array[i]
			if key != "" {
				newKey = key + " : " + array[i]
			}
			t := MakeJSON(jsonMap, newKey)
			if i > 0 {
				arraysN = arraysN + "," + `"` + array[i] + `":` + t
			} else {
				if array[i] != "" {
					arraysN = `"` + array[i] + `":` + t
				}
			}
		}
		s = "{" + arraysN + "}"
	} else {
		s = value
	}
	return s
}

func GetMap(bytJSON []byte, key string) (jsonMap map[string]string, err error) {
	get := make(map[string]string)
	var jsonArr []interface{}
	var jsonObj map[string]interface{}
	var jsonStr string
	var jsonBool bool
	var jsonInt int
	var jsonFloat float64
	err = json.Unmarshal(bytJSON, &jsonArr)
	if err != nil {
		err = json.Unmarshal(bytJSON, &jsonObj)
		if err != nil {
			err = json.Unmarshal(bytJSON, &jsonInt)
			if err != nil {
				err = json.Unmarshal(bytJSON, &jsonFloat)
				if err != nil {
					err = json.Unmarshal(bytJSON, &jsonBool)
					if err != nil {
						err = json.Unmarshal(bytJSON, &jsonStr)
						if err != nil {
							err = errors.New("unnable to get this json file into a map")
						} else {
							get[key] = `"` + jsonStr + `"`
						}
					} else {
						get[key] = strconv.FormatBool(jsonBool)
					}
				} else {
					get[key] = strconv.FormatFloat(jsonFloat, 'f', -1, 64)
				}
			} else {
				get[key] = strconv.Itoa(jsonInt)
			}
		} else {
			get[key] = "json:object "
			z := 0
			for n, el := range jsonObj {
				if z == 0 {
					get[key] += n
				} else {
					get[key] += " " + n
				}
				byt, _ := json.Marshal(el)
				rJSON := make(map[string]string)
				rJSON, err = GetMap(byt, newKey(key, n))
				for n, el := range rJSON {
					get[n] = el
				}
				z++
			}
		}
	} else {
		get[key] = "json:array "
		z := 0
		for n, el := range jsonArr {
			if z == 0 {
				get[key] += strconv.Itoa(n)
			} else {
				get[key] += " " + strconv.Itoa(n)
			}
			byt, _ := json.Marshal(el)
			rJSON := make(map[string]string)
			rJSON, err = GetMap(byt, newKey(key, strconv.Itoa(n)))
			for n, el := range rJSON {
				get[n] = el
			}
			z++
		}
	}

	return get, err
}

func newKey(old string, add string) string {
	new := ""
	if old == "" {
		new = add
	} else {
		new = old + " : " + add
	}
	return new
}
