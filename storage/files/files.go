package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/kavshevnova/telegrambot_tuzov/lib/e"
	"github.com/kavshevnova/telegrambot_tuzov/storage"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774 //у всех пользователей будут права на чтение и запись
//Permissions - права доступа к файлам и системные ограничения

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save page", err) }()
	//Формируем путь до директории куда будет сохраняться файл
	filePath := filepath.Join(s.basePath, page.Username)
	//s.basePath — путь из storage  (например, "/data/users").
	//page.Username — имя пользователя (например, "john_doe").
	//filepath.Join() корректно соединяет пути с учётом ОС (добавляет / или \).
	//Итоговый путь: "/data/users/john_doe"
	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		//Создаёт все нужные директории в этом пути.
		//Если директории users не существует — она будет создана (например user)
		//Рекурсивно создаёт все недостающие папки в пути.
		//Задаёт права доступа 0774
		return err
	}
	fName, err := fileName(page)
	//формируем имя файла
	if err != nil {
		return err
	}
	filePath = filepath.Join(filePath, fName)
	//дописываем имя файла к пути

	file, err := os.Create(filePath)
	//создаем файл с помощью функции которой передаем путь до файла
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	//выполняем сериализацию (преобразование в бинарный формат структуры page и запись ее в файл c использованием GOB (Go Binary))
	//записываем в файл нашу страницу в нужном формате
	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}
	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random page", err) }()
	path := filepath.Join(s.basePath, userName)
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))
	file := files[n]
	return s.decodePage(filepath.Join(path, file.Name()))

}

func (s Storage) PickAll(userName string) (pages []*storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick all pages", err) }()
	path := filepath.Join(s.basePath, userName)
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	for _, file := range files {
		page, err := s.decodePage(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}
	return pages, nil
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.WrapIfErr("can't remove page", err)
	}
	path := filepath.Join(s.basePath, fileName)
	if err := os.Remove(path); err != nil {
		return e.WrapIfErr(fmt.Sprintf("can't remove page %s", path), err)
	}
	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.WrapIfErr("can't check if file %s exists", err)
	}
	path := filepath.Join(s.basePath, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, e.WrapIfErr(fmt.Sprintf("can't check if file %s exists", path), err)
	}
	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
