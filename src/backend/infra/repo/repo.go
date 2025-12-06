package repo

import (
	"local/config"
	"local/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	User() UserRepo
	Conversation() ConversationRepo
	Participant() ParticipantRepo
	Message() MessageRepo
}

type Repository struct {
	db           *gorm.DB
	user         UserRepo
	conversation ConversationRepo
	participant  ParticipantRepo
	message      MessageRepo
}

func (r *Repository) User() UserRepo {
	return r.user
}

func (r *Repository) Conversation() ConversationRepo {
	return r.conversation
}

func (r *Repository) Participant() ParticipantRepo {
	return r.participant
}

func (r *Repository) Message() MessageRepo {
	return r.message
}

// NewRepositoryWithDB creates a repository instance with the provided database
// This can be used for both production and testing
func NewRepositoryWithDB(db *gorm.DB) (RepositoryInterface, error) {
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
		db:           db,
		user:         userRepo,
		conversation: conversationRepo,
		participant:  participantRepo,
		message:      messageRepo,
	}, nil
}

// NewRepository creates a repository instance with MySQL connection from config
// This is used for production
func NewRepository() (RepositoryInterface, error) {
	// Create database connection
	dsn := config.Config.DBUser + ":" + config.Config.DBPassword + "@tcp(" + config.Config.DBHost + ":" + config.Config.DBPort + ")/" + config.Config.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return NewRepositoryWithDB(db)
}

