package trimmomatic

import "fmt"

func (trimm *Trimmomatic) BuildMainCommand() (string, error) {
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
