package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

// update brew and install postgresql
// brew update
// brew install postgresql
// get the postgresql dump file
// pg_dump -U PGUSER -h PGHOST -p PGPORT -W -F t PGDATABASE > YOUR_FILENAME_HERE
// to restore
// pg_restore -U <username> -h <host> -p <port> -W -F t -d <db_name> <dump_file_name>
// OR in other words: pg_restore -U postgres -h containers-us-west-15.railway.app -p 6473 -W -F t -d railway mydatabasebackup
func main() {	
	err := godotenv.Load()
	if err != nil {
			log.Fatal("Error loading .env file:", err)
	}
	
	brewExt := os.Getenv("brew_EXT")

	// terminal command
	brewUpdate := exec.Command(brewExt, "update")

	brewInstallPostgresql := exec.Command(brewExt, "update")

	// pg_dump := os.Getenv("pg_dump")
	PGUSER := os.Getenv("PGUSER")
	PGHOST := os.Getenv("PGHOST")
	PGPORT := os.Getenv("PGPORT")
	PGDATABASE := os.Getenv("PGDATABASE")
	YOUR_FILENAME_HERE := os.Getenv("YOUR_FILENAME_HERE")
	PGPASSWORD := os.Getenv("PGPASSWORD")
	
	pg_dump_command := exec.Command("pg_dump", "-U", PGUSER, "-h", PGHOST, "-p", PGPORT, "-W", "-F", "t", PGDATABASE, ">", YOUR_FILENAME_HERE)

	// Run the command and capture its outputBrewUpdate
	outputBrewUpdate, errBrewUpdate := brewUpdate.CombinedOutput()
	if errBrewUpdate != nil {
		fmt.Printf("Error: %v\n", errBrewUpdate)
		return
	}

	outputBrewInstallPostgresql, errBrewInstallPostgresql := brewInstallPostgresql.CombinedOutput()
	if errBrewInstallPostgresql != nil {
		fmt.Printf("Error: %v\n", errBrewInstallPostgresql)
		return
	}

	output_pg_dump_command, err_pg_dump_command := pg_dump_command.CombinedOutput()
	if(err_pg_dump_command != nil){
		fmt.Printf("Error: %v\n", err_pg_dump_command)
		return
	}

	// Split the terminal output into lines and print each line
	outputLinesBrewUpdate := strings.Split(string(outputBrewUpdate), "\n")
	for _, line := range outputLinesBrewUpdate {
		fmt.Println(line)
	}
	outputLinesBrewInstallPostgresql := strings.Split(string(outputBrewInstallPostgresql), "\n")
	for _, line := range outputLinesBrewInstallPostgresql {
		fmt.Println(line)
	}
	output_lines_pg_dump_command := strings.Split(string(output_pg_dump_command),"\n")
	for _, line := range output_lines_pg_dump_command {
		fmt.Println(line)
	}
	// logic for pg_dump
	// Open a file for writing the output
	file, err := os.Create(YOUR_FILENAME_HERE)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer file.Close()
	pg_dump_command.Stdout = file

	// Create a reader with the password and set it as the Stdin for the command
	passwordReader := strings.NewReader(PGPASSWORD)
	pg_dump_command.Stdin = passwordReader

	// Capture and print the standard error output
	stderr, err := pg_dump_command.StderrPipe()
	if err != nil {
		fmt.Printf("Error creating stderr pipe: %v\n", err)
		return
	}

	// Start the command
	if err := pg_dump_command.Start(); err != nil {
		fmt.Printf("Error starting pg_dump: %v\n", err)
		return
	}

	// Print any errors encountered during execution
	errOutput, _ := io.ReadAll(stderr)
	fmt.Printf("Errors encountered during pg_dump:\n%s\n", string(errOutput))

	// Wait for the command to finish
	if err := pg_dump_command.Wait(); err != nil {
		fmt.Printf("Error running pg_dump: %v\n", err)
		return
	}

	fmt.Printf("pg_dump completed successfully. Output saved to %s\n", YOUR_FILENAME_HERE)

}