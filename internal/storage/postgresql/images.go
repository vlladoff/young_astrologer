package postgresql

import (
	"bytes"
	"errors"
	"fmt"
)

func (s *Storage) GetImageData(imageName string) ([]byte, error) {
	const op = "storage.images.GetImageData"

	if imageName == "" {
		return nil, fmt.Errorf("%s: %w", op, errors.New("imageName is empty"))
	}

	query := `SELECT
		CASE
			WHEN $1 = url THEN images.image
			WHEN $1 = hd_url THEN images.hd_image
			END AS image_data
		FROM data
				 JOIN images ON data.images_id = images.id
		WHERE url = $1 OR hd_url = $1;`

	var image []byte

	err := s.db.QueryRow(query, imageName).Scan(&image)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return image, nil
}

func (s *Storage) SaveImages(image, hdImage *bytes.Buffer) (int64, error) {
	const op = "storage.postgresql.images.SaveImages"

	query := "INSERT INTO images(image, hd_image) VALUES($1, $2) RETURNING id"

	var id int64

	err := s.db.QueryRow(query, image.Bytes(), hdImage.Bytes()).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
