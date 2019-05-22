package goose

func ExampleMigration() MigrationInterface {
	return &Migration{
		Name: "example_migration",
		Up: "CREATE TABLE example",
		Down: "",
	}
}


func Up() {
	table := CreateTable("example")
	table.SetColumns([]Column{
		{
			Name: "id",
			Type: "binary",
			Length: 10,
		},
	})
}


