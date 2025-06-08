package seeder

import (
	"time"

	db "github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/models"
)

// SeedBasicAchievements создает базовые достижения
func SeedBasicAchievements() error {
	achievements := []models.Achievement{
		{
			Title:    "First Green Action",
			Subtitle: "Complete your first eco-friendly action",
			Points:   25,
		},
		{
			Title:    "Eco Traveler",
			Subtitle: "Use eco-friendly transportation",
			Points:   20,
		},
		{
			Title:    "Point Collector",
			Subtitle: "Collect 100 eco points",
			Points:   50,
		},
		{
			Title:    "Distance Master",
			Subtitle: "Travel 50km using eco transport",
			Points:   30,
		},
		{
			Title:    "Green Warrior",
			Subtitle: "Complete 10 green actions",
			Points:   75,
		},
		{
			Title:    "Tree Planter",
			Subtitle: "Plant your first tree",
			Points:   100,
		},
		{
			Title:    "Energy Saver",
			Subtitle: "Save energy 5 times",
			Points:   40,
		},
		{
			Title:    "Recycling Hero",
			Subtitle: "Recycle items 5 times",
			Points:   35,
		},
	}

	for _, achievement := range achievements {
		var existing models.Achievement
		if db.DB.Where("title = ?", achievement.Title).First(&existing).Error != nil {
			// Создаем достижение
			if err := db.DB.Create(&achievement).Error; err != nil {
				return err
			}

			// Назначаем всем существующим пользователям
			if err := assignAchievementToAllUsers(achievement.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

// SeedBasicMiniChallenges создает базовые мини-челленджи
func SeedBasicMiniChallenges() error {
	miniChallenges := []models.MiniChallenge{
		{
			Title:    "Daily Walker",
			Subtitle: "Walk at least 2km today",
			Points:   15,
		},
		{
			Title:    "Lights Out",
			Subtitle: "Turn off unnecessary lights",
			Points:   10,
		},
		{
			Title:    "Public Transport",
			Subtitle: "Use public transport today",
			Points:   20,
		},
		{
			Title:    "Recycling Day",
			Subtitle: "Recycle something today",
			Points:   15,
		},
		{
			Title:    "Water Saver",
			Subtitle: "Practice water conservation",
			Points:   12,
		},
	}

	for _, challenge := range miniChallenges {
		var existing models.MiniChallenge
		if db.DB.Where("title = ?", challenge.Title).First(&existing).Error != nil {
			// Создаем мини-челлендж
			if err := db.DB.Create(&challenge).Error; err != nil {
				return err
			}

			// Назначаем всем существующим пользователям
			if err := assignMiniChallengeToAllUsers(challenge.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

// SeedBasicWeeklyChallenges создает базовые недельные челленджи
func SeedBasicWeeklyChallenges() error {
	startOfWeek := getStartOfCurrentWeek()
	endOfWeek := getEndOfCurrentWeek()

	weeklyChallenges := []models.WeeklyChallenge{
		{
			Title:       "Eco Transport Week",
			Subtitle:    "Travel 100km using eco-friendly transport",
			Points:      100,
			TargetValue: 100.0,
			IsActive:    true,
			StartDate:   startOfWeek,
			EndDate:     endOfWeek,
		},
		{
			Title:       "Green Actions Week",
			Subtitle:    "Complete 15 green actions this week",
			Points:      75,
			TargetValue: 15.0,
			IsActive:    true,
			StartDate:   startOfWeek,
			EndDate:     endOfWeek,
		},
		{
			Title:       "Point Challenge",
			Subtitle:    "Earn 200 points this week",
			Points:      150,
			TargetValue: 200.0,
			IsActive:    true,
			StartDate:   startOfWeek,
			EndDate:     endOfWeek,
		},
		{
			Title:       "Carbon Reduction",
			Subtitle:    "Reduce your carbon footprint by 50 units",
			Points:      120,
			TargetValue: 50.0,
			IsActive:    true,
			StartDate:   startOfWeek,
			EndDate:     endOfWeek,
		},
	}

	for _, challenge := range weeklyChallenges {
		var existing models.WeeklyChallenge
		if db.DB.Where("title = ? AND start_date = ?", challenge.Title, startOfWeek).First(&existing).Error != nil {
			// Создаем недельный челлендж
			if err := db.DB.Create(&challenge).Error; err != nil {
				return err
			}

			// Назначаем всем существующим пользователям
			if err := assignWeeklyChallengeToAllUsers(challenge.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

// Вспомогательные функции
func assignAchievementToAllUsers(achievementID uint) error {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		userAchievement := models.UserAchievement{
			UserID:        user.ID,
			AchievementID: achievementID,
		}
		db.DB.FirstOrCreate(&userAchievement, models.UserAchievement{
			UserID:        user.ID,
			AchievementID: achievementID,
		})
	}
	return nil
}

func assignMiniChallengeToAllUsers(challengeID uint) error {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		userChallenge := models.UserMiniChallenge{
			UserID:          user.ID,
			MiniChallengeID: challengeID,
			AssignedAt:      time.Now(),
		}
		db.DB.FirstOrCreate(&userChallenge, models.UserMiniChallenge{
			UserID:          user.ID,
			MiniChallengeID: challengeID,
		})
	}
	return nil
}

func assignWeeklyChallengeToAllUsers(challengeID uint) error {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		userChallenge := models.UserWeeklyChallenge{
			UserID:            user.ID,
			WeeklyChallengeID: challengeID,
			CurrentValue:      0.0,
			AssignedAt:        time.Now(),
		}
		db.DB.FirstOrCreate(&userChallenge, models.UserWeeklyChallenge{
			UserID:            user.ID,
			WeeklyChallengeID: challengeID,
		})
	}
	return nil
}

func getStartOfCurrentWeek() time.Time {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	return now.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
}

func getEndOfCurrentWeek() time.Time {
	start := getStartOfCurrentWeek()
	return start.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
}

// InitializeBasicData инициализирует базовые данные
func InitializeBasicData() error {
	if err := SeedBasicAchievements(); err != nil {
		return err
	}

	if err := SeedBasicMiniChallenges(); err != nil {
		return err
	}

	if err := SeedBasicWeeklyChallenges(); err != nil {
		return err
	}

	return nil
}
