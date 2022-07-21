package ports

type Statistics interface {
	GetCountAcceptedTask() (count int, err error)
	GetCountDeclinedTask() (count int, err error)
	GetSumReaction(objectId uint32) (count int, err error)
}
