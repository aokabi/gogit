/*
.gitディレクトリ内のファイルやディレクトリの操作をする
.gitディレクトリが親ディレクトリや更に上のディレクトリにあっても機能するようにする
*/
package pkg

import (
	"os"
	"path/filepath"
)

const (
	gitDirName = ".git"
)

// ルートまで遡ってgit projectのディレクトリを探す
// .gitディレクトリのpathを返す
func findGitDir() (string, error) {
	// get current dir
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// find .git directory
	currentDir := workingDir
	for {
		if _, err := os.Stat(filepath.Join(currentDir, gitDirName)); err == nil {
			return filepath.Join(currentDir, gitDirName), nil
		}
		// rootまで辿っても見つからなかった
		if currentDir == "/" {
			return "", nil
		}
		// 探すディレクトリを更新
		currentDir = filepath.Dir(currentDir)
	}
}

// git projectの中にディレクトリをつくる
func CreateDir(path string) error {
	gitDir, err := findGitDir()
	if err != nil {
		return err
	}

	return os.Mkdir(filepath.Join(gitDir, path), 0755)
}

func CreateFile(path string) (*os.File, error) {
	gitDir, err := findGitDir()
	if err != nil {
		return nil, err
	}

	return os.Create(filepath.Join(gitDir, path))
}

func IsNotExist(path string) bool {
	gitDir, err := findGitDir()
	if err != nil {
		return true
	}

	if _, err := os.Stat(filepath.Join(gitDir, path)); os.IsNotExist(err) {
		return true
	} else {
		return false
	}
}

func Open(path string) (*os.File, error) {
	gitDir, err := findGitDir()
	if err != nil {
		return nil, err
	}

	return os.Open(filepath.Join(gitDir, path))
}

func OpenFile(path string, flag int, mode os.FileMode) (*os.File, error) {
	gitDir, err := findGitDir()
	if err != nil {
		return nil, err
	}

	return os.OpenFile(filepath.Join(gitDir, path), flag, mode)
}

func Truncate(path string, size int64) error {
	gitDir, err := findGitDir()
	if err != nil {
		return err
	}

	return os.Truncate(filepath.Join(gitDir, path), size)
}