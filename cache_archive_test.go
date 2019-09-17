package main

import (
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/pathutil"
)

type nopWriteCloser struct{}

func (writer nopWriteCloser) Write(b []byte) (int, error) {
	return len(b), nil
}

func (writer nopWriteCloser) Close() error {
	return nil
}

func TestNewArchive(t *testing.T) {
	tests := []struct {
		name     string
		compress bool
		wantGzip bool
		wantErr  bool
	}{
		{
			name:     "no compress",
			compress: false,
			wantGzip: false,
			wantErr:  false,
		},
		{
			name:     "compress",
			compress: true,
			wantGzip: true,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var writer nopWriteCloser
			got, err := NewArchive(writer, tt.compress)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewArchive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			hasGzip := got != nil && got.gzip != nil
			if tt.wantGzip != hasGzip {
				t.Errorf("NewArchive() has gzip = %v, want %v", hasGzip, tt.wantGzip)
			}
		})
	}
}

func TestArchive_Write(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("cache")
	if err != nil {
		t.Fatalf("failed to create tmp dir: %s", err)
	}

	fileToArchive := filepath.Join(tmpDir, "file")
	createDirStruct(t, map[string]string{fileToArchive: ""})

	t.Log("no compress")
	{
		var writer nopWriteCloser
		archive, err := NewArchive(writer, false)
		if err != nil {
			t.Fatalf("failed to create archive: %s", err)
		}

		if err := archive.Write([]string{fileToArchive}, false); err != nil {
			t.Fatalf("failed to write archive: %s", err)
		}
	}

	t.Log("compress")
	{
		var writer nopWriteCloser
		archive, err := NewArchive(writer, true)
		if err != nil {
			t.Fatalf("failed to create archive: %s", err)
		}

		if err := archive.Write([]string{fileToArchive}, false); err != nil {
			t.Fatalf("failed to write archive: %s", err)
		}
	}
}

func TestArchive_WriteHeader(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("cache")
	if err != nil {
		t.Fatalf("failed to create tmp dir: %s", err)
	}

	fileToArchive := filepath.Join(tmpDir, "file")
	createDirStruct(t, map[string]string{fileToArchive: ""})

	var writer nopWriteCloser
	archive, err := NewArchive(writer, false)
	if err != nil {
		t.Fatalf("failed to create archive: %s", err)
	}

	if err := archive.WriteHeader(map[string]string{"file/to/cache": "indicator/file"}, cacheInfoFilePath); err != nil {
		t.Fatalf("failed to write archive header: %s", err)
	}
}

func TestArchive_Close(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("cache")
	if err != nil {
		t.Fatalf("failed to create tmp dir: %s", err)
	}

	fileToArchive := filepath.Join(tmpDir, "file")
	createDirStruct(t, map[string]string{fileToArchive: ""})

	t.Log("no compress")
	{
		var writer nopWriteCloser
		archive, err := NewArchive(writer, false)
		if err != nil {
			t.Fatalf("failed to create archive: %s", err)
		}

		if err := archive.Write([]string{fileToArchive}, false); err != nil {
			t.Fatalf("failed to write archive: %s", err)
		}

		if err := archive.Close(); err != nil {
			t.Fatalf("failed to close archive: %s", err)
		}
	}

	t.Log("compress")
	{
		var writer nopWriteCloser
		archive, err := NewArchive(writer, true)
		if err != nil {
			t.Fatalf("failed to create archive: %s", err)
		}

		if err := archive.Write([]string{fileToArchive}, false); err != nil {
			t.Fatalf("failed to write archive: %s", err)
		}

		if err := archive.Close(); err != nil {
			t.Fatalf("failed to close archive: %s", err)
		}
	}
}
