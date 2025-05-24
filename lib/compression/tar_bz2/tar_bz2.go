package tar_bz2

import (
	"archive/tar"
	"fmt"
	"github.com/dsnet/compress/bzip2"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type EncodeDecoder struct {
	OutputPath string
}

func New(outPaht string) *EncodeDecoder {
	return &EncodeDecoder{
		OutputPath: outPaht,
	}
}

func (ed *EncodeDecoder) Encode(sourcePaths []string) error {
	file, err := os.Create(ed.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", ed.OutputPath, err)
	}
	defer file.Close()

	bz2Writer, err := bzip2.NewWriter(file, nil) // nil для параметрів за замовчуванням
	if err != nil {
		return fmt.Errorf("failed to create bzip2 writer: %w", err)
	}
	defer bz2Writer.Close()

	tarWriter := tar.NewWriter(bz2Writer)
	defer tarWriter.Close()

	for _, source := range sourcePaths {
		err = filepath.Walk(source, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error walking through %s: %w", filePath, err)
			}

			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return fmt.Errorf("failed to create tar header for %s: %w", filePath, err)
			}

			relPath, err := filepath.Rel(filepath.Dir(source), filePath)
			if err != nil {
				return fmt.Errorf("failed to get relative path for %s: %w", filePath, err)
			}
			header.Name = strings.ReplaceAll(relPath, string(os.PathSeparator), "/")

			if err := tarWriter.WriteHeader(header); err != nil {
				return fmt.Errorf("failed to write tar header for %s: %w", filePath, err)
			}

			if !info.IsDir() {
				file, err := os.Open(filePath)
				if err != nil {
					return fmt.Errorf("failed to open file %s: %w", filePath, err)
				}
				defer file.Close()

				_, err = io.Copy(tarWriter, file)
				if err != nil {
					return fmt.Errorf("failed to write file %s to tar.bz2: %w", filePath, err)
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (ed *EncodeDecoder) Decode(outputDir string) error {
	if outputDir == "" {
		return fmt.Errorf("output directory path is empty")
	}

	if info, err := os.Stat(outputDir); err == nil && !info.IsDir() {
		return fmt.Errorf("output path %s is a file, not a directory", outputDir)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output directory %s: %v\n", outputDir, err)
		return fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
	}

	file, err := os.Open(ed.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to open archive %s: %w", ed.OutputPath, err)
	}

	defer file.Close()
	bz2Reader, err := bzip2.NewReader(file, nil)
	if err != nil {
		return fmt.Errorf("failed to create bzip2 reader: %w", err)
	}
	defer bz2Reader.Close()
	tarReader := tar.NewReader(bz2Reader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		targetPath := filepath.Join(outputDir, header.Name)
		targetPath = filepath.Clean(targetPath)
		if strings.Contains(header.Name, "..") {
			fmt.Fprintf(os.Stderr, "Skipping potentially unsafe path: %s\n", header.Name)
			continue
		}
		if header.Typeflag == tar.TypeDir {
			fmt.Printf("Creating directory: %s\n", targetPath)
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
			continue
		}
		if header.Typeflag == tar.TypeReg {
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", targetPath, err)
			}
			targetFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", targetPath, err)
			}

			_, err = io.Copy(targetFile, tarReader)
			if err != nil {
				targetFile.Close()
				return fmt.Errorf("failed to write file %s: %w", targetPath, err)
			}

			if err := targetFile.Close(); err != nil {
				return fmt.Errorf("failed to close file %s: %w", targetPath, err)
			}
		}
	}

	return nil
}
