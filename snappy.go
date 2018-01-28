package slog

import (
	"fmt"
	"io"
	"os"

	"github.com/golang/snappy"
)

// Snappy file compression.
func cSnappy(sfile *os.File) (int64, error) {
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
		os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600,
	)
	if err != nil {
		return 0, err
	}
	defer dfile.Close()

	writer := snappy.NewWriter(dfile)
	defer writer.Close()
	defer writer.Flush()

	return io.Copy(writer, sfile)
}
