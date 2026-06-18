package command

type ModuleName = string

type Metadata struct {
	Module ModuleName
	Alias  string
}

type FileType string

const (
	FileTypeTXT  FileType = "txt"
	FileTypeJSON FileType = "json"
	FileTypeHTML FileType = "html"
)
