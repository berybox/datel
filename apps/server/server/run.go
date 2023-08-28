package server

import (
	"net/http"

	"github.com/berybox/datel/apps/server/fiberutils"
	"github.com/berybox/datel/apps/server/handle"
	"github.com/berybox/datel/apps/server/middleware/dummyheader"
	"github.com/berybox/datel/apps/server/middleware/useroverride"
	"github.com/berybox/datel/apps/server/static"
	"github.com/berybox/datel/pkg/mongodb"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
)

const (
	defaultStaticDir  = "../static"
	templateExtension = ".html"
	debugUsername     = "debug"
)

// Run run Datel server
func Run(mongoURI, addr, overrideid string, isDebugMode bool) error {
	var templateEngine *html.Engine

	if isDebugMode {
		templateEngine = html.New(defaultStaticDir, templateExtension)
		templateEngine.Reload(true)
	} else {
		templateEngine = html.NewFileSystem(http.FS(static.FS), templateExtension)
	}

	_, err := mongodb.GetInstance(mongoURI)
	if err != nil {
		return err
	}

	fiberApp := fiber.New(fiber.Config{
		Views:        templateEngine,
		ErrorHandler: handle.ErrorGET,
	})

	if overrideid != "" {
		fiberApp.Use(useroverride.New(useroverride.Config{Name: overrideid}))
	}

	if isDebugMode {
		fiberApp.Use(dummyheader.New(
			dummyheader.Config{Key: fiberutils.UserIDHeader, Value: debugUsername},
		))
	}

	fiberApp.Get("/", handle.HomeGET).Name("home")

	fiberApp.Get("/add-collection", handle.AddCollectionGET)
	fiberApp.Post("/add-collection", handle.AddCollectionPOST)
	fiberApp.Get("/delete-collection/:collection", handle.DeleteCollectionGET)

	fiberApp.Get("/add-item/:collection", handle.AddItemGET)
	fiberApp.Get("/add-item/:collection/:dbid", handle.AddItemUpdateGET)
	fiberApp.Post("/add-item/:collection", handle.AddItemPOST)
	fiberApp.Post("/add-item/:collection/:dbid", handle.AddItemUpdatePOST)
	fiberApp.Delete("/add-item/:collection/:dbid", handle.AddItemDELETE)

	fiberApp.Get("/add-field/:collection", handle.AddFieldGET)
	fiberApp.Get("/add-field/:collection/:field", handle.AddFieldUpdateGET)
	fiberApp.Post("/add-field/:collection", handle.AddFieldPOST)

	fiberApp.Get("/edit-fields/:collection", handle.EditFieldsGET)
	fiberApp.Post("/edit-fields/:collection", handle.EditFieldsPOST)

	fiberApp.Get("/show-items/:collection", handle.ShowItemsGET)

	if isDebugMode {
		fiberApp.Static("/static", defaultStaticDir, fiber.Static{
			Browse: true,
		})
	} else {
		fiberApp.Use("/static", filesystem.New(filesystem.Config{
			Root:   http.FS(static.FS),
			Browse: false,
		}))
	}

	err = fiberApp.Listen(addr)
	if err != nil {
		return err
	}

	return nil
}
