package trimmomatic

import (
	"fmt"
)

func (trimm *Trimmomatic) BuildMainCommand() (string, error) {
	trimm.Params.Threads = 20
	trimm.Params.Phred = 33
	trimm.Params.Paired = "PE"
	trimm.Params.Output = "Output"
	if err := trimm.Params.Validate(); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%v %v %v -phred%v -threads %v %v %v",
		trimm.Params.Prefix,
		trimm.Params.Path,
		trimm.Params.Paired,
		trimm.Params.Threads,
		trimm.Params.Phred,
		trimm.Params.Input,
		trimm.Params.Output,
	), nil
}

func (trimm *Trimmomatic) BuildInputFiles() (string, error) {

	return "", nil
}
