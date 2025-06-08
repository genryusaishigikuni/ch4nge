package routes

import "github.com/gin-gonic/gin"
import Auth "github.com/genryusaishigikuni/ch4nge/handlers/auth"
import Middleware "github.com/genryusaishigikuni/ch4nge/middleware"
import User "github.com/genryusaishigikuni/ch4nge/handlers/user"
import Achievement "github.com/genryusaishigikuni/ch4nge/handlers/achievement"
import MiniChallenges "github.com/genryusaishigikuni/ch4nge/handlers/challenges"
import Activity "github.com/genryusaishigikuni/ch4nge/handlers/activity"
import Actions "github.com/genryusaishigikuni/ch4nge/handlers/action"
import Post "github.com/genryusaishigikuni/ch4nge/handlers/post"
import Admin "github.com/genryusaishigikuni/ch4nge/handlers/admin"

func SetupRoutes(r *gin.Engine) {

	// Serve static files (uploaded images)
	r.Static("/uploads", "./uploads")

	// Public routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", Auth.Register)
		auth.POST("/login", Auth.Login)
		auth.POST("/logout", Middleware.AuthMiddleware(), Auth.Logout)
	}

	// Protected routes
	api := r.Group("/", Middleware.AuthMiddleware())
	{
		// User routes
		users := api.Group("/users")
		{
			users.GET("", User.GetAllUsers)
			users.GET("/:userId", User.GetUserDetails)
			users.GET("/:userId/friends", User.GetUserFriends)
			users.PUT("/:userId/friends", User.UpdateUserFriends)
			users.POST("/:userId/profile-pic", User.UploadProfilePicture)

			// NEW: User's own content endpoints
			users.GET("/:userId/actions", User.GetUserActions)           // User's action history
			users.GET("/:userId/actions/stats", User.GetUserActionStats) // Action statistics
			users.GET("/:userId/activities", Activity.GetUserActivities) // User's own activities
			users.GET("/:userId/posts", Post.GetUserPosts)               // User's own posts
			users.GET("/:userId/dashboard", User.GetUserDashboard)       // Overall dashboard

			// Achievement routes
			users.GET("/:userId/achievements", Achievement.GetAllAchievements)
			users.GET("/:userId/achievements/next", Achievement.GetNextAchievement)
			users.GET("/:userId/achievements/progress", Achievement.GetAchievementProgress)

			// Challenge routes
			users.GET("/:userId/mini-challenges", MiniChallenges.GetMiniChallenges)
			users.GET("/:userId/weekly-challenge", MiniChallenges.GetWeeklyChallenge)
			users.GET("/:userId/challenges/completed", MiniChallenges.GetCompletedChallenges) // NEW: Completed challenges
			users.PUT("/:userId/weekly-challenge/:challengeId/progress", MiniChallenges.UpdateUserWeeklyChallengeProgress)
		}

		// Activity routes
		activities := api.Group("/activities")
		{
			activities.POST("/friends", Activity.GetFriendsActivities)
		}

		// Action routes
		actions := api.Group("/actions")
		{
			actions.POST("/green", Actions.UploadGreenAction)
			actions.POST("/transportation", Actions.UploadTransportationAction)
		}

		// Post routes
		posts := api.Group("/posts")
		{
			posts.POST("", Post.UploadPost)
			posts.GET("/recent", Post.GetRecentPosts)
			posts.POST("/:postId/like", Post.LikePost)
			posts.POST("/:postId/share", Post.SharePost)
		}
	}

	// Admin routes
	admin := r.Group("/admin", Middleware.AuthMiddleware(), Middleware.AdminMiddleware())
	{
		// Achievement management
		achievements := admin.Group("/achievements")
		{
			achievements.GET("", Admin.GetAllAchievementsAdmin)
			achievements.POST("", Admin.CreateAchievement)
			achievements.PUT("/:id", Admin.UpdateAchievement)
			achievements.DELETE("/:id", Admin.DeleteAchievement)
			achievements.POST("/assign", Admin.AssignAchievementToUser)
		}

		// Mini challenge management
		miniChallenges := admin.Group("/mini-challenges")
		{
			miniChallenges.GET("", Admin.GetAllMiniChallengesAdmin)
			miniChallenges.POST("", Admin.CreateMiniChallenge)
			miniChallenges.PUT("/:id", Admin.UpdateMiniChallenge)
			miniChallenges.DELETE("/:id", Admin.DeleteMiniChallenge)
			miniChallenges.POST("/assign", Admin.AssignMiniChallengeToUser)
		}

		// Weekly challenge management
		weeklyChallenges := admin.Group("/weekly-challenges")
		{
			weeklyChallenges.GET("", Admin.GetAllWeeklyChallengesAdmin)
			weeklyChallenges.POST("", Admin.CreateWeeklyChallenge)
			weeklyChallenges.PUT("/:id", Admin.UpdateWeeklyChallenge)
			weeklyChallenges.DELETE("/:id", Admin.DeleteWeeklyChallenge)
			weeklyChallenges.POST("/assign", Admin.AssignWeeklyChallengeToUser)
		}
	}
}
