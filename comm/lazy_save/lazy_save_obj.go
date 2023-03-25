package lazy_save

type lazySaveObj interface {
	GetLsoId() string
	SaveOrUpdate(callback func())
}
