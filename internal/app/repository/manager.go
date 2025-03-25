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
	reader, err := CreateReader(sm.config.FileStoragePath, sm.log)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	state, err := reader.LoadState()

	for k, v := range state.state {
		sm.log.Info("state to load from file", zap.String("shortUrl", k), zap.String("origUrl", v.OriginalURL))
	}
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (sm *StateManager) SaveToFile(state *URLRepositoryState) error {
	writer, err := CreateWriter(sm.config.FileStoragePath, sm.log)
	if err != nil {
		return err
	}
	defer writer.Close()
	for k, v := range state.state {
		sm.log.Info("state to save", zap.String("shortUrl", k), zap.String("origUrl", v.OriginalURL))
	}
	err = writer.SaveState(state)
	if err != nil {
		return err
	}

	return nil
}

func CreateReader(fileName string, log zap.Logger) (*FileReader, error) {
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
	state := make(map[string]Record)

	if err := reader.Reset(); err != nil {
		return nil, err
	}

	for reader.scanner.Scan() {
		data := reader.scanner.Bytes()
		record := Record{}

		err := record.UnmarshalJSON(data)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling error: %s", err)
		}

		state[record.ShortURL] = record
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

func CreateWriter(fileName string, log zap.Logger) (*FileWriter, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, fmt.Errorf("error opening %v: %s", fileName, err)
	}

	return &FileWriter{
		file: file,
		log:  log,
	}, nil
}

func (writer *FileWriter) SaveState(state *URLRepositoryState) error {
	urls := state.GetURLRepositoryState()

	buffWriter := bufio.NewWriter(writer.file)
	defer buffWriter.Flush()

	for _, originalURLRecord := range urls {
		record := originalURLRecord

		if _, err := easyjson.MarshalToWriter(record, buffWriter); err != nil {
			return fmt.Errorf("marshalling error: %s", err)
		}

		if _, err := buffWriter.WriteString("\n"); err != nil {
			return fmt.Errorf("write file error: %s", err)
		}

		writer.log.Info("record preserved",
			zap.Int("uuid", record.ID),
			zap.String("short_url", record.ShortURL),
			zap.String("original_url", record.OriginalURL),
		)

	}

	return nil
}

func (reader *FileReader) Reset() error {
	if _, err := reader.file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("file reader reset: %w", reader.scanner.Err())
	}

	reader.scanner = bufio.NewScanner(reader.file)

	return nil
}

func (reader *FileReader) Close() error {
	err := reader.file.Close()
	if err != nil {
		return fmt.Errorf("error closing the file: %s", err)
	}

	return nil
}

func (writer *FileWriter) Close() error {
	err := writer.file.Close()
	if err != nil {
		return fmt.Errorf("error closing the file: %s", err)
	}

	return nil
}
