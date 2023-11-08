package file

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/Asliddin/zoomda/configs"
	"github.com/Asliddin/zoomda/models"
)

type FilesService struct {
	cfg configs.Config
}

func NewFilesService(cfg configs.Config) *FilesService {
	return &FilesService{cfg: cfg}
}

func (f *FilesService) Save(ctx context.Context, file models.File) (string, error) {
	src, err := file.File.Open()
	if err != nil {
		return "", err
	}
	src.Close()
	createPath := f.cfg.StaticFilePath + file.Path
	fmt.Println("path", createPath)
	if _, err := os.Stat(createPath); os.IsNotExist(err) {
		err = os.Mkdir(createPath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	pattern := regexp.MustCompile("\\.[0-9a-z]+$")

	extension := pattern.FindString(file.File.Filename)
	if extension == "" {
		return "", models.ErrFileName
	}

	now := strconv.FormatInt(time.Now().UnixNano(), 10)

	newName := now + extension

	dst := createPath + "/" + newName

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}

	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return "", err
	}

	return newName, nil
}

func (f *FilesService) Delete(ctx context.Context, path, filename string) error {
	err := os.Remove(f.cfg.StaticFilePath + path + "/" + filename)
	if err != nil {
		return err
	}
	return nil
}
