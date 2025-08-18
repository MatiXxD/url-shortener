package usecase

import "errors"

var (
	ErrNoBatchShorten         = errors.New("failed to create shorten urls for all batch")
	ErrSomeBatchShortenFailed = errors.New("failed to create shorten urls for part of the batch")
)
