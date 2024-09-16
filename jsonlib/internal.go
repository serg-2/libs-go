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
