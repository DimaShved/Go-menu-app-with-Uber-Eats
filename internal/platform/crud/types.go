package crud

type Preloader interface {
	PreloadRelations() []string
}
