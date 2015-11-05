package roll

type AdminRepo interface {
	IsAdmin(subject string) (bool, error)
}
