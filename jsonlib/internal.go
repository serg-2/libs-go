package jsonlib

import (
	"encoding/json"
	"log"
)

func marshallStructure(structureToMarshall interface{}) []byte {
	dataM, err := json.MarshalIndent(structureToMarshall, "", " ")
	if err != nil {
		log.Fatalln("Cant' marshall data.\nError:" + err.Error())
	}
	return dataM
}

func UnmarshallString[T any](source string) T {
	var loadedStructure T
	if err := json.Unmarshal([]byte(source), &loadedStructure); err != nil {
		log.Println("Can't unmarshall string: " + source)
	}
	return loadedStructure
}
