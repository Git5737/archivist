package cmd

import (
	"archivist/lib/compression"
	"archivist/lib/compression/tar"
	"archivist/lib/compression/tar_bz2"
	"archivist/lib/compression/tar_gz"
	"archivist/lib/compression/tar_xz"
	"archivist/lib/compression/zip"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var unpackcmd = &cobra.Command{
	Use:   "unpack",
	Short: "Unpack file",
	Run:   unpack,
}

func unpack(cmd *cobra.Command, args []string) {
	archivePath := args[0]
	if archivePath == "" {
		fmt.Errorf("archive path is not specified")
	}

	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		fmt.Errorf("archive %s does not exist: %w", archivePath, err)
	}

	var outputDir string
	if outputDir == "" {
		base := strings.TrimSuffix(filepath.Base(archivePath), filepath.Ext(archivePath))
		base = strings.TrimSuffix(base, ".tar")
		outputDir = filepath.Join(filepath.Dir(archivePath), base)
	}

	if info, err := os.Stat(outputDir); err == nil && !info.IsDir() {
		fmt.Errorf("output path %s is a file, not a directory", outputDir)
	}

	method, err := cmd.Flags().GetString("method")
	if err != nil {
		fmt.Errorf("failed to get method flag: %w", err)
	}

	if method == "" {
		switch {
		case strings.HasSuffix(archivePath, ".zip"):
			method = "zip"
		case strings.HasSuffix(archivePath, ".tar"):
			method = "tar"
		case strings.HasSuffix(archivePath, ".tar.gz"):
			method = "tar.gz"
		case strings.HasSuffix(archivePath, ".tar.bz"):
			method = "tar.bz"
		case strings.HasSuffix(archivePath, ".tar.xz"):
			method = "tar.xz"
		default:
			fmt.Errorf("cannot determine compression method from file ext ension: %s", archivePath)
		}
	}

	var decode compression.Decoder
	switch method {
	case "zip":
		decode = zip.New(archivePath)
	case "tar":
		decode = tar.New(archivePath)
	case "tar.gz":
		decode = tar_gz.New(archivePath)
	case "tar.bz2":
		decode = tar_bz2.New(archivePath)
	case "tar.xz":
		decode = tar_xz.New(archivePath)
	case "tar.bz":
		decode = tar_bz2.New(archivePath)
	default:
		fmt.Errorf("unknown compression method: %s", method)
	}

	err = decode.Decode(outputDir)
	if err != nil {
		fmt.Errorf("failed to decode %s: %w", archivePath, err)
	}
}

func init() {
	rootCmd.AddCommand(unpackcmd)

	unpackcmd.Flags().StringP("method", "m", "", "decompression method: vlc")
	if err := packcmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}
}
