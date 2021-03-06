// based on https://doc.getqor.com/get_started.html
package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"plugin"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/qor/admin"
	"github.com/qor/assetfs"
	"github.com/qor/qor/utils"

	_ "github.com/sergolius/qor_bindatafs_example/config/bindatafs"
)

// Define a GORM-backend model
type User struct {
	gorm.Model
	Name string
}

// Define another GORM-backend model
type Product struct {
	gorm.Model
	Name        string
	Description string
}

var DB *gorm.DB

func main() {
	// Set up the database
	DB, _ = gorm.Open("sqlite3", "demo.db")
	DB.AutoMigrate(&User{}, &Product{})

	// Initialize AssetFS
	AssetFS := assetfs.AssetFS().NameSpace("admin")
	// Register custom paths to manually saved views
	AssetFS.RegisterPath(filepath.Join(utils.AppRoot, "qor/admin/views"))

	// Initialize Admin
	Admin := admin.New(&admin.AdminConfig{
		DB:      DB,
		AssetFS: AssetFS,
	})

	// Create resources from GORM-backend model
	Admin.AddResource(&User{})
	Admin.AddResource(&Product{})

	// Plugins
 
    // The plugins (the *.so files) must be in a 'plugins' sub-directory
    all_plugins, err := filepath.Glob("./plugins/*.so")
    if err != nil {
        panic(err)
    }
 
    for _, filename := range (all_plugins) {
        p, err := plugin.Open(filename)
        if err != nil {
            panic(err)
        }

        symbol, err := p.Lookup("Migrate")
        if err != nil {
            panic(err)
        }
 
        migrateFunc, ok := symbol.(func() []interface{})
        if !ok {
            panic("Plugin has no 'Migrate() []interface{}' function")
        }

        tables := migrateFunc()
        for _, table := range tables {
        	DB.AutoMigrate(table)
        	Admin.AddResource(table)
        }

    }

	// Initialize an HTTP request multiplexer
	mux := http.NewServeMux()

	// Mount admin to the mux
	Admin.MountTo("/admin", mux)

	fmt.Println("Listening on: 8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalln(err)
	}
}

func TruncateTables(tables ...interface{}) {
	for _, table := range tables {
		if err := DB.DropTableIfExists(table).Error; err != nil {
			panic(err)
		}
		DB.AutoMigrate(table)
	}
}
