package config

import "time"

type Datasource struct {
	Connector      string        `json:"connector"`
	Url            string        `json:"url"`
	Host           string        `json:"host"`
	Port           int16         `json:"port"`
	User           string        `json:"user"`
	Password       string        `json:"password"`
	Database       string        `json:"database"`
	ConnectTimeout time.Duration `json:"connectTimeout"`
}

type Datasources struct {
	Default string                `json:"default"`
	List    map[string]Datasource `json:"list"`
}

type ErrNotFoundDatasource struct {
	name string
}

func (e ErrNotFoundDatasource) Error() string {
	return "datasource not found: " + e.name
}

func (ds Datasources) Get(names ...string) Datasource {
	name := ds.Default
	if len(names) > 0 {
		name = names[0]
	}
	if value, isExist := ds.List[name]; isExist {
		return value
	} else {
		panic(ErrNotFoundDatasource{name: name})
	}
}
