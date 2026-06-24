package web

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"sort"
)

func GetStaticFilesHash() (string, error) {
	var files []string
	err := fs.WalkDir(Assets, "static", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	sort.Strings(files)

	h := sha256.New()
	for _, f := range files {
		data, err := Assets.ReadFile(f)
		if err != nil {
			return "", err
		}

		h.Write([]byte(f))
		h.Write([]byte{0})
		h.Write(data)
		h.Write([]byte{0})
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
