package cmd

import (
	"archivist/lib/compression"
	tar2 "archivist/lib/compression/tar"
	"archivist/lib/compression/tar_bz2"
	"archivist/lib/compression/tar_gz"
	"archivist/lib/compression/tar_xz"
	"archivist/lib/compression/zip"
	"errors"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

var packcmd = &cobra.Command{
	Use:   "pack",
	Short: "Pack file",
	Run:   pack,
}

var ErrEmptyPath = errors.New("path to file is not specified")

func pack(cmd *cobra.Command, args []string) {
	filePath := args[0]

	if len(args) == 0 || args[0] == "" {
		handleErr(ErrEmptyPath)
	}

	var encode compression.Encoder

	method := cmd.Flag("method").Value.String()

	packedName := packedFileName(filePath, method)

	switch method {
	case "zip":
		encode = zip.New(packedName)
	case "tar":
		encode = tar2.New(packedName)
	case "tar.gz":
		encode = tar_gz.New(packedName)
	case "tar.xz":
		encode = tar_xz.New(packedName)
	case "tar.bz":
		encode = tar_bz2.New(packedName)

	default:
		cmd.PrintErr("unknown method")
	}

	err := encode.Encode([]string{filePath})
	if err != nil {
		handleErr(err)
	}
}

func packedFileName(path string, packedExtension string) string {
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
