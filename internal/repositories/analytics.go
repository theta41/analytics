package repositories

type Analytics interface {
	GetCountAcceptedTask() (count int, err error)
	GetCountDeclinedTask() (count int, err error)
	GetSumReaction(objectId uint32) (count int, err error)
	CreateTask(objectId uint32) (taskId int, err error)
	FinishTask(objectId uint32) (err error)
	CreateLetter(objectId uint32, email string) (letterId int, err error)
	AcceptedLetter(objectId uint32, email string) error
	DeclinedLetter(objectId uint32, email string) error
}
