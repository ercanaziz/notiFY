package alarm

import "time"

type AlarmService struct {
	Repo *AlarmRepository
}

func NewAlarmService(repo *AlarmRepository) *AlarmService {
	return &AlarmService{Repo: repo}
}

// Yeni alarm oluşturma mantığı
func (s *AlarmService) Create(userID string, input AlertInput) (Alert, error) {
	newAlert := Alert{
		UserID:      userID,
		ProductID:   input.ProductID,
		TargetPrice: input.TargetPrice,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}

	// Veritabanına kaydetmesi için repo'ya gönderiyoruz
	return s.Repo.CreateAlert(newAlert)
}

// Aktif alarmları getirme mantığı
func (s *AlarmService) GetActiveAlerts(userID string) ([]Alert, error) {
	return s.Repo.GetUserAlerts(userID)
}

func (s *AlarmService) DeleteAlert(id string) error {
	// Buradaki 's.repo' (veya senin verdiğin isim) Repository'ye gider
	return s.Repo.DeleteAlert(id)
}
