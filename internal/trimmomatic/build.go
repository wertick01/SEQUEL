package trimmomatic

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

func (trimm *Trimmomatic) BuildMainCommand() (string, error) {
	if err := trimm.Params.Validate(); err != nil {
		return "", err
	}
	log.Println(trimm.Params)

	var outputPrefix string
	var outputFiles string
	var separator string

	if runtime.GOOS == "windows" {
		separator = `\`
	} else {
		separator = "/"
	}

	if trimm.Params.Paired == "PE" {
		mass := strings.Split(trimm.Params.Input, " ")
		if trimm.Params.BaseOutput {
			mass2 := strings.Split(mass[0], separator)
			inputFile := mass2[len(mass2)-1]
			mass3 := strings.Split(inputFile, ".")
			outputPrefix = mass3[0]
			format := ""
			for _, val := range mass3[1:] {
				format += "." + val
			}
			outputPrefix = outputPrefix + "_output" + format
			outputFiles = fmt.Sprintf("-baseout %v%v%v", trimm.Params.Output, separator, outputPrefix)
		} else {
			mass2, mass3 := strings.Split(mass[0], separator), strings.Split(mass[1], separator)
			fow, rev := mass2[len(mass2)-1], mass3[len(mass3)-1]
			fowFileNameSlice, revFileNameSlice := strings.Split(fow, "."), strings.Split(rev, ".")
			fowFileName, revFileName := fowFileNameSlice[0], revFileNameSlice[0]
			formatFow, formatRev := "", ""
			for _, val := range fowFileNameSlice[1:] {
				formatFow += "." + val
			}
			for _, val := range revFileNameSlice[1:] {
				formatRev += "." + val
			}
			fowPaired, fowUnpaired := trimm.Params.Output+separator+fowFileName+"_paired_output"+formatFow, trimm.Params.Output+separator+fowFileName+"_unpaired_output"+formatFow
			revPaired, revUnpaired := trimm.Params.Output+separator+revFileName+"_paired_output"+formatRev, trimm.Params.Output+separator+revFileName+"_unpaired_output"+formatRev
			outputFiles = fmt.Sprintf("%v %v %v %v", fowPaired, fowUnpaired, revPaired, revUnpaired)
		}
	} else {
		mass := strings.Split(trimm.Params.Input, ".")
		inputFileName := mass[0]
		outputFormat := ""
		for _, val := range mass[1:] {
			outputFormat += "." + val
		}
		outputFiles = trimm.Params.Output + separator + inputFileName + "_output" + outputFormat
	}

	res := fmt.Sprintf(
		"%v %v %v -threads %v -phred%v -trimlog %v %v %v",
		trimm.Params.Prefix,
		trimm.Params.Path,
		trimm.Params.Paired,
		trimm.Params.Threads,
		trimm.Params.Phred,
		trimm.Params.Logfile,
		trimm.Params.Input,
		outputFiles,
	)

	log.Println(res)

	return res, nil
}

func (trimm *Trimmomatic) BuildInputFiles() (string, error) {

	return "", nil
}
