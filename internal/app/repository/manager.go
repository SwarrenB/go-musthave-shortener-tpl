package repository

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/config"
	easyjson "github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type StateManager struct {
	config *config.Config
	log    zap.Logger
}

type FileReader struct {
	file    *os.File
	scanner *bufio.Scanner
	log     zap.Logger
}

type FileWriter struct {
	file *os.File
	log  zap.Logger
}

func CreateStateManager(config *config.Config, log zap.Logger) *StateManager {
	return &StateManager{
		log:    log,
		config: config,
	}
}

func (sm *StateManager) LoadFromFile() (*URLRepositoryState, error) {
	reader, err := NewReader(sm.config.FileStoragePath, sm.log)
	if err != nil {
		return nil, err
	}
	defer reader.file.Close()

	state, err := reader.LoadState()
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (sm *StateManager) SaveToFile(state *URLRepositoryState) error {
	FileWriter, err := NewWriter(sm.config.FileStoragePath, sm.log)
	if err != nil {
		return err
	}
	defer FileWriter.file.Close()

	err = FileWriter.SaveState(state)
	if err != nil {
		return err
	}

	return nil
}

func NewReader(fileName string, log zap.Logger) (*FileReader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open %v error: %s", fileName, err)
	}

	return &FileReader{
		file:    file,
		scanner: bufio.NewScanner(file),
		log:     log,
	}, nil
}

func (reader *FileReader) LoadState() (*URLRepositoryState, error) {
	state := make(map[string]string)

	if err := reader.Reset(); err != nil {
		return nil, err
	}
	reader.scanner = bufio.NewScanner(reader.file)

	for reader.scanner.Scan() {
		data := reader.scanner.Bytes()
		record := FileRecord{}

		err := record.UnmarshalJSON(data)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling error: %s", err)
		}

		state[record.ShortURL] = record.OriginalURL
		reader.log.Info("record loaded",
			zap.Int("uuid", record.ID),
			zap.String("short_url", record.ShortURL),
			zap.String("original_url", record.OriginalURL))
	}

	if reader.scanner.Err() == nil {
		return CreateURLRepositoryState(state), nil
	}

	return nil, fmt.Errorf("error reading file: %s", reader.scanner.Err())
}

func NewWriter(fileName string, log zap.Logger) (*FileWriter, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, fmt.Errorf("error opening %v: %s", fileName, err)
	}

	return &FileWriter{
		file: file,
		log:  log,
	}, nil
}

func (FileWriter *FileWriter) SaveState(state *URLRepositoryState) error {
	urls := state.GetURLRepositoryState()

	index := 1
	for shortURL, OriginalURL := range urls {
		record := &FileRecord{
			ID:          index,
			ShortURL:    shortURL,
			OriginalURL: OriginalURL,
		}

		if _, err := easyjson.MarshalToWriter(record, FileWriter.file); err != nil {
			return fmt.Errorf("marshalling error: %s", err)
		}

		if _, err := FileWriter.file.WriteString("\n"); err != nil {
			return fmt.Errorf("write file error: %s", err)
		}

		FileWriter.log.Info("record preserved",
			zap.Int("uuid", record.ID),
			zap.String("short_url", record.ShortURL),
			zap.String("original_url", record.OriginalURL),
		)

		index++
	}

	return nil
}

func (r *FileReader) Reset() error {
	if _, err := r.file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("file reader reset: %w", r.scanner.Err())
	}

	r.scanner = bufio.NewScanner(r.file)

	return nil
}
