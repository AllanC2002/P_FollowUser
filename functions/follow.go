package functions

import (
	"errors"

	"gorm.io/gorm"

	"github.com/AllanC2002/P_FollowUser/models"
)

func FollowUser(db *gorm.DB, idFollower int, idFollowing int) (string, int, error) {
	if idFollower == idFollowing {
		return "", 400, errors.New("You cannot follow yourself")
	}

	// Search active profiles
	var followerProfile models.Profile
	var followingProfile models.Profile

	if err := db.Where("Id_User = ? AND Status_account = 1", idFollower).First(&followerProfile).Error; err != nil {
		return "", 404, errors.New("Follower profile not found or inactive")
	}

	if err := db.Where("Id_User = ? AND Status_account = 1", idFollowing).First(&followingProfile).Error; err != nil {
		return "", 404, errors.New("Following profile not found or inactive")
	}

	// Verify before following
	var existing models.Followers
	err := db.Where("Id_Follower = ? AND Id_Following = ?", idFollower, idFollowing).First(&existing).Error
	if err == nil {
		if existing.Status == 1 {
			return "Already following", 200, nil
		} else {
			existing.Status = 1
			db.Save(&existing)
			return "Follow re-activated", 200, nil
		}
	} else if err != gorm.ErrRecordNotFound {
		return "", 500, err
	}

	// Create new follow
	newFollow := models.Followers{
		IdFollower:  idFollower,
		IdFollowing: idFollowing,
		Status:      1,
	}

	if err := db.Create(&newFollow).Error; err != nil {
		return "", 500, err
	}

	return "Followed successfully", 201, nil
}
