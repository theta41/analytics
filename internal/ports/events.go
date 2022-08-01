package ports

type Events interface {
	CreateTask(objectId uint32) error
	DeleteTask(objectId uint32) error
	FinishTask(objectId uint32) error
	CreateLetter(objectId uint32, email string) error
	AcceptedLetter(objectId uint32, email string) error
	DeclinedLetter(objectId uint32, email string) error
}
