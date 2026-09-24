package initializers

import "booking-app/model"

func SyncDb() error {
	if err := DB.AutoMigrate(&model.User{}); err != nil {
		return err
	}

	if err := DB.AutoMigrate(&model.Work{}); err != nil {
		return err
	}

	if err := DB.AutoMigrate(&model.Event{}); err != nil {
		return err
	}

	return nil
}
