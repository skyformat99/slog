package slog

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func cGzip(sfile *os.File) (int64, error) {
	file := fmt.Sprintf("%s.snappy", sfile.Name())

	for i := 0; ; i++ {
		_, err := os.Stat(file)
		if err != nil {
			break
		}
		file = fmt.Sprintf("%s.snappy.%d", sfile.Name(), i)
	}

	dfile, err := os.OpenFile(
		file,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600,
	)
	if err != nil {
		return 0, nil
	}
	defer dfile.Close()

	writer, err := gzip.NewWriterLevel(dfile, gzip.BestCompression)
	if err != nil {
		return 0, err
	}
	defer writer.Close()
	defer writer.Flush()

	return io.Copy(writer, sfile)
}
