package controller

import (
	"database/sql"

	"job/domain/repository"
	"job/presentation/core/config"
	"job/presentation/core/logger"
)

type Environment struct {
	Balances *sql.DB
	logger   interface {
		Info(message string, source string)
		Error(message string, source string)
		Warning(message string, source string)
		Fatal(message string, source string)
	}
}

func (env *Environment) SetUsersDatabase(db *sql.DB) *Environment {
	env.Balances = db
	return env
}

func (env *Environment) SetLogger(logger *logger.Logger) *Environment {
	env.logger = logger
	return env
}

func SetMaxConnections(db *sql.DB, max int) {
	db.SetMaxOpenConns(max)
}

func NewEnvironment() (*Environment, error) {
	env := new(Environment)
	conf := config.Get()
	users, err := repository.NewPgDatabase(conf.Database.User, conf.Database.Password, conf.Database.Host, conf.Database.Name, conf.Database.Port)
	if err != nil {
		return nil, err
	}
	logger, err := logger.NewLogger().
		SetApp(conf.Application.Name).
		SetVersion(conf.Application.Version).
		CreateLogger()
	if err != nil {
		return nil, err
	}

	SetMaxConnections(users, 10)
	env.SetLogger(logger)
	env.SetUsersDatabase(users)
	return env, nil
}
