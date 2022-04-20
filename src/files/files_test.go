package files

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsts(t *testing.T) {
	assert.Equal(t, "[%s](https://%s/%s/%s/-/commit/%s)", changeLogCommitHashLinkFormat)
}

func TestVars(t *testing.T) {
	assert.Equal(t, "__version__", versionPyVariable)
	assert.Equal(t, "CHANGELOG.md", changeLogDefaultFile)
	assert.Equal(t, "setup.py", setupPythonDefaultFile)
}

func TestExists(t *testing.T) {

	t.Run("File exists", func(t *testing.T) {
		file := New(File{
			OriginPath: "mock/mock_file.txt",
		})

		assert.True(t, file.Exists())
	})

	t.Run("File does not exists", func(t *testing.T) {
		file := New(File{
			OriginPath: "mock/mock_file_404.txt",
		})

		assert.False(t, file.Exists())
	})
}

func TestOpenFile(t *testing.T) {
	t.Run("Successfully opened", func(t *testing.T) {
		file := New(File{
			OriginPath: "mock/mock_file.txt",
		})

		openedFile, err := file.OpenFile()
		assert.NoError(t, err)
		assert.NotNil(t, openedFile)
	})

	t.Run("File does not exists", func(t *testing.T) {
		file := New(File{
			OriginPath: "mock/mock_file_404.txt",
		})

		openedFile, err := file.OpenFile()
		assert.Error(t, err)
		assert.Nil(t, openedFile)
	})
}

func TestUpgradeVersionInSetupPyFile(t *testing.T) {

	t.Run("Error while oppening file", func(t *testing.T) {
		file := File{
			OriginPath:        "mock/setup_404.py",
			OutputPath:        "mock/setup_404.py",
			NewReleaseVersion: "1.0.1",
		}
		filesVersion := New(file)

		err := filesVersion.UpgradeVersionInSetupPyFile()
		assert.Error(t, err)
		assert.Equal(t, "no such file or directory", err.Error())
	})

	t.Run("Error while writing to file", func(t *testing.T) {
		file := File{
			OriginPath:        "mock/setup_mock.py",
			OutputPath:        "mock/test/setup_mock_404.py",
			NewReleaseVersion: "1.0.1",
		}
		filesVersion := New(file)

		err := filesVersion.UpgradeVersionInSetupPyFile()
		assert.Error(t, err)
		assert.Equal(t, "open mock/test/setup_mock_404.py: no such file or directory", err.Error())
	})

	t.Run("No error", func(t *testing.T) {
		file := File{
			OriginPath:        "mock/setup_mock.py",
			OutputPath:        "mock/setup_mock.py",
			NewReleaseVersion: "1.0.1",
		}
		filesVersion := New(file)
		err := filesVersion.UpgradeVersionInSetupPyFile()
		assert.NoError(t, err)

		file.NewReleaseVersion = "1.0.0"
		filesVersion = New(file)
		err = filesVersion.UpgradeVersionInSetupPyFile()
		assert.NoError(t, err)
	})
}
