package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"TelegramBotGolang/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// создаём хранилище
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open databases: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}
	return &Storage{db: db}, nil
}

// далее объявляем методы структуры Storage (которые прописаны в интерфейсе Storage)

func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)` // на месте знаков вопроса будут значения, которые подставляются в функцию ExecContext

	// функция выпоняет запрос
	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	// сортируем статьи пользователя по рандому и берём первый элемент
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)

	// ошибка ниже говорит о том что пользователь ничего не сохранял, и мы ничего не возвращаем
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}

	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}

	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`
	if _, err := s.db.ExecContext(ctx, q, page.URL, page.UserName); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}
	return nil
}

func (s *Storage) IsExists(ctx context.Context, page *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return count > 0, nil
}

// создние таблицы
func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages(url TEXT, user_name TEXT)`
	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}
