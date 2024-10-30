package repository

import (
	"os"
	"testing"

	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/config"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/logger"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

func TestStateManager_LoadFromFile(t *testing.T) {
	t.Parallel()

	testRepoState := CreateURLRepositoryState(
		map[string]string{
			"XXXYZZZZ": "http://ya.com",
			"XXXYYZZZ": "http://ya.com",
			"XXXYYYZZ": "http://ya.com",
		},
	)

	t.Run("test save/load service state", func(t *testing.T) {
		t.Parallel()

		config := config.CreateDefaultConfig()

		config.FileStoragePath = "records.json"
		defer os.Remove(config.FileStoragePath)

		manager := CreateStateManager(config, *logger.Log)

		err := manager.SaveToFile(testRepoState)
		require.NoError(t, err)

		repoState, err := manager.LoadFromFile()
		require.NoError(t, err)
		assert.Equal(t, len(testRepoState.GetURLRepositoryState()), len(repoState.GetURLRepositoryState()))
	})
}
