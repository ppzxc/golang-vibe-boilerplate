//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/ppzxc/golang-vibe-boilerplate/internal/adapter/eventbus"
	"github.com/ppzxc/golang-vibe-boilerplate/internal/adapter/httphandler"
	"github.com/ppzxc/golang-vibe-boilerplate/internal/adapter/postgresrepo"
	"github.com/ppzxc/golang-vibe-boilerplate/internal/app"
	"github.com/ppzxc/golang-vibe-boilerplate/internal/app/todo"
	"github.com/ppzxc/golang-vibe-boilerplate/internal/config"
	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
)

func InitializeServer(cfg *config.Config) (*Server, func(), error) {
	wire.Build(
		config.GetDSN,
		postgresrepo.Open,
		postgresrepo.NewTodoRepository,
		wire.Bind(new(domain.Repository), new(*postgresrepo.TodoRepository)),
		eventbus.NewInMemoryEventBus,
		wire.Bind(new(app.EventBus), new(*eventbus.InMemoryEventBus)),
		todo.NewService,
		wire.Bind(new(httphandler.TodoService), new(*todo.Service)),
		httphandler.NewRouter,
		wire.FieldsOf(new(*config.Config), "Server"),
		wire.FieldsOf(new(config.ServerConfig), "Port"),
		NewServer,
	)
	return nil, nil, nil
}
