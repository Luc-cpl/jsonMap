package jsonMap

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

// MakeJSON from an map[string]string
func MakeJSON(jsonMap map[string]string, key string) (s string) {
	value := jsonMap[key]
	if arrayJSON := strings.TrimPrefix(value, "json:array "); arrayJSON != value {
		array := strings.Split(arrayJSON, " ")
		arraysN := ""
		for i := 0; i < len(array); i++ {
			newKey := array[i]
			if key != "" {
				newKey = key + " " + array[i]
			}
			t := MakeJSON(jsonMap, newKey)
			if arraysN != "" {
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
				newKey = key + " " + array[i]
			}
			t := MakeJSON(jsonMap, newKey)
			if arraysN != "" {
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

//GetMap gets a json file and transform into a map[string]string
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
		new = old + " " + add
	}
	return new
}

//UpdateValue updates the values from each key on the original jsonMap but doesn't update arrays or objects
func UpdateValue(jsonMap map[string]string, updateMap map[string]string) (newMap map[string]string, err error) {
	s := make(map[string]string)
	s = jsonMap
	for key := range updateMap {
		t1 := strings.TrimPrefix(jsonMap[key], "json:array")
		t2 := strings.TrimPrefix(jsonMap[key], "json:object")
		if jsonMap[key] != "" && t1 != jsonMap[key] && t2 != jsonMap[key] {
			s[key] = updateMap[key]
		} else if jsonMap[key] != "" && (t1 != jsonMap[key] || t2 != jsonMap[key]) {
			err = errors.New("can't update array or object, the map return untouched")
			return jsonMap, err
		} else {
			err = errors.New("some key value thoes not exist, the map return untouched")
			return jsonMap, err
		}
	}
	return s, nil
}

//Create a new object, array or value from each key passed hoes have an firt point on the original jsonMap and does not exists
//to create from an array, just put the key value as a letter or word (ensure that the key does not exists on original jsonMap)
//example of an createMap:
//	["0 user"] = "contacts"								 - add a contacts field on the object
//	["0 user contacts"] = "json:array 0"				 - place the contacts field like an array and put an key inside
//	["0 user contacts 0"] = "json:object contactName" 	 - take the field from the contacts array like an object and put a contactName field inside
//	["0 user contacts 0 contactName"] = "Luke Skywalker" - place a value inside contactName field
func Create(jsonMap map[string]string, createMap map[string]string) (newMap map[string]string, err error) {
	s := make(map[string]string)
	s = jsonMap
	for key := range createMap {
		t1 := strings.TrimPrefix(s[key], "json:array")
		t2 := strings.TrimPrefix(s[key], "json:object")
		if jsonMap[key] != "" && (t1 != jsonMap[key] || t2 != jsonMap[key]) {
			s[key] = s[key] + " " + createMap[key]
		} else if jsonMap[key] == "" {
			s[key] = createMap[key]
		} else {
			err = errors.New("some key value alread exists, the map return untouched")
			return jsonMap, err
		}
	}
	return s, nil
}

//Delete a object or array from each key hoes eist on the original jsonMap, if the jsonMap becomes empyt, the key is removed from the json result
//exemple of an deleteMap key and value: ["0 user"] = "deleteValue1 deleteValue2"
func Delete(jsonMap map[string]string, deleteMap map[string]string) map[string]string {
	s := make(map[string]string)
	s = jsonMap
	for key := range deleteMap {
		if jsonMap[key] != "" {
			deleteArray := strings.Split(deleteMap[key], " ")
			for i := 0; i < len(deleteArray); i++ {
				s[key] = strings.Replace(s[key], " "+deleteArray[i], "", 1)
				if s[key] == "json:array" || s[key] == "json:object" {
					s[key] = ""
				}
			}
		}
	}
	return s
}
