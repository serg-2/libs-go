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

func UnmarshallString(source string, sourceInterface interface{}) interface{} {
	if err := json.Unmarshal([]byte(source), &sourceInterface); err != nil {
		log.Println("Can't unmarshall string: " + source)
	}
	return sourceInterface
}
