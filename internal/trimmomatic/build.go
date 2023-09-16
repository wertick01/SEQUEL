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

	if strings.Contains(trimm.Params.Output, " ") {
		trimm.Params.Output = DeleteSpaces(trimm.Params.Output, separator)
	}
	if strings.Contains(trimm.Params.Logfile, " ") {
		trimm.Params.Logfile = DeleteSpaces(trimm.Params.Logfile, separator)
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
		mass := strings.Split(trimm.Params.Input, separator)
		mass1 := strings.Split(mass[len(mass)-1], ".")
		inputFileName := mass1[0]
		outputFormat := ""
		for _, val := range mass1[1:] {
			outputFormat += "." + val
		}
		outputFiles = trimm.Params.Output + separator + inputFileName + "_output" + outputFormat
	}

	subParams, err := trimm.BuildSubParams()
	if err != nil {
		log.Println(err)
	}

	res := fmt.Sprintf(
		"%v %v %v -threads %v -phred%v -trimlog %v %v %v %v",
		trimm.Params.Prefix,
		trimm.Params.Path,
		trimm.Params.Paired,
		trimm.Params.Threads,
		trimm.Params.Phred,
		trimm.Params.Logfile,
		trimm.Params.Input,
		outputFiles,
		subParams,
	)

	log.Println(res)

	return res, nil
}

func (trimm *Trimmomatic) BuildSubParams() (string, error) {
	subParams := ""

	if len(trimm.Params.SubParams.IlluminaClip) > len("ILLUMINACLIP::::") {
		subParams = fmt.Sprintf("%v %v", subParams, trimm.Params.SubParams.IlluminaClip)
	}
	if len(trimm.Params.SubParams.SlidingWindow) > len("SLIDINGWINDOW::") {
		subParams = fmt.Sprintf("%v %v", subParams, trimm.Params.SubParams.SlidingWindow)
	}
	if len(trimm.Params.SubParams.Leading) > len("LEADING:") {
		subParams = fmt.Sprintf("%v %v", subParams, trimm.Params.SubParams.Leading)
	}
	if len(trimm.Params.SubParams.Trailing) > len("TRAILING:") {
		subParams = fmt.Sprintf("%v %v", subParams, trimm.Params.SubParams.Trailing)
	}
	if len(trimm.Params.SubParams.Crop) > len("CROP:") {
		subParams = fmt.Sprintf("%v %v", subParams, trimm.Params.SubParams.Crop)
	}
	if len(trimm.Params.SubParams.HeadCrop) > len("HEADCROP:") {
		subParams = fmt.Sprintf("%v %v", subParams, trimm.Params.SubParams.HeadCrop)
	}
	if len(trimm.Params.SubParams.MinLen) > len("MINLEN:") {
		subParams = fmt.Sprintf("%v %v", subParams, trimm.Params.SubParams.MinLen)
	}

	return subParams, nil
}

func DeleteSpaces(path, separator string) string {
	splited := strings.Split(path, separator)
	outputDir := ""
	for _, val := range splited {
		if strings.Contains(val, " ") {
			val = fmt.Sprintf("'%v'", val)
		}
		// outputDir = fmt.Sprintf("%v%v%v", outputDir, separator, val)
		outputDir += separator + val
	}
	path = outputDir[1:]
	return path
}
