package postgresql

import "github.com/vlladoff/young_astrologer/internal/http-server/handlers/astro"

func (s *Storage) GetAllAstroData() ([]astro.AstroData, error) {
	const op = "storage.postgresql.GetAll"

	return []astro.AstroData{}, nil
}

func (s *Storage) GetAstroDataByDay(date string) (astro.AstroData, error) {
	const op = "storage.postgresql.GetByDay"

	return astro.AstroData{}, nil
}

func (s *Storage) SaveAstroData(data astro.AstroData, image, hdImage []byte) error {
	const op = "storage.postgresql.SaveAstroData"

	//todo if need return last id
	return nil
}
