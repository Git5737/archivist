package zip

import (
	"archive/zip"
	"fmt"
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
	zipFile, err := os.Create(ed.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file %s: %w", ed.OutputPath, err)
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	for _, source := range sourcePaths {
		err = filepath.Walk(source, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error walking through %s: %w", filePath, err)
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return fmt.Errorf("failed to create zip header for %s: %w", filePath, err)
			}

			relPath, err := filepath.Rel(filepath.Dir(source), filePath)
			if err != nil {
				return fmt.Errorf("failed to get relative path for %s: %w", filePath, err)
			}
			header.Name = strings.ReplaceAll(relPath, string(os.PathSeparator), "/")

			header.Method = zip.Deflate

			if info.IsDir() {
				header.Name += "/"
			} else {
				writer, err := archive.CreateHeader(header)
				if err != nil {
					return fmt.Errorf("failed to create zip entry for %s: %w", filePath, err)
				}

				file, err := os.Open(filePath)
				if err != nil {
					return fmt.Errorf("failed to open file %s: %w", filePath, err)
				}
				defer file.Close()

				_, err = io.Copy(writer, file)
				if err != nil {
					return fmt.Errorf("failed to write file %s to zip: %w", filePath, err)
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

	reader, err := zip.OpenReader(ed.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to open zip archive %s: %w", ed.OutputPath, err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		targetPath := filepath.Join(outputDir, file.Name)
		targetPath = filepath.Clean(targetPath)

		if strings.Contains(file.Name, "..") {
			fmt.Fprintf(os.Stderr, "Skipping potentially unsafe path: %s\n", file.Name)
			continue
		}

		if file.FileInfo().IsDir() {
			fmt.Printf("Creating directory: %s\n", targetPath)
			if err := os.MkdirAll(targetPath, file.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
			continue
		}

		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %s in zip: %w", file.Name, err)
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			rc.Close()
			return fmt.Errorf("failed to create parent directory for %s: %w", targetPath, err)
		}

		targetFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
		if err != nil {
			rc.Close()
			return fmt.Errorf("failed to create file %s: %w", targetPath, err)
		}
		
		_, err = io.Copy(targetFile, rc)
		if err != nil {
			rc.Close()
			targetFile.Close()
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}

		rc.Close()
		if err := targetFile.Close(); err != nil {
			return fmt.Errorf("failed to close file %s: %w", targetPath, err)
		}
	}

	return nil
}
