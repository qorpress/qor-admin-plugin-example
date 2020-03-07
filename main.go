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

func main() {
	// Plugins
	   numbers := []int{5, 2, 7, 6, 1, 3, 4, 8}
 
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
         
        symbol, err := p.Lookup("Sort")
        if err != nil {
            panic(err)
        }
 
        sortFunc, ok := symbol.(func([]int) *[]int)
        if !ok {
            panic("Plugin has no 'Sort([]int) []int' function")
        }
 
        sorted := sortFunc(numbers)
        fmt.Println(filename, sorted)
    }


	// Set up the database
	DB, _ := gorm.Open("sqlite3", "demo.db")
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
