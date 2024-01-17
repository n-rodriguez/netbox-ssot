package utils

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
)

// Helper function to determine if a given reflect.Value contains an embedded objects.Choice
func isChoiceEmbedded(v reflect.Value) bool {
	vType := v.Type()
	return vType.Field(0).Type == reflect.TypeOf(objects.Choice{})
	// for i := 0; i < v.NumField(); i++ {
	// 	if vType.Field(i).Type == reflect.TypeOf(objects.Choice{}) {
	// 		return true
	// 	}
	// }
	// return false
}

func choiceValue(v reflect.Value) interface{} {
	return v.Field(0).FieldByName("Value").Interface()
}

// Struct used for patching objects, when attributes are structs, or slices
// we only need object with an id of original object for patching.
type IDObject struct {
	ID int `json:"id"`
}

// JsonDiffMapExceptID checks if any field except ID field, from two objects, are different.
// It returns a map of fields (represented by their JSON tag names) that are different with their values from newObj.
func JsonDiffMapExceptId(newObj, existingObj interface{}) (map[string]interface{}, error) {
	diff := make(map[string]interface{})

	newObject := reflect.ValueOf(newObj)
	existingObject := reflect.ValueOf(existingObj)

	// Ensure that both objects are of the same kind (e.g. struct)
	if newObject.Kind() != existingObject.Kind() {
		return nil, fmt.Errorf("arguments are not of the same type")
	}

	// Check if the values are pointers and get the element they point to
	if newObject.Kind() == reflect.Ptr {
		newObject = newObject.Elem()
		existingObject = existingObject.Elem()
	}

	// Ensure that we are dealing with structs
	if newObject.Kind() != reflect.Struct {
		return nil, fmt.Errorf("arguments are not structs")
	}

	for i := 0; i < newObject.NumField(); i++ {
		fieldName := newObject.Type().Field(i).Name
		jsonTag := newObject.Type().Field(i).Tag.Get("json")

		if fieldName == "Id" {
			continue
		}

		// Custom logic for all objects that inherit from NetboxObject
		if fieldName == "NetboxObject" {
			netboxObjectDiffMap, err := JsonDiffMapExceptId(newObject.Field(i).Interface(), existingObject.Field(i).Interface())
			if err != nil {
				return nil, fmt.Errorf("error processing JsonDiffMapExceptId when processing NetboxObject %s", err)
			}
			for k, v := range netboxObjectDiffMap {
				diff[k] = v
			}
			continue
		}

		// Get json tag, so we can get the json field name
		if jsonTag == "" || jsonTag == "-" {
			jsonTag = fieldName
		} else {
			jsonTag = strings.Split(jsonTag, ",")[0]
		}

		// Ensure that both fields are of the same kind
		if newObject.Field(i).Kind() != existingObject.Field(i).Kind() {
			return nil, fmt.Errorf("field %s is not of the same type in both objects", jsonTag)
		}

		// Check if elements are pointers, in that case get the elements they are pointing to
		newObjectField := newObject.Field(i)
		existingObjectField := existingObject.Field(i)
		if newObjectField.Kind() == reflect.Ptr {
			newObjectField = newObjectField.Elem()
			existingObjectField = existingObjectField.Elem()
		}

		switch newObjectField.Kind() {

		case reflect.Invalid:
			// the newObjectField is nil, in this case we don't do anything (like omitempty), see tests
			continue

		case reflect.Slice:
			err := addSliceDiff(newObjectField, existingObjectField, jsonTag, diff)
			if err != nil {
				return nil, fmt.Errorf("error processing JsonDiffMapExceptId when processing slice %s", err)
			}

		case reflect.Struct:
			err := addStructDiff(newObjectField, existingObjectField, jsonTag, diff)
			if err != nil {
				return nil, fmt.Errorf("error processing JsonDiffMapExceptId when processing struct %s", err)
			}

		case reflect.Map:
			err := addMapDiff(newObjectField, existingObjectField, jsonTag, diff)
			if err != nil {
				return nil, fmt.Errorf("error processing JsonDiffMapExceptId when processing map %s", err)
			}

		default:
			if !newObjectField.IsValid() {
				if existingObjectField.IsValid() {
					diff[jsonTag] = nil
				}
				continue
			}
			if newObjectField.Interface() != existingObjectField.Interface() {
				diff[jsonTag] = newObjectField.Interface()
			}
		}
	}

	return diff, nil
}

// Function that takes two objects (of type slice) and returns a map
// that can be easily used with json.Marshal
// To achieve this the map is of the following format:
// map[jsonTag] = [id1, id2, id3] // If the slice contains objects with ID field
// map[jsonTag] = [value1, value2, value3] // If the slice contains strings
func addSliceDiff(newSlice reflect.Value, existingSlice reflect.Value, jsonTag string, diffMap map[string]interface{}) error {

	// If first slice is nil, that means that we reset the value
	if !newSlice.IsValid() || newSlice.Len() == 0 {
		if existingSlice.IsValid() && existingSlice.Len() > 0 {
			diffMap[jsonTag] = []interface{}{}
		}
		return nil
	}

	// There are going to be 2 kinds of comparison.
	// One where slice will contain objects, in that case we
	// will compare ids of the objects, else we will compare
	// the values of the slice
	switch newSlice.Index(0).Kind() {
	case reflect.String:
		strSet := make(map[string]bool)
		for j := 0; j < newSlice.Len(); j++ {
			strSet[newSlice.Index(j).Interface().(string)] = true
		}
		if len(strSet) != existingSlice.Len() {
			diffMap[jsonTag] = newSlice.Interface()
		} else {
			for j := 0; j < existingSlice.Len(); j++ {
				if !strSet[existingSlice.Index(j).Interface().(string)] {
					diffMap[jsonTag] = newSlice.Interface()
					return nil
				}
			}
		}

	default:
		newIdSet := make(map[int]bool, newSlice.Len())
		var id int
		for j := 0; j < newSlice.Len(); j++ {
			element := newSlice.Index(j)
			if newSlice.Index(j).IsNil() {
				continue
			}
			if element.Kind() == reflect.Ptr {
				element = element.Elem()
			}
			id := element.FieldByName("Id").Interface().(int)
			newIdSet[id] = true
		}

		newIdSlice := make([]int, 0, len(newIdSet))
		for id := range newIdSet {
			newIdSlice = append(newIdSlice, id)
		}
		slices.Sort(newIdSlice)

		if len(newIdSet) != existingSlice.Len() {
			diffMap[jsonTag] = newIdSlice
		} else {
			for j := 0; j < existingSlice.Len(); j++ {
				existingSliceEl := existingSlice.Index(j)
				if existingSlice.Index(j).Kind() == reflect.Ptr {
					existingSliceEl = existingSliceEl.Elem()
				}
				id = existingSliceEl.FieldByName("Id").Interface().(int)
				if _, ok := newIdSet[id]; !ok {
					diffMap[jsonTag] = newIdSlice
					return nil
				}
			}
		}
	}
	return nil
}

// Returns json form for patching the difference e.g. { "Id": 1 }
func addStructDiff(newObj reflect.Value, existingObj reflect.Value, jsonTag string, diffMap map[string]interface{}) error {

	// If first struct is nil, that means that we reset the attribute to nil
	if !newObj.IsValid() {
		diffMap[jsonTag] = nil
		return nil
	}

	// We check if struct is a objects.Choice (special netbox struct)
	if isChoiceEmbedded(newObj) {
		if !existingObj.IsValid() || newObj.Interface() != existingObj.Interface() {
			diffMap[jsonTag] = choiceValue(newObj)
		}
		return nil
	}

	// We use ids for comparison between structs, because for patching objects, all we need is id of attribute
	idField := newObj.FieldByName("Id")

	// If objects don't have ID field, compare them by their values
	if !idField.IsValid() {
		if !existingObj.IsValid() {
			diffMap[jsonTag] = newObj.Interface()
		} else {
			if newObj.Interface() != existingObj.Interface() {
				if isChoiceEmbedded(newObj) {
					diffMap[jsonTag] = choiceValue(newObj)
				} else {
					diffMap[jsonTag] = newObj.Interface()
				}
			}
		}
	} else {
		if !existingObj.IsValid() {
			diffMap[jsonTag] = IDObject{ID: idField.Interface().(int)}
		} else {
			// Objects have ID field, compare their ids
			if newObj.FieldByName("Id").Interface() != existingObj.FieldByName("Id").Interface() {
				id := newObj.FieldByName("Id").Interface().(int)
				diffMap[jsonTag] = IDObject{ID: id}
			}
		}
	}
	return nil
}

func addMapDiff(newMap reflect.Value, existingMap reflect.Value, jsonTag string, diffMap map[string]interface{}) error {

	// If the new map is not set, we don't change anything
	if !newMap.IsValid() {
		return nil
	}

	// Go through all keys in new map, and check if they are in existing map
	// If they are not, add them to diff map
	mapsDiff := make(map[string]interface{})
	for _, key := range newMap.MapKeys() {
		// Keys have to be strings
		if key.Kind() != reflect.String {
			return fmt.Errorf("map keys have to be strings. Not implemented for anything else yet")
		}
		if !existingMap.MapIndex(key).IsValid() {
			mapsDiff[key.Interface().(string)] = newMap.MapIndex(key).Interface()
		} else if newMap.MapIndex(key).Interface() != existingMap.MapIndex(key).Interface() {
			mapsDiff[key.Interface().(string)] = newMap.MapIndex(key).Interface()
		}
	}

	if len(mapsDiff) > 0 {
		for _, key := range existingMap.MapKeys() {
			if !newMap.MapIndex(key).IsValid() {
				mapsDiff[key.Interface().(string)] = existingMap.MapIndex(key).Interface()
			}
		}
		diffMap[jsonTag] = mapsDiff
	}
	return nil
}
