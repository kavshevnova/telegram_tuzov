package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kavshevnova/telegrambot_tuzov/storage"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	//уточняем с какой базой данных нам предстоит работать и передаем путь до файла
	//функция опен возвращает ошибку и некую сущность с помощью которой мы будем взаимодействовать с базой
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	//с помощью функции ping проверим удалось ли установить соединение с файлом и корректно ли все работает
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect database: %w", err)
	}
	//если все прошло успешно возвращаем результат storage
	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	//пишем sql запрос который будет сохранять запись в базу данных
	q := `INSERT INTO pages (URL, Username) VALUES(?, ?)`
	//выполняем запрос с помощью сущности db которую мы получили в функции new
	//в разработке использовать контексты считается хорошим тоном (лучше сразу о них позаботиться чем встраивать потом)
	if _, err := s.db.ExecContext(ctx, q, p.URL, p.Username); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}
	return nil
}

func (s *Storage) PickRandom(ctx context.Context, userName string) (p *storage.Page, err error) {
	q := `SELECT URL FROM pages WHERE Username = ? ORDER BY RANDOM() LIMIT 1`
	//в функцию скан передается ссылка на переменную в которуюю мы хотим положить результаты, так как у нас 1 запись то переменная одна
	var url string
	s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}
	return &storage.Page{
		URL:      url,
		Username: userName,
	}, nil
}

func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	q := `DELETE FROM pages WHERE URL = ? AND Username = ?`
	if _, err := s.db.ExecContext(ctx, q, p.URL, p.Username); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}
	return nil
}

func (s *Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE URL = ? AND Username = ?`
	var count int
	if err := s.db.QueryRowContext(ctx, q, p.URL, p.Username).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}
	return count > 0, nil
}

// инициализация нашей базы
func (s *Storage) Init(ctx context.Context) error {
	//эта функция будет создавать таблицу которая будет хранить наши странички
	q := `CREATE TABLE IF NOT EXISTS pages (URL TEXT, Username TEXT)`
	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}
