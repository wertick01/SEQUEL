package fastqc

import (
	"fmt"
	"log"
	"reflect"
)

func (fastQC *FastQC) BuildCommand() string {
	command := ""
	v := reflect.ValueOf(*fastQC)
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsNil() {
			name := reflect.TypeOf(*fastQC).Field(i).Name
			if !contains([]string{"Version", "Icons", "InputFile", "OutputDir"}, name) {
				if v.Field(i).Elem().FieldByName("IsUsed").Interface().(bool) {
					value := v.Field(i).Elem().FieldByName("Value").Interface().(string)
					if len(value) > 0 {
						command = fmt.Sprintf("%v %v %v", command, v.Field(i).Elem().FieldByName("Flag").Interface(), value)
					} else {
						command = fmt.Sprintf("%v %v", command, v.Field(i).Elem().FieldByName("Flag").Interface())
					}
				}
			}
		}
	}
	command = fmt.Sprintf("fastqc %v %v %v %v", fastQC.InputFile.Value, fastQC.OutputDir.Flag, fastQC.OutputDir.Value, command)
	log.Println(command)
	return ""
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
