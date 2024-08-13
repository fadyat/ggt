package main

import (
	"io"
	"log/slog"
)

type writerCreator func(string) (io.WriteCloser, error)

func renderPackage(pkgInfo *packageInfo, wcreator writerCreator) {
	for path, fi := range pkgInfo.files {
		if err := renderFile(path, fi, wcreator); err != nil {
			slog.Error("rendering file", slog.String("path", path), slog.Any("error", err))
		} else {
			slog.Info("generated tests", slog.String("path", path))
		}
	}
}

func renderFile(filePath string, fileInfo *fileInfo, wcreator writerCreator) error {
	// out, err := wcreator(filePath)
	//if err != nil {
	//	return fmt.Errorf("creating file: %w", err)
	//}
	//defer out.Close()

	// ...

	return nil
}
