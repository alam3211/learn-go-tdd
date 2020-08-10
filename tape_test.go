package main

import (
	"io/ioutil"
	"testing"
)

func TestTape_Write(t *testing.T) {
	file, clean := createTempFile(t, "1234")
	defer clean()
	t.Run("it records win when POST", func(t *testing.T) {
		tape := &Tape{file}
		tape.Write([]byte("ABC"))

		file.Seek(0, 0)
		newFileContents, _ := ioutil.ReadAll(file)
		got := string(newFileContents)
		expected := "ABC"

		if got != expected {
			t.Errorf("Write file failed got %s, expected %s", got, expected)
		}

	})
}
