package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type TrimmomaticParams struct {
	Prefix     string
	Path       string
	Input      string
	Output     string
	Paired     string
	Logfile    string
	BaseOutput bool
	Phred      int
	Threads    int
	SubParams  *TrimmomaticSubParams
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

type TrimmomaticSubParams struct {
	IlluminaClip  string
	SlidingWindow string
	Leading       string
	Trailing      string
	MinLen        string
	Crop          string
	HeadCrop      string
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
