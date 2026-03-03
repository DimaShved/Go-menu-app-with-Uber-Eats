package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"uber-go-menu/internal/repository"
	"uber-go-menu/internal/routes/rest"
	"uber-go-menu/internal/service"
)

func NewHTTPServer(database *gorm.DB, vld *validator.Validate) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	app.Use(requestLogger())

	restaurantRepo := repository.NewRestaurantRepo(database)
	restaurantService := service.NewRestaurantService(restaurantRepo)

	menuSectionRepo := repository.NewMenuSectionRepo(database)
	menuSectionService := service.NewMenuSectionService(menuSectionRepo)

	menuCategoryRepo := repository.NewMenuCategoryRepo(database)
	menuCategoryService := service.NewMenuCategoryService(menuCategoryRepo)

	menuItemRepo := repository.NewMenuItemRepo(database)
	menuItemService := service.NewMenuItemService(menuItemRepo)

	menuAvailabilityRepo := repository.NewMenuAvailabilityRepo(database)
	menuAvailabilityService := service.NewMenuAvailabilityService(menuAvailabilityRepo)

	variationOptionRepo := repository.NewVariationOptionRepo(database)
	variationOptionService := service.NewVariationOptionService(variationOptionRepo)

	variationRepo := repository.NewVariationRepo(database)
	variationService := service.NewVariationService(variationRepo, variationOptionService)

	modifierOptionRepo := repository.NewModifierOptionRepo(database)
	modifierOptionService := service.NewModifierOptionService(modifierOptionRepo)

	modifierRepo := repository.NewModifierRepo(database)
	modifierService := service.NewModifierService(modifierRepo, modifierOptionService)

	rest.SetupRestaurantRoutes(app, restaurantService, vld)
	rest.SetupMenuSectionRoutes(app, menuSectionService, vld)
	rest.SetupMenuCategoryRoutes(app, menuCategoryService, vld)
	rest.SetupMenuItemRoutes(app, menuItemService, vld)
	rest.SetupMenuAvailabilityRoutes(app, menuAvailabilityService, vld)
	rest.SetupVariationRoutes(app, variationService, vld)
	rest.SetupVariationOptionRoutes(app, variationOptionService, vld)
	rest.SetupModifierOptionRoutes(app, modifierOptionService, vld)
	rest.SetupModifierRoutes(app, modifierService, vld)

	return app
}
