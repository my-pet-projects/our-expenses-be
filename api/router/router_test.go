package router

// import (
// 	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/testing/mocks"
// 	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/web/api/controllers"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestInitializeRoutes_ReturnsRouter(t *testing.T) {
// 	categoryCtrl := new(mocks.CategoryControllerInterface)
// 	router := &Router{
// 		categoryCtrl: categoryCtrl,
// 	}

// 	result := router.InitializeRoutes()

// 	assert.NotNil(t, result, "Router should not be nil.")
// }

// func TestProvideRouter_ReturnsRouter(t *testing.T) {
// 	controller := new(controllers.CategoryController)

// 	results := ProvideRouter(controller)

// 	assert.NotNil(t, results, "Router should not be nil.")
// }
