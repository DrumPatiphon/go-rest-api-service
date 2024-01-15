package appinfoUsecases

import (
	"github.com/DrumPatiphon/go-rest-api-service/modules/appInfo"
	appinfoRepositories "github.com/DrumPatiphon/go-rest-api-service/modules/appInfo/appInfoRepositories"
)

type IAppInfoUsecase interface {
	FindCategory(req *appInfo.CategoryFilter) ([]*appInfo.Category, error)
	InsertCagetory(req []*appInfo.Category) error
	DeleteCategory(categoryId int) error
}

type appinfoUsecase struct {
	appInfoRepository appinfoRepositories.IAppInfoRepository
}

func AppInfoUsecase(appInfoRepository appinfoRepositories.IAppInfoRepository) IAppInfoUsecase {
	return &appinfoUsecase{
		appInfoRepository: appInfoRepository,
	}
}

func (u *appinfoUsecase) FindCategory(req *appInfo.CategoryFilter) ([]*appInfo.Category, error) {
	category, err := u.appInfoRepository.FindCategory(req)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (u *appinfoUsecase) InsertCagetory(req []*appInfo.Category) error {
	if err := u.appInfoRepository.InsertCategory(req); err != nil {
		return err
	}
	return nil
}

func (u *appinfoUsecase) DeleteCategory(categoryId int) error {
	if err := u.appInfoRepository.DeleteCategory(categoryId); err != nil {
		return err
	}
	return nil
}
