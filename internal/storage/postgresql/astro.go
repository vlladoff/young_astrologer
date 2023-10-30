package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/vlladoff/young_astrologer/internal/http-server/handlers/astro"
)

func (s *Storage) GetAllAstroData() ([]astro.AstroData, error) {
	const op = "storage.postgresql.astro.GetAllAstroData"

	query := `SELECT datetime date, title, explanation, url, hd_url FROM data;`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var astroData []astro.AstroData

	for rows.Next() {
		var data astro.AstroData
		if err := rows.Scan(&data.Date, &data.Title, &data.Explanation, &data.Url, &data.HdUrl); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		astroData = append(astroData, data)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return astroData, nil
}

func (s *Storage) GetAstroDataByDay(date string) (astro.AstroData, error) {
	const op = "storage.postgresql.astro.GetByDay"

	if date == "" {
		return astro.AstroData{}, fmt.Errorf("%s: %w", op, errors.New("empty date parameter"))
	}

	query := `SELECT datetime date, title, explanation, url, hd_url FROM data WHERE datetime = $1;`

	var data astro.AstroData

	err := s.db.QueryRow(query, date).Scan(&data.Date, &data.Title, &data.Explanation, &data.Url, &data.HdUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return data, fmt.Errorf("%s: %w", op, errors.New("data not found"))
		}
		return data, fmt.Errorf("%s: %w", op, err)
	}

	return data, nil
}

func (s *Storage) SaveAstroData(data astro.AstroData, imagesId int64) error {
	const op = "storage.postgresql.astro.SaveAstroData"

	stmt, err := s.db.Prepare(
		"INSERT INTO data(datetime, title, explanation, url, hd_url, images_id) VALUES($1, $2, $3, $4, $5, $6)",
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(data.Date, data.Title, data.Explanation, data.Url, data.HdUrl, imagesId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
