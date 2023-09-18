package fastqc

import "biolink-nipt-gui/internal/models"

type FastQC struct {
	Version      *Version
	InputFile    *InputFile
	OutputDir    *OutputDir
	Casava       *Casava
	Nano         *Nano
	NoFilter     *NoFilter
	Extract      *Extract
	Java         *Java
	NoExtract    *NoExtract
	NoGroup      *NoGroup
	MinLength    *MinLength
	Format       *Format
	Threads      *Threads
	Contaminants *Contaminants
	Adapters     *Adapters
	Limits       *Limits
	KMers        *KMers
	Quiet        *Quiet
	Dir          *Dir
	Icons        *models.Icons
}

func New() *FastQC {
	return &FastQC{
		Version: &Version{
			Flag:        "--version",
			Value:       "FastQC v0.11.9",
			Description: "Print the version of the program and exit",
			IsUsed:      true,
		},
		InputFile: &InputFile{
			IsUsed: true,
		},
		OutputDir: &OutputDir{
			Name:        "Output Directory",
			Flag:        "--outdir",
			Description: "Create all output files in the specified output directory.\nPlease note that this directory must exist as the program\nwill not create it.  If this option is not set then the\noutput file for each sequence file is created in the same\ndirectory as the sequence file which was processed.",
			IsUsed:      false,
		},
		Casava: &Casava{
			Name:        "Casava",
			Flag:        "--casava",
			Description: "Files come from raw casava output. Files in the same sample\ngroup (differing only by the group number) will be analysed\nas a set rather than individually. Sequences with the filter\nflag set in the header will be excluded from the analysis.\nFiles must have the same names given to them by casava\n(including being gzipped and ending with .gz) otherwise they\nwon't be grouped together correctly.",
			IsUsed:      false,
		},
		Nano: &Nano{
			Name:        "Nanopore",
			Flag:        "--nano",
			Description: "Files come from nanopore sequences and are in fast5 format. In\nthis mode you can pass in directories to process and the program\nwill take in all fast5 files within those directories and produce\na single output file from the sequences found in all files.",
			IsUsed:      false,
		},
		NoFilter: &NoFilter{
			Name:        "No Filter",
			Flag:        "--nofilter",
			Description: "If running with --casava then don't remove read flagged by casava as poor quality when performing the QC analysis.",
			IsUsed:      false,
		},
		Extract: &Extract{
			Name:        "Extract",
			Flag:        "--extract",
			Description: "If set then the zipped output file will be uncompressed in the same directory after it has been created.\nBy default this option will be set if fastqc is run in non-interactive mode.",
			IsUsed:      false,
		},
		Java: &Java{
			Name:        "Java",
			Flag:        "--java",
			Description: "Provides the full path to the java binary you want to use to launch fastqc. If not supplied then java is assumed to be in your path.",
			IsUsed:      false,
		},
		NoExtract: &NoExtract{
			Name:        "No Extract",
			Flag:        "--noextract",
			Description: "Do not uncompress the output file after creating it. You should set this option if you do not\nwish to uncompress the output when running in non-interactive mode.",
			IsUsed:      false,
		},
		NoGroup: &NoGroup{
			Name:        "No Group",
			Flag:        "--nogroup",
			Description: "Disable grouping of bases for reads >50bp. All reports will show data for every base in the read.\nWARNING: Using this option will cause fastqc to crash and burn if you use it on really long reads,\nand your plots may end up a ridiculous size. You have been warned!",
			IsUsed:      false,
		},
		MinLength: &MinLength{
			Name:        "Minimal Length",
			Flag:        "--min_length",
			Description: "Sets an artificial lower limit on the length of the sequence to be shown in the report.\nAs long as you set this to a value greater or equal to your longest read length then this will be\nthe sequence length used to create your read groups.  This can be useful for making directly comaparable statistics\nfrom datasets with somewhat variable read lengths.",
			IsUsed:      false,
		},
		Format: &Format{
			Name:        "Format",
			Flag:        "--format",
			Description: "Bypasses the normal sequence file format detection and forces the program to use the\nspecified format. Valid formats are bam,sam,bam_mapped,sam_mapped and fastq",
			Formats: []string{
				"bam",
				"sam",
				"bam_mapped",
				"sam_mapped",
				"fastq",
			},
			IsUsed: false,
		},
		Threads: &Threads{
			Name:        "Threads",
			Flag:        "--threads",
			Description: "Specifies the number of files which can be processed simultaneously.  Each thread will be allocated 250MB\nof memory so you shouldn't run more threads than your available memory will cope with,\nand not more than 6 threads on a 32 bit machine",
			IsUsed:      false,
		},
		Contaminants: &Contaminants{
			Name:        "Adapters",
			Flag:        "--contaminants",
			Description: "Specifies a non-default file which contains the list of contaminants to screen overrepresented sequences against.\nThe file must contain sets of named contaminants in the form name[tab]sequence.\nLines prefixed with a hash will be ignored.",
			IsUsed:      false,
		},
		Adapters: &Adapters{
			Name:        "Contaminants",
			Flag:        "--adapters",
			Description: "Specifies a non-default file which contains the list of adapter sequences which will be explicity searched\nagainst the library. The file must contain sets of named adapters in the form name[tab]sequence.\nLines prefixed with a hash will be ignored.",
			IsUsed:      false,
		},
		Limits: &Limits{
			Name:        "Limits",
			Flag:        "--limits",
			Description: "Specifies a non-default file which contains a set of criteria which will be used to determine the warn/error\nlimits for the various modules.  This file can also be used to selectively remove some modules from the output all together.\nThe format needs to mirror the default limits.txt file found in the Configuration folder.",
			IsUsed:      false,
		},
		KMers: &KMers{
			Name:        "K-Mers",
			Flag:        "--kmers",
			Description: "Specifies the length of Kmer to look for in the Kmer content module. Specified Kmer length must be between 2 and 10. Default length is 7 if not specified.",
			IsUsed:      false,
		},
		Quiet: &Quiet{
			Name:        "Quiet",
			Flag:        "--quiet",
			Description: "Supress all progress messages on stdout and only report errors.",
			IsUsed:      false,
		},
		Dir: &Dir{
			Name:        "Temporary Directory",
			Flag:        "--dir",
			Description: "Selects a directory to be used for temporary files written when generating report images. Defaults to system temp directory if not specified.",
			IsUsed:      false,
		},
	}
}

type Version struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}

type InputFile struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}

type OutputDir struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}

type Casava struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}

type Nano struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}

type NoFilter struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}

type Extract struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}

type Java struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}

type NoExtract struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}

type NoGroup struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}

type MinLength struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}

type Format struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
	Formats     []string
}

type Threads struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}

type Contaminants struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}

type Adapters struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}

type Limits struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}

type KMers struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}

type Quiet struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}

type Dir struct {
	Name        string
	Flag        string
	Value       string
	Description string
	IsUsed      bool
}
