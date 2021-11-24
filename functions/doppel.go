package functions

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/jheuel/asar"
)

type AsarBoundaries struct {
	Name     string
	IsDir    bool
	Flag     asar.Flag
	Parent   *AsarBoundaries
	Children []*AsarBoundaries
	Content  []byte
}

func Decode(path string) (*AsarBoundaries, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Could not open file: %v", err)
	}
	defer f.Close()

	archive, err := asar.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("Could not decode archive: %v", err)
	}

	return toMemory(archive), nil
}

// Brings all file to the memory
func toMemory(e *asar.Entry) *AsarBoundaries {
	n := &AsarBoundaries{}
	n.Name = e.Name
	n.IsDir = e.Flags&asar.FlagDir != 0
	n.Content = e.Bytes()
	n.Flag = e.Flags
	for _, c := range e.Children {
		child := toMemory(c)
		child.Parent = n
		n.Children = append(n.Children, child)
	}
	return n
}

// populate -> the directory with shining new asar file, replacing the modified one.
func populate(n *AsarBoundaries, entries *asar.Builder) {
	for _, c := range n.Children {
		if c.IsDir {
			e := entries.AddDir(c.Name, asar.FlagDir)
			populate(c, e)
			entries.Parent()
		} else {
			entries.Add(c.Name, bytes.NewReader(c.Content), int64(len(c.Content)), c.Flag)
		}
	}
}

// EncodeTo ->  Build the fale again, then @populate() it
func EncodeTo(archive *AsarBoundaries, asarFileName string) error {
	asarArchive, err := os.Create(asarFileName)
	if err != nil {
		return fmt.Errorf("Não pude abrir o arquivo: %v", err)
	}
	defer asarArchive.Close()

	entries := asar.Builder{}

	populate(archive, &entries)
	if _, err := entries.Root().EncodeTo(asarArchive); err != nil {
		return fmt.Errorf("could not create: %s, the error was %v", asarFileName, err)
	}
	return nil
}

func Modify(n *AsarBoundaries, pattern, value string) {
	if strings.HasSuffix(n.Name, ".js") {
			n.Content = []byte(strings.ReplaceAll(string(n.Content), pattern, value))
	}
	for _, c := range n.Children {
		Modify(c, pattern, value)
	}
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}