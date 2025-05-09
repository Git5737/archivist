package cmd

import (
	"archivist/lib/compression"
	"archivist/lib/compression/vlc"
	"errors"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var packcmd = &cobra.Command{
	Use:   "pack",
	Short: "Pack file",
	Run:   pack,
}

const packedExtension = "vlc"

var ErrEmptyPath = errors.New("path to file is not specified")

func pack(cmd *cobra.Command, args []string) {
	var encode compression.Encoder

	if len(args) == 0 || args[0] == "" {
		handleErr(ErrEmptyPath)
	}

	method := cmd.Flag("method").Value.String()

	switch method {
	case "vlc":
		encode = vlc.New()
	default:
		cmd.PrintErr("unknown method")
	}

	filePath := args[0]

	r, err := os.Open(filePath)
	if err != nil {
		handleErr(err)
	}
	defer r.Close()

	//var w io.Writer

	data, err := io.ReadAll(r)
	if err != nil {
		handleErr(err)
	}

	packed := encode.Encode(string(data))

	err = os.WriteFile(packedFileName(filePath), packed, 0644)
	if err != nil {
		handleErr(err)
	}
}

func packedFileName(path string) string {
	fileName := filepath.Base(path)

	return strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "." + packedExtension
}

func init() {
	rootCmd.AddCommand(packcmd)

	packcmd.Flags().StringP("method", "m", "", "compression method: vlc")
	if err := packcmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}
}
