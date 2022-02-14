package repository

type Repositories struct {
	Category Category
}

func NewRepositories(cateogry Category) *Repositories {
	return &Repositories{
		Category: cateogry,
	}
}
