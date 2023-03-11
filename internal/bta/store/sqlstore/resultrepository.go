package sqlstore

import (
	"BrainTApp/internal/bta/store"
	entity "BrainTApp/internal/model/entity"
	"database/sql"
	"strings"
	"time"
)

type ResultRepository struct {
	store *Store
}

func (r *ResultRepository) CreateArithmetic(ar *entity.Result) error {
	if err := ar.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO res_test (user_id, test_id, date_test, result_inter, result_1, result_2) VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id",
		ar.UserID,
		ar.TestID,
		ar.DateTest,
		ar.ResultInter,
		ar.ResultOne,
		ar.ResultTwo,
	).Scan(&ar.UserID)
}

func (r *ResultRepository) CreateMemorize(mem *entity.Result) error {
	if err := mem.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO res_test (user_id, test_id, date_test, result_inter, result_1, result_2, result_3, test_set) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING user_id",
		mem.UserID,
		mem.TestID,
		mem.DateTest,
		mem.ResultInter,
		mem.ResultOne,
		mem.ResultTwo,
		mem.Results,
		mem.TestSet,
	).Scan(&mem.UserID)
}

func (r *ResultRepository) Find(user_id int, test_id int) (*entity.Result, error) {
	res := &entity.Result{}
	if err := r.store.db.QueryRow(
		"SELECT * FROM res_test WHERE user_id = $1 and test_id = $2",
		user_id, test_id).Scan(
		&res.UserID,
		&res.TestID,
		&res.DateTest,
		&res.ResultInter,
		&res.ResultOne,
		&res.ResultTwo,
		&res.Results,
		&res.TestSet,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return res, nil
}

func (r *ResultRepository) FindAll(user_id int) ([]entity.Result, error) {

	sqls := "SELECT user_id, test_id, date_test, result_inter,\n" +
		"CASE WHEN result_1 IS NULL THEN 0 ELSE result_1 END, \n" +
		"CASE WHEN result_2 IS NULL THEN 0 ELSE result_2 END, \n" +
		"CASE WHEN result_3 IS NULL THEN '' ELSE result_3 END, \n" +
		"CASE WHEN test_set IS NULL THEN '' ELSE test_set  END \n" +
		"FROM res_test WHERE user_id = $1 ORDER BY test_id , date_test"

	rows, err := r.store.db.Query(sqls, user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	defer rows.Close()

	results := make([]entity.Result, 0)
	for rows.Next() {
		res := &entity.Result{}
		resultInterStr := ""

		if err := rows.Scan(
			&res.UserID,
			&res.TestID,
			&res.DateTest,
			&resultInterStr,
			&res.ResultOne,
			&res.ResultTwo,
			&res.Results,
			&res.TestSet,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		parts := strings.Split(resultInterStr, ":")
		durationStr := parts[0] + "h" + parts[1] + "m" + parts[2] + "s"
		duration, _ := time.ParseDuration(durationStr)
		res.ResultInter = time.Duration(duration.Seconds())
		results = append(results, *res)
	}

	return results, nil
}

func (r *ResultRepository) CreateShulte(sh *entity.Result) error {

	if err := sh.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO res_test (user_id, test_id, result_inter, result_3, date_test) VALUES ($1, $2, $3, $4, $5) RETURNING user_id",
		sh.UserID,
		sh.TestID,
		sh.ResultInter,
		sh.Results,
		sh.DateTest,
	).Scan(&sh.UserID)

}

func (r *ResultRepository) TestSetWords(num int, partSpeech int) ([]string, error) {
	testSet := []string{}
	rows, err := r.store.db.Query(
		"SELECT word FROM words WHERE part_ps_id = $1 ORDER BY random() LIMIT $2", partSpeech, num)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var word string
		if err := rows.Scan(&word); err != nil {
			return nil, err
		}
		testSet = append(testSet, word)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return testSet, nil
}
