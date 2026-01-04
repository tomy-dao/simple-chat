package repo

import (
	"local/config"
	"local/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Repository struct {
	db              *gorm.DB
	UserRepo        UserRepo
	ConversationRepo ConversationRepo
	ParticipantRepo  ParticipantRepo
	MessageRepo      MessageRepo
}

// NewRepositoryWithDB creates a repository instance with the provided database
// This can be used for both production and testing
func NewRepositoryWithDB(db *gorm.DB) (*Repository, error) {
	// Auto migrate all models
	err := db.AutoMigrate(
		&model.User{},
		&model.Conversation{},
		&model.ConversationParticipant{},
		&model.Message{},
	)
	if err != nil {
		return nil, err
	}

	// Create repositories
	userRepo := &userRepository{db: db}
	conversationRepo := &conversationRepository{db: db}
	participantRepo := &participantRepository{db: db}
	messageRepo := &messageRepository{db: db}

	return &Repository{
		db:              db,
		UserRepo:        userRepo,
		ConversationRepo: conversationRepo,
		ParticipantRepo:  participantRepo,
		MessageRepo:      messageRepo,
	}, nil
}

// NewRepository creates a repository instance with MySQL connection from config
// This is used for production
func NewRepository() (*Repository, error) {
	// Create database connection
	dsn := config.Config.DBUser + ":" + config.Config.DBPassword + "@tcp(" + config.Config.DBHost + ":" + config.Config.DBPort + ")/" + config.Config.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return NewRepositoryWithDB(db)
}
