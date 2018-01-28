package slog

import "testing"
import "time"

func TestNewLog(t *testing.T) {
	log, err := NewLog("example.txt", 0, 0, Snappy)
	if err != nil {
		t.Fatal(err)
	}
	defer log.Close()

	log.Info("Yee, whats up!")
	time.Sleep(time.Second)
	log.Compress()

	log.Warning("Hehe. Check log file man")
	time.Sleep(time.Second)
	log.Compress()

	time.Sleep(time.Second)
}
