package store

type Store interface {
	GetUserRepo() UserRepository
	GetEventRepo() EventRepository
}
