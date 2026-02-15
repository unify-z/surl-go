package ShortURLService

import (
	"github.com/unify-z/go-surl/internal/errors"
	ShortURL2 "github.com/unify-z/go-surl/internal/repository/surl"
	"github.com/unify-z/go-surl/internal/utils"
)

type ShortURLManager struct {
	repo *ShortURL2.ShortURLRepo
}

func NewShortURLManager(shortURLRepo *ShortURL2.ShortURLRepo) *ShortURLManager {
	return &ShortURLManager{repo: shortURLRepo}
}

func (sm *ShortURLManager) CreateShortURL(originalURL string, userID uint) (*ShortURL2.ShortURL, error) {
	//var shortCode string = utils.Base62Encode(originalURL + utils.RandStr(3))
	var shortCode string = utils.RandStr(4)
	url := &ShortURL2.ShortURL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		UserID:      userID,
	}
	err := sm.repo.CreateShortURL(url)
	if err != nil {
		if utils.IsUniqueConstraintError(err) {
			return sm.CreateShortURL(originalURL, userID)
		}
		return nil, err
	}
	return url, nil
}

func (sm *ShortURLManager) GetShortURLByShortCode(shortCode string) (*ShortURL2.ShortURL, error) {
	url, err := sm.repo.GetShortURLByShortCode(shortCode)
	if err != nil {
		if url == nil {
			return nil, errors.ErrShortURLNotFound
		}
		return nil, err
	}
	return url, nil
}

func (sm *ShortURLManager) GetOriginalURLFromShortCode(shortCode string) (string, error) {
	urlData, err := sm.repo.GetShortURLByShortCode(shortCode)
	if err != nil {
		if urlData == nil {
			return "", errors.ErrShortURLNotFound
		}
		return "", err
	}
	return urlData.OriginalURL, nil
}

func (sm *ShortURLManager) ListShortURLsByUserID(userID uint) ([]ShortURL2.ShortURL, error) {
	urls, err := sm.repo.GetShortURLsByUserID(userID)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func (sm *ShortURLManager) DeleteShortURL(id uint) error {
	return sm.repo.DeleteShortURL(id)
}

func (sm *ShortURLManager) UpdateShortURL(url *ShortURL2.ShortURL) error {
	return sm.repo.UpdateShortURL(url)
}

func (sm *ShortURLManager) ListAllShortURLsWithPagination(page int, pageSize int) ([]ShortURL2.ShortURL, error) {
	return sm.repo.ListShortURLsWithPagination(page, pageSize)
}

func (sm *ShortURLManager) CountShortURLs() (int64, error) {
	return sm.repo.CountSURLs()
}
