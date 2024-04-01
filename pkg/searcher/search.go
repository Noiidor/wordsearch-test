package searcher

import (
	"bufio"
	"errors"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"wordsearch/pkg/dir"
	scanextension "wordsearch/pkg/scan-extension"
)

type Searcher struct {
	FS fs.FS
}

func (s *Searcher) Search(word string) ([]string, error) {
	word = strings.ToLower(word)

	filePaths, err := dir.FilesFS(s.FS, "")
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(filePaths))

	// Канал чтоб ошибки в горутинах не пропадали в бездне
	errChan := make(chan error, len(filePaths))
	var wg sync.WaitGroup
	for _, v := range filePaths {
		wg.Add(1)
		go func() {
			defer wg.Done()

			file, err := s.FS.Open(v)
			if err != nil {
				errChan <- err
			}

			// Разделяет текст по отдельным словам(знаки игнорируются) и записывает в мапу как ключи
			scanner := bufio.NewScanner(file)
			scanner.Split(scanextension.ScanWordsOnly)

			words := make(map[string]any)

			for scanner.Scan() {
				words[strings.ToLower(scanner.Text())] = nil
			}

			// Поиск по хэшу О(1)
			if _, ok := words[word]; ok {
				info, err := file.Stat()
				if err != nil {
					errChan <- err
				}
				fullname := info.Name()
				name := strings.TrimSuffix(fullname, filepath.Ext(fullname))
				result = append(result, name)
			}
		}()
	}

	wg.Wait()

	if len(errChan) == 0 {
		sort.Strings(result)
		return result, nil
	}

	asyncErrors := make([]error, 0, len(filePaths))
	for v := range errChan {
		asyncErrors = append(asyncErrors, v)
	}
	return nil, errors.Join(asyncErrors...)
}
