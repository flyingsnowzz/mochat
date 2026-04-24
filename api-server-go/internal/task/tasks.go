package task

import (
	"time"

	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/logger"
)

func RegisterCorpDataTask(db *gorm.DB) {
	_ = AddFunc("0 0 2 * * *", func() {
		logger.Sugar.Info("running corp data task...")
		var corps []model.Corp
		if err := db.Find(&corps).Error; err != nil {
			logger.Sugar.Errorf("fetch corps failed: %v", err)
			return
		}

		today := time.Now().Format("2006-01-02")
		for _, corp := range corps {
			dayData := model.CorpDayData{
				CorpID: corp.ID,
				Date:   today,
			}
			db.Where("corp_id = ? AND date = ?", corp.ID, today).FirstOrCreate(&dayData)
		}
	})
}

func RegisterMediaIdUpdateTask() {
	_ = AddFunc("0 0 3 * * *", func() {
		logger.Sugar.Info("running media id update task...")
	})
}

func RegisterEmployeeStatisticTask() {
	_ = AddFunc("0 30 2 * * *", func() {
		logger.Sugar.Info("running employee statistic task...")
	})
}

func RegisterSyncWorkAgentTask() {
	_ = AddFunc("0 0 4 * * *", func() {
		logger.Sugar.Info("running sync work agent task...")
	})
}

func RegisterAllTasks(db *gorm.DB) {
	RegisterCorpDataTask(db)
	RegisterMediaIdUpdateTask()
	RegisterEmployeeStatisticTask()
	RegisterSyncWorkAgentTask()
}
