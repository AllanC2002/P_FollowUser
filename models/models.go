package models

type Profile struct {
	IdUser        int    `gorm:"column:Id_User;primaryKey"`
	UserMail      string `gorm:"column:User_mail;unique"`
	Name          string `gorm:"column:Name"`
	Lastname      string `gorm:"column:Lastname"`
	Description   string `gorm:"column:Description"`
	IdPreferences int    `gorm:"column:Id_preferences"`
	IdType        int    `gorm:"column:Id_type"`
	StatusAccount int    `gorm:"column:Status_account"`
}

func (Profile) TableName() string {
	return "Profile"
}

type Followers struct {
	IdFollows   int `gorm:"column:Id_Follows;primaryKey;autoIncrement"`
	IdFollower  int `gorm:"column:Id_Follower"`
	IdFollowing int `gorm:"column:Id_Following"`
	Status      int `gorm:"column:Status"`
}

func (Followers) TableName() string {
	return "Followers"
}
