package goose

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

var (
	productionOperations    []string
	productionOperationsSet bool
	isProduction            bool
	envSet                  bool
)

const (
	MigrateOperationInstall = "install"
	MigrateOperationUp      = "up"
	MigrateOperationDrop    = "drop"
	MigrateOperationReset   = "reset"
)

// Set which migrate operations are allowed in production.
// Use this to avoid accidentally running destructive migration operations
// against production data.
func SetProductionOperations(allowed []string) {
	productionOperations = allowed
	productionOperationsSet = true
}

// Get which migrate operations are allowed in production.
func getProductionOperations() []string {
	if !productionOperationsSet {
		productionOperations = []string{
			MigrateOperationUp,
			MigrateOperationInstall,
		}
	}
	return productionOperations
}

// Set whether the current environment is production.
func SetEnvProduction(p bool) {
	isProduction = p
	envSet = true
}

func operationAllowed(operation string) bool {
	// If we have been explicitly told we aren't in production, all commands are allowed.
	if envSet && !isProduction {
		return true
	}
	o := getProductionOperations()
	return stringInSlice(operation, o)
}

func HandleMigrate(schema SchemaInterface, args []string, dryRun *bool) (migrated []string, dropped []string, err error) {

	// Make sure we know what environment we are running in before doing anything.
	if !envSet {
		return migrated, dropped, errors.New("you must set the current environment by calling goose.SetEnvProduction before running any migration operations")
	}

	// If allowed production operations haven't been explicitly set, show a warning.
	if !productionOperationsSet {
		println("production-safe migration operations have not been set; only 'up' and 'install' will be allowed if the current environment has been set to production")
	}

	if *dryRun {
		schema.SetDebug(true)
	}

	var operation string
	if len(args) < 1 {
		operation = MigrateOperationUp
	} else {
		operation = args[0]
	}

	switch operation {
	case MigrateOperationReset:
		allowed := operationAllowed(operation)
		if !allowed {
			msg := fmt.Sprintf("'%s' operation is not allowed in a production environment", operation)
			return migrated, dropped, errors.New(msg)
		}

		dropped, migrated, err = schema.Reset()
		if err != nil {
			return
		}

		dmsg := strings.Join(dropped, "\n")
		println("dropped tables: \n" + dmsg)

		mmsg := strings.Join(migrated, "\n")
		println("migrated tables: \n" + mmsg)

	case MigrateOperationInstall:
		allowed := operationAllowed(operation)
		if !allowed {
			msg := fmt.Sprintf("'%s' operation is not allowed in a production environment", operation)
			return migrated, dropped, errors.New(msg)
		}

		if err = schema.Install(); err != nil {
			return
		}

		println("created migration table")

	case MigrateOperationDrop:
		allowed := operationAllowed(operation)
		if !allowed {
			msg := fmt.Sprintf("'%s' operation is not allowed in a production environment", operation)
			return migrated, dropped, errors.New(msg)
		}

		dropped, err = schema.Drop()
		if err != nil {
			return
		}

		msg := strings.Join(dropped, "\n")
		println("dropped tables: \n" + msg)

	default:
		allowed := operationAllowed(operation)
		if !allowed {
			msg := fmt.Sprintf("'%s' operation is not allowed in a production environment", operation)
			return migrated, dropped, errors.New(msg)
		}

		migrated, err = schema.Migrate()
		if err != nil {
			return
		}

		msg := strings.Join(migrated, "\n")
		println("migrated tables: \n" + msg)

	}

	return
}
