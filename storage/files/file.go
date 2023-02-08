package files

import (
	"Telegram_bot/lib/e"
	storage_ "Telegram_bot/storage"
	"encoding/gob"
	"errors"
	"math/rand"
	"os"
	filepath2 "path/filepath"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Remove(p *storage_.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("Can't remove a file", err)
	}

	path := filepath2.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		return e.Wrap("Can't remove a file", err)
	}

	return nil
}

func (s Storage) IsExists(p *storage_.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("Can't remove a file", err)
	}

	path := filepath2.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, e.Wrap("Can't check if file exists", err)

	}

	return true, nil

}

func (s Storage) Save(page *storage_.Page) (err error) {
	defer func() { err = e.Wrap("Can't save", err) }()

	fPath := filepath2.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath2.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil

}

func (s Storage) PickRandom(username string) (page *storage_.Page, err error) {
	defer func() { err = e.Wrap("Can't pick a random page", err) }()

	path := filepath2.Join(s.basePath, username)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage_.ErrNoSavePage
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath2.Join(path, file.Name()))

}

func (s Storage) decodePage(filePath string) (*storage_.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("Can't open a page", err)
	}
	defer func() { _ = f.Close() }()

	var p storage_.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("Can't decode a page", err)
	}

	return &p, nil
}

func fileName(p *storage_.Page) (string, error) {
	return p.Hash()
}
