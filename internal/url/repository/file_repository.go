package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/pkg/logger"
)

type FileRepository struct {
	file       *os.File
	cache      map[string]*models.URL
	logger     *logger.Logger
	mu         sync.RWMutex
	isSaveMode bool
}

func NewFileRepository(filename string, logger *logger.Logger) (*FileRepository, error) {
	// empty filename -> disable saving
	if filename == "" {
		return &FileRepository{
			file:       nil,
			cache:      make(map[string]*models.URL),
			logger:     logger,
			mu:         sync.RWMutex{},
			isSaveMode: false,
		}, nil
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Errorf("failed to open file %v: %v", filename, err)
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	fr := &FileRepository{
		file:       file,
		cache:      make(map[string]*models.URL),
		logger:     logger,
		mu:         sync.RWMutex{},
		isSaveMode: true,
	}

	if err := fr.initCache(); err != nil {
		logger.Errorf("failed to init cache %v: %v", filename, err)
		return nil, fmt.Errorf("failed to init cache: %w", err)
	}

	return fr, nil
}

func (fr *FileRepository) AddURL(url, shortURL string) (string, error) {
	fr.mu.RLock()
	if got, ok := fr.cache[url]; ok {
		fr.logger.Infof("cache hit for %s: %v", url, got)
		fr.mu.RUnlock()
		return got.ShortURL, nil
	}
	fr.mu.RUnlock()

	fr.mu.Lock()
	defer fr.mu.Unlock()

	u := models.NewURL(url, shortURL)
	fr.cache[url] = u

	if fr.isSaveMode {
		if err := fr.saveURL(u); err != nil {
			fr.logger.Errorf("failed to save url %v: %v", url, err)
			return "", fmt.Errorf("failed to save url: %w", err)
		}
	}

	return shortURL, nil
}

func (fr *FileRepository) GetURL(shortURL string) (string, bool) {
	fr.mu.RLock()
	defer fr.mu.RUnlock()

	for k, v := range fr.cache {
		if v.ShortURL == shortURL {
			return k, true
		}
	}

	return "", false
}

func (fr *FileRepository) initCache() error {
	file, err := os.OpenFile(fr.file.Name(), os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		fr.logger.Errorf("failed to open file %v: %v", fr.file.Name(), err)
		return err
	}

	// each model should be on new line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var u models.URL
		if err := json.Unmarshal(scanner.Bytes(), &u); err != nil {
			fr.logger.Errorf("failed to unmarshal json: %v", err)
			return err
		}
		fr.cache[u.BaseURL] = &u
	}

	if err := scanner.Err(); err != nil {
		fr.logger.Errorf("failed to scan file %v: %v", fr.file.Name(), err)
		return err
	}

	return nil
}

func (fr *FileRepository) saveURL(url *models.URL) error {
	data, err := json.Marshal(url)
	if err != nil {
		fr.logger.Errorf("failed to marshal url %v: %v", url, err)
		return err
	}

	// write url model with new linea
	file, err := os.OpenFile(fr.file.Name(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fr.logger.Errorf("failed to open file %v: %v", fr.file.Name(), err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(string(data) + "\n")
	if err != nil {
		fr.logger.Errorf("failed to write file %v: %v", fr.file.Name(), err)
		return err
	}

	return nil
}
