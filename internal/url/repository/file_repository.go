package repository

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

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

func (fr *FileRepository) AddURL(ctx context.Context, shortenURL *models.URL) (string, error) {
	fr.mu.RLock()
	if got, ok := fr.cache[shortenURL.BaseURL]; ok {
		fr.logger.Infof("cache hit for %s: %v", shortenURL.BaseURL, got)
		fr.mu.RUnlock()
		return got.ShortURL, nil
	}
	fr.mu.RUnlock()

	fr.mu.Lock()
	defer fr.mu.Unlock()

	url := &models.URL{
		CorrelationID: shortenURL.CorrelationID,
		BaseURL:       shortenURL.BaseURL,
		ShortURL:      shortenURL.ShortURL,
		CreateAt:      time.Now(),
	}

	fr.cache[shortenURL.BaseURL] = url

	if fr.isSaveMode {
		if err := fr.saveURL(url); err != nil {
			fr.logger.Errorf("failed to save url %s: %v", url.BaseURL, err)
			return "", fmt.Errorf("failed to save url: %w", err)
		}
	}

	return shortenURL.ShortURL, nil
}

func (mr *FileRepository) BatchAddURL(ctx context.Context, urls []*models.URL) ([]*models.URL, error) {
	res := make([]*models.URL, 0, len(urls))

	for _, u := range urls {
		shortUrl, err := mr.AddURL(ctx, u)
		if err != nil {
			return res, fmt.Errorf("failed to add url=%s: %w", u.BaseURL, err)
		}

		res = append(res, &models.URL{
			CorrelationID: u.CorrelationID,
			BaseURL:       u.BaseURL,
			ShortURL:      shortUrl,
		})
	}

	return res, nil
}

func (fr *FileRepository) GetURL(ctx context.Context, shortURL string) (*models.URL, error) {
	fr.mu.RLock()
	defer fr.mu.RUnlock()

	for _, v := range fr.cache {
		if v.ShortURL == shortURL {
			return &models.URL{
				CorrelationID: v.CorrelationID,
				BaseURL:       v.BaseURL,
				ShortURL:      v.ShortURL,
			}, nil
		}
	}

	return nil, fmt.Errorf("url was not found")
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
