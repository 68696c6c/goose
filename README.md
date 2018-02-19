# Goose

Goose is a basic Golang migration library for MySQL built using GORM and Logrus.

## Current Version
0.1.0

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

Goose uses GORM for database interactions.  Your project will need to provide a
GORM database connection to Goose.  Future developments will remove this 
dependency in favor of using the built-in sql driver.

Goose uses Logrus for logging.  You can pass a logrus.Logger instance to the 
schema constructor, or pass nil to have Goose create it's own log at the 
basePath provded in the second argument to the constructor.

```
import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
  db, _ := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
  defer db.Close()
  
  // Goose will create a log file at /path/to/your/project/goose.log
  schema := goose.NewSchema(db, "/path/to/your/project", nil)
}
```

### Installing

Using dep:

```
dep ensure -add github.com/68696c6c/goose
```

If you're using `go get` for dependency management, add an import for 
"github.com/68696c6c/goose" and run:

```
go get ./...
```

## Running the tests

Tests are coming soon...

## Built With

* [GORM](github.com/jinzhu/gorm) - database library
* [Logrus](github.com/Sirupsen/logrus) - logging library

## Contributing

Goose is still in the early stages of development.  If you have a feature or bug
fix, feel free to make a pull request.

## Authors

* **Aaron Hill** [68696c6c](https://github.com/68696c6c)
* **Wes Curtis** 

## License

This project is licensed under the MIT License
