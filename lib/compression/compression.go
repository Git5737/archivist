package compression

type Encoder interface {
	Encode(sourcePaths []string) error
}

type Decoder interface {
	Decode(outputDir string) error
}
