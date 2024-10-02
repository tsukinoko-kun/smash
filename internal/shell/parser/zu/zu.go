package zu

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"smash/internal/assert"
	"smash/internal/env"
	"sort"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
)

var (
	mut = sync.RWMutex{}
)

func getZu() (*Zu, error) {
	mut.Lock()
	defer mut.Unlock()
	file, err := os.OpenFile(env.ZuFile, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, errors.Join(errors.New("failed to open zu file"), err)
	}
	defer file.Close()

	// read the file
	z := &Zu{}
	b, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Join(errors.New("failed to read zu file"), err)
	}
	if err := proto.Unmarshal(b, z); err != nil {
		return nil, errors.Join(errors.New("failed to parse zu file"), err)
	}
	return z, nil
}

func saveZu(z *Zu) error {
	z.Sort()

	mut.RLock()
	defer mut.RUnlock()
	file, err := os.Create(env.ZuFile)
	if err != nil {
		return errors.Join(errors.New("failed to create zu file"), err)
	}
	defer file.Close()

	// write the file
	b, err := proto.Marshal(z)
	if err != nil {
		return errors.Join(errors.New("failed to marshal zu file"), err)
	}
	if _, err := file.Write(b); err != nil {
		return errors.Join(errors.New("failed to write zu file"), err)
	}
	return nil
}

func (z *Zu) FindOrCreateExact(path string) *ZuEntry {
	mut.Lock()
	defer mut.Unlock()

	for _, e := range z.Entries {
		if e.Path == path {
			return e
		}
	}
	e := &ZuEntry{Path: path, AccessCount: 0}
	z.Entries = append(z.Entries, e)
	return e
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (z *Zu) delete(path string) {
	mut.Lock()
	defer mut.Unlock()
	for i, e := range z.Entries {
		if e.Path == path {
			z.Entries = append(z.Entries[:i], z.Entries[i+1:]...)
			break
		}
	}
}

func (z *Zu) FindMostCommonMatch(target string) *ZuEntry {
	mut.RLock()
	defer mut.RUnlock()

	var best *ZuEntry
	for _, e := range z.Entries {
		if filepath.Base(e.Path) == target {
			if !exists(e.Path) {
				// remove non-existing Entries

				continue
			}
			best = e
		}
	}
	if best != nil {
		return best
	}
	// no exact match found, try partial match
	for _, e := range z.Entries {
		if strings.Contains(strings.ToLower(filepath.Base(e.Path)), strings.ToLower(target)) {
			best = e
		}
	}
	return best
}

// Sorts the entries by access count descending
func (z *Zu) Sort() {
	mut.Lock()
	defer mut.Unlock()
	sort.Slice(z.Entries, func(i, j int) bool {
		return z.Entries[i].AccessCount > z.Entries[j].AccessCount
	})
	if len(z.Entries) > 64 {
		// cut off the least used Entries
		z.Entries = z.Entries[:32]
	}
}

// Clear clears the zu
func Clear() error {
	mut.Lock()
	z := &Zu{}
	mut.Unlock()
	return saveZu(z)
}

func List() ([]string, error) {
	mut.RLock()
	defer mut.RUnlock()
	z, err := getZu()
	if err != nil {
		return nil, errors.Join(errors.New("failed to get zu"), err)
	}
	assert.NotNil(z, "*Zu is nil")

	var paths []string
	for _, e := range z.Entries {
		paths = append(paths, fmt.Sprintf("%s (%d)", e.Path, e.AccessCount))
	}
	return paths, nil
}

func To(target string) error {
	// get the current zu
	z, err := getZu()
	if err != nil {
		return errors.Join(errors.New("failed to get zu"), err)
	}
	assert.NotNil(z, "*Zu is nil")

	// look if the target exists
	stat, err := os.Stat(target)
	if err == nil {
		if stat.IsDir() {
			// add the target to the zu
			abs, err := filepath.Abs(target)
			if err != nil {
				return errors.Join(errors.New("failed to get absolute path"), err)
			}
			e := z.FindOrCreateExact(abs)
			if e.AccessCount < math.MaxInt64 {
				e.AccessCount++
			}
			_ = saveZu(z)
			// change to the target directory
			return os.Chdir(abs)
		} else {
			return errors.New("target exists but is not a directory")
		}
	}

	// find matching entry in zu
	e := z.FindMostCommonMatch(target)
	if e == nil {
		return errors.New("no matching entry found")
	}

	// change to the target directory
	e.AccessCount++
	_ = saveZu(z)
	return os.Chdir(e.Path)
}
