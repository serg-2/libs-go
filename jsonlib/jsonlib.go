package jsonlib

import (
	"encoding/json"
	"log"
	"os"
	"reflect"

	"github.com/go-playground/validator/v10"
)

func JsonAsString(doc interface{}) string {
	return string(marshallStructure(doc))
}

func LoadJsonFromFile[T any](file string) T {
	var loadedStructure T
	fileToRead, err := os.Open(file)
	if err != nil {
		log.Fatalln("Can't open file to load" + file + "\nError:" + err.Error())
	}
	defer fileToRead.Close()

	jsonParser := json.NewDecoder(fileToRead)
	err = jsonParser.Decode(&loadedStructure)
	if err != nil {
		log.Fatalln("Bad JSON in " + file + "\nError:" + err.Error())
	}

	// Validate only structs!
	if !(reflect.TypeOf(loadedStructure).Kind() == reflect.Slice || reflect.TypeOf(loadedStructure).Kind() == reflect.Map) {
		err = validator.New().Struct(loadedStructure)
		if err != nil {
			log.Fatalln("Error of validation: " + err.Error())
		}
	}

	return loadedStructure
}

func SaveJsonToFile(fileName string, structureToWrite interface{}) {
	dataM := marshallStructure(structureToWrite)

	err := os.WriteFile(fileName, dataM, 0644)
	if err != nil {
		log.Fatalln("Can't write file: " + fileName + "\nError:" + err.Error())
	}
}
