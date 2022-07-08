package providers

type DbConfig struct {
	Host      string
	Port      string
	User      string
	Password  string
	DbName    string
	DbTimeOut int
}

func (db DbConfig) ConfigureDbConnection() DbConfig {
	db.Host = "localhost"
	db.Port = "5432"
	db.User = "postgres"
	db.Password = "admin123"
	db.DbName = "local_test_db"
	db.DbTimeOut = 5
	return db
}
