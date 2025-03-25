package repository

// import (
// 	"os"
// 	"testing"

// 	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/config"
// 	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/logger"
// 	"github.com/go-playground/assert/v2"
// 	"github.com/stretchr/testify/require"
// )

// func TestStateManager_LoadFromFile(t *testing.T) {
// 	t.Parallel()
// 	//TODO it will be fixed later
// 	testRepoState := CreateURLRepositoryState(
// 		map[string]Record{
// 			"XXXYZZZZ": {0, "XXXYZZZZ", "ZZZZXXXYYY", "1"},
// 			"XXXYYZZZ": {1, "XXXYZZZZ", "ZZZZXXXYYY", "1"},
// 			"XXXYYYZZ": {2, "XXXYZZZZ", "ZZZZXXXYYY", "1"},
// 		},
// 	)

// 	t.Run("test save/load service state", func(t *testing.T) {
// 		t.Parallel()

// 		config := config.CreateDefaultConfig()

// 		config.FileStoragePath = "records.json"
// 		defer os.Remove(config.FileStoragePath)

// 		manager := CreateStateManager(config, logger.CreateLogger("Info").GetLogger())

// 		err := manager.SaveToFile(testRepoState)
// 		require.NoError(t, err)

// 		repoState, err := manager.LoadFromFile()
// 		require.NoError(t, err)
// 		assert.Equal(t, len(testRepoState.GetURLRepositoryState()), len(repoState.GetURLRepositoryState()))
// 	})
// }
