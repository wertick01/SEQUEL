package models

import (
	"fyne.io/fyne/v2"
	validation "github.com/go-ozzo/ozzo-validation"
)

type TrimmomaticParams struct {
	Prefix      string
	Path        string
	Input       string
	Output      string
	Paired      string
	Logfile     string
	BaseOutput  bool
	Phred       int
	Threads     int
	SubParams   *TrimmomaticSubParams
	Description *TrimmomaticParamsDescription
}

func (params *TrimmomaticParams) Validate() error {
	return validation.ValidateStruct(
		params,
		validation.Field(&params.Prefix, validation.Required),
		validation.Field(&params.Path, validation.Required),
		validation.Field(&params.Input, validation.Required),
		validation.Field(&params.Output, validation.Required),
		// validation.Field(&params.Paired, validation.Required),
		validation.Field(&params.Paired, validation.Required),
		validation.Field(&params.Phred, validation.Required),
		validation.Field(&params.Threads, validation.Required),
		// validation.Field(params.SubParams, validation.Required),
	)
}

type TrimmomaticParamsDescription struct {
	Phred string
}

func NewMainParamsDescription() *TrimmomaticParamsDescription {
	return &TrimmomaticParamsDescription{
		Phred: "> Convert quality scores to Phred-33 or Phred-64.",
	}
}

type TrimmomaticSubParams struct {
	IlluminaClip  string
	SlidingWindow string
	Leading       string
	Trailing      string
	MinLen        string
	Crop          string
	HeadCrop      string
	Description   *TrimmomaticSubParamsDescription
	Names         *TrimmomaticSubParamsNames
	QuestionWidth float32
}

type Icons struct {
	Submit   fyne.Resource
	Warning  fyne.Resource
	Question fyne.Resource
}

type TrimmomaticSubParamsDescription struct {
	IlluminaClip  string
	SlidingWindow string
	Leading       string
	Trailing      string
	MinLen        string
	Crop          string
	HeadCrop      string
}

func NewDescription() *TrimmomaticSubParamsDescription {
	return &TrimmomaticSubParamsDescription{
		IlluminaClip:  "> ILLUMINACLIP:<fastaWithAdaptersEtc>:<seed mismatches>:<palindrome clip threshold>:<simple clip threshold>\n> Cut adapter and other illumina-specific sequences from the read.\n> fastaWithAdaptersEtc: specifies the path to a fasta file containing all the adapters, PCR sequences etc. The naming of the various sequences within this file determines how they are used. See below.\n> seedMismatches: specifies the maximum mismatch count which will still allow a full match to be performed\n> palindromeClipThreshold: specifies how accurate the match between the two 'adapter ligated' reads must be for PE palindrome read alignment.\n> simpleClipThreshold: specifies how accurate the match between any adapter etc. sequence must be against a read.",
		SlidingWindow: "> SLIDINGWINDOW:<windowSize>:<requiredQuality>\n> Perform a sliding window trimming, cutting once the average quality within the window falls below a threshold.\n> windowSize: specifies the number of bases to average across\n> requiredQuality: specifies the average quality required.",
		Leading:       "> LEADING:<quality>\n> Cut bases off the start of a read, if below a threshold quality\n> quality: Specifies the minimum quality required to keep a base",
		Trailing:      "> TRAILING:<quality>\n> Cut bases off the end of a read, if below a threshold quality\n> quality: Specifies the minimum quality required to keep a base.",
		MinLen:        "> MINLEN:<length>\n> Cut the read to a specified length\n> length: Specifies the minimum length of reads to be kept.",
		Crop:          "> CROP:<length>\n> Cut the read to a specified length\n> length: The number of bases to remove from the start of the read.",
		HeadCrop:      "> HEADCROP:<length>\n> Cut the specified number of bases from the start of the read\n> length: The number of bases to remove from the start of the read.",
	}
}

type TrimmomaticSubParamsNames struct {
	IlluminaClip  map[int]string
	SlidingWindow map[int]string
	Leading       string
	Trailing      string
	MinLen        string
	Crop          string
	HeadCrop      string
}

func NewSubparamsNames() *TrimmomaticSubParamsNames {
	return &TrimmomaticSubParamsNames{
		IlluminaClip: map[int]string{
			0: "fastaWithAdaptersEtc",
			1: "seed mismatches",
			2: "palindrome clip threshold",
			3: "simple clip threshold",
		},
		SlidingWindow: map[int]string{
			0: "windowSize",
			1: "requiredQuality",
		},
		Leading:  "quality",
		Trailing: "quality",
		MinLen:   "length",
		Crop:     "length",
		HeadCrop: "length",
	}
}

// func (params *TrimmomaticSubParams) Validate() error {
// 	return validation.ValidateStruct(
// 		params,
// 		validation.Field(&params.IlluminaClip, validation.Required),
// 		validation.Field(&params.Leading, validation.Required),
// 		validation.Field(&params.Trailing, validation.Required),
// 		validation.Field(&params.MinLen, validation.Required),
// 		validation.Field(&params.Crop, validation.Required),
// 		validation.Field(&params.HedaCrop, validation.Required),
// 		validation.Field(&params.SlidingWindow, validation.Required),
// 	)
// }
