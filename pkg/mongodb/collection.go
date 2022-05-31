package mongodb

type Collection interface {
	ColName() string
}

type DynCollection interface {
	DynColName(resourceId string) string
}
