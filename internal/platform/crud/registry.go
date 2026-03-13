package crud

import "github.com/gofiber/fiber/v3"

type RouteRegistrar interface {
	RegisterRoutes(app *fiber.App)
}

type Registry struct {
	resources []RouteRegistrar
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (r *Registry) Register(resource RouteRegistrar) {
	r.resources = append(r.resources, resource)
}

func (r *Registry) Mount(app *fiber.App) {
	for _, resource := range r.resources {
		resource.RegisterRoutes(app)
	}
}
