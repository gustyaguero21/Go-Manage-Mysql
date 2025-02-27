package config

//db queries

const (
	ExistsDB = "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?"
	CreateDB = "CREATE DATABASE IF NOT EXISTS %s"
)

// db test queries

const (
	ExistsTestQuery    = "SELECT \\* FROM `users`"
	SearchTestQuery    = "SELECT \\* FROM `users`"
	SaveTestQuery      = "INSERT INTO `users`"
	UpdateTestQuery    = "UPDATE `users` SET"
	DeleteTestQuery    = "DELETE FROM `users`"
	ChangePwdTestQuery = "UPDATE `users` SET"
)
