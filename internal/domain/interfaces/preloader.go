package interfaces

type Preloader interface {
	PreloadRelations() []string
}
