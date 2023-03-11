package sqlstore

import (
	"BrainTApp/internal/bta/store"
	"BrainTApp/internal/model/entity"
	"database/sql"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *entity.User) error {

	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id_user",
		u.Username,
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID)
}

func (r *UserRepository) Find(id int) (*entity.User, error) {
	u := &entity.User{}
	if err := r.store.db.QueryRow(
		"SELECT * FROM users WHERE  id_user = $1",
		id,
	).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	u := &entity.User{}
	if err := r.store.db.QueryRow(
		"SELECT * FROM users WHERE  email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}
