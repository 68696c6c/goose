package goose

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var (
	installFileName = "install"
	dropFileName    = "drop"
	migrationsDir   = "migrations"
)

type SchemaInterface interface {
	SetDebug(b bool)
	Reset() ([]string, []string, error)
	Install() error
	Drop() ([]string, error)
	Migrate() ([]string, error)
	CreateMigration(name string) (string, error)
	DropTable(table string) []error
	RunMigrationFile(fileName string) []error
	SetForeignKeyChecks(enabled bool) []error
}

// Migration table record.
type migration struct {
	ID       uint64
	FileName string
}

// Schema interaction struct.
type Schema struct {
	migrationsPath string
	installFile    string
	dropFile       string
	log            *logrus.Logger
	gorm           *gorm.DB
	debug          bool
}

func NewSchema(gorm *gorm.DB, basePath string, log *logrus.Logger) (SchemaInterface, error) {
	basePath = strings.TrimRight(basePath, "/")
	if basePath == "" {
		return nil, errors.New("you must provide a base path")
	}
	if !writableDir(basePath) {
		return nil, errors.New("base path does not exist or is not writable")
	}
	migrationPath := basePath + "/" + migrationsDir
	err := createDirIfNotExist(migrationPath)
	if err != nil {
		return nil, errors.New("failed to create migrations directory at " + migrationPath + " Try creating the directory manually.")
	}
	if log == nil {
		log = getLogger(basePath)
	}
	return &Schema{
		migrationsPath: migrationPath,
		installFile:    installFileName + ".sql",
		dropFile:       dropFileName + ".sql",
		log:            log,
		gorm:           gorm,
		debug:          false,
	}, nil
}

// Returns a logrus logger instance, pointed to the provided path
func getLogger(basePath string) *logrus.Logger {
	logger := logrus.New()
	fname, err := filepath.Abs(basePath + "/goose.log")
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		panic(err)
	}
	logger.Out = file
	return logger
}

// Returns the path to the specified migration file.
func (s *Schema) getMigrationPath(file string) (path string) {
	return fmt.Sprintf("%s/%s", s.migrationsPath, file)
}

func (s *Schema) getMigrationHead() (m migration, errs []error) {
	result := s.gorm.First(&m)
	if result.RecordNotFound() {
		s.log.Info("Migrations table is empty")
		m = migration{}
	} else if errs = result.GetErrors(); len(errs) > 0 {
		s.log.Error("Failed to read migrations table.")
		s.log.Errorf("errors: %v\n", errs)
		errs = prependErrors(errs, errors.New("failed to read migrations table"))
		return
	} else {
		s.log.Infof("Migration head: %s", m.FileName)
	}

	return
}

func (s *Schema) setMigrationHead(file string) (errs []error) {
	if file == "" {
		return []error{}
	}
	s.log.Infof("Setting migration head to: %v", file)
	errs = s.exec("TRUNCATE migrations")
	if len(errs) > 0 {
		return
	}
	insert := fmt.Sprintf(`INSERT INTO migrations (file_name) VALUES ("%s")`, file)
	return s.exec(insert)
}

func (s *Schema) exec(statement string) (errs []error) {
	if s.debug {
		s.log.Infof("migration debug: %s\n", statement)
	} else {
		errs = s.gorm.Exec(statement).GetErrors()
	}
	return
}

// Set Schema debug mode.
// In debug mode, migrations will be logged, but not executed.
func (s *Schema) SetDebug(b bool) {
	s.debug = b
	if s.debug {
		s.log.Info("DEBUG MODE ON")
	}
}

// Drop all tables, rebuild the migrations table, and run all migrations.
func (s *Schema) Reset() ([]string, []string, error) {
	dropped, err := s.Drop()
	if err != nil {
		return dropped, nil, err
	}
	s.log.Info("All tables dropped.")
	s.log.Info("Rebuilding migrations table.")
	err = s.Install()
	if err != nil {
		return dropped, nil, err
	}
	s.log.Info("Migrations table created.")

	migrated, err := s.Migrate()
	return dropped, migrated, err
}

// Rebuild the migrations table.
func (s *Schema) Install() error {
	errs := s.exec(`
CREATE TABLE IF NOT EXISTS migrations (
  id        INT(10) UNSIGNED                        NOT NULL AUTO_INCREMENT,
  file_name VARCHAR(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (id)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
`)
	if len(errs) > 0 {
		err := errorsToError(errs)
		s.log.Error(err)
		return err
	}

	return nil
}

// Run the drop file.
func (s *Schema) Drop() ([]string, error) {
	var dropped []string
	path := s.getMigrationPath(s.dropFile)
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	// Run all drop statements.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		// Read the line.
		statement := scanner.Text()
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		// Execute the migration.
		errs := s.exec(statement)
		if len(errs) > 0 {
			return dropped, errorsToError(errs)
		}
		dropped = append(dropped, statement)
	}

	// Record the first migration as the new head, if the migrations table exists
	if s.gorm.HasTable("migrations") {
		errs := s.setMigrationHead(s.installFile)
		if len(errs) > 0 {
			return dropped, errorsToError(errs)
		}
	}

	return dropped, nil
}

// Run any migrations that haven't been run yet.
func (s *Schema) Migrate() ([]string, error) {
	var migrated []string

	// Load the migration files.
	files, err := ioutil.ReadDir(s.getMigrationPath(""))
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	// Get the current migration step.
	head, errs := s.getMigrationHead()
	if len(errs) > 0 {
		s.log.Error(errs)
		return nil, errorsToError(errs)
	}

	var last string
	var runErr error
	for _, file := range files {
		fileName := file.Name()

		// Don't run the drop or install migrations.
		if fileName == s.dropFile || fileName == s.installFile {
			s.log.Infof("Skipping file: %s", fileName)
			continue
		}

		// Only run new migrations.
		if head.FileName != s.installFile && head.FileName >= fileName {
			s.log.Infof("Migration has already been run: %s", fileName)
			continue
		}

		// Read the migration.
		s.log.Infof("Running migration: %s", fileName)
		filePath := s.getMigrationPath(fileName)
		statement, err := ioutil.ReadFile(filePath)
		if err != nil {
			s.log.Error(err)
			runErr = fmt.Errorf("failed on migration %s: %s", fileName, err.Error())
			break
		}

		migrated = append(migrated, fileName)

		// Execute the migration.
		errs := s.exec(string(statement))
		if len(errs) > 0 {
			err = errorsToError(errs)
			s.log.Error(err)
			runErr = fmt.Errorf("failed on migration %s: %s", fileName, err.Error())
			break
		}

		// Remember the last migration that was successfully ran.
		last = fileName
	}

	// Record the last migration as the new head
	errs = s.setMigrationHead(last)
	if len(errs) > 0 {
		s.log.Error(errs)
		return migrated, errorsToError(errs)
	}
	if runErr != nil {
		return migrated, runErr
	}

	return migrated, nil
}

// Create a new migration file with a unique prefix.
func (s *Schema) CreateMigration(name string) (string, error) {
	date := time.Now().Format("20060102150405")
	fullPath, err := filepath.Abs(fmt.Sprintf("%s/%s_%s.sql", s.migrationsPath, date, name))
	if err != nil {
		return "", err
	}

	f, err := os.OpenFile(fullPath, os.O_RDONLY|os.O_CREATE|os.O_EXCL, 0664)
	if err != nil {
		return "", err
	}

	if err := f.Close(); err != nil {
		return "", err
	}

	return fullPath, nil
}

func (s *Schema) DropTable(table string) []error {
	statement := fmt.Sprintf("DROP TABLE IF EXISTS %s", table)
	return s.exec(statement)
}

func (s *Schema) RunMigrationFile(fileName string) []error {
	statement, err := ioutil.ReadFile(s.getMigrationPath(fileName))
	if err != nil {
		s.log.Error(err)
		return []error{err}
	}
	return s.exec(string(statement))
}

func (s *Schema) SetForeignKeyChecks(enabled bool) []error {
	value := 0
	if enabled {
		value = 1
	}
	statement := fmt.Sprintf("SET FOREIGN_KEY_CHECKS=%v", value)
	return s.exec(statement)
}
