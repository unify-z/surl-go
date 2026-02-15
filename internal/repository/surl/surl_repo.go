package surl

import "gorm.io/gorm"

type ShortURLRepo struct {
	db *gorm.DB
}

func NewShortURLRepo(db *gorm.DB) *ShortURLRepo {
	return &ShortURLRepo{db: db}
}

func (s *ShortURLRepo) CountSURLs() (int64, error) {
	var count int64
	err := s.db.Model(&ShortURL{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *ShortURLRepo) CreateShortURL(url *ShortURL) error {
	return s.db.Create(url).Error
}

func (s *ShortURLRepo) GetShortURLByShortCode(shortCode string) (*ShortURL, error) {
	var url ShortURL
	err := s.db.Where("short_code = ?", shortCode).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (s *ShortURLRepo) GetShortURLsByUserID(userID uint) ([]ShortURL, error) {
	var urls []ShortURL
	err := s.db.Where("user_id = ?", userID).Find(&urls).Error
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func (s *ShortURLRepo) DeleteShortURL(id uint) error {
	return s.db.Delete(&ShortURL{}, id).Error
}

func (s *ShortURLRepo) UpdateShortURL(url *ShortURL) error {
	return s.db.Save(url).Error
}

func (s *ShortURLRepo) ListAllShortURLs() ([]ShortURL, error) {
	var urls []ShortURL
	err := s.db.Find(&urls).Error
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func (s *ShortURLRepo) ListShortURLsWithPagination(page int, pageSize int) ([]ShortURL, error) {
	var urls []ShortURL
	offset := (page - 1) * pageSize

	err := s.db.
		Order("id ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&urls).Error

	if err != nil {
		return nil, err
	}
	return urls, nil
}
