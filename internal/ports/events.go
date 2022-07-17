package ports

type Events interface {
	CreateTask(objectId uint32) (taskId int, err error)
	FinishTask(objectId uint32) (err error)
	CreateLetter(objectId uint32, email string) (letterId int, err error)
	AcceptedLetter(objectId uint32, email string) (err error)
	DeclinedLetter(objectId uint32, email string) (err error)
}
