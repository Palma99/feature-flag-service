package main

import (
	"bufio"
	"database/sql"
	"errors"
	"os"
	"regexp"
	"strconv"

	_ "github.com/lib/pq"

	"fmt"
	"strings"

	"github.com/Palma99/feature-flag-service/internals/application/services"
	usecase "github.com/Palma99/feature-flag-service/internals/application/usecase"
	infrastructure "github.com/Palma99/feature-flag-service/internals/infrastructure/repository"
	"github.com/Palma99/feature-flag-service/internals/interfaces"
)

var db *sql.DB

func init() {
	dbConn, err := sql.Open("postgres", "postgresql://root:root@localhost:5432/local_feature_flag?sslmode=disable")
	if err != nil {
		panic(err)
	}

	db = dbConn
}

func main() {

	keyService := services.NewKeyService()

	environmentRepository := infrastructure.NewPgEnvironmentRepository(db)
	environmentInteractor := usecase.NewEnvironmentInteractor(environmentRepository, keyService)
	environmentController := interfaces.NewEnvironmentInteractor(environmentInteractor)

	flagRepository := infrastructure.NewPgFlagRepository(db)
	flagInteractor := usecase.NewFlagInteractor(flagRepository, environmentRepository, keyService)

	flagController := interfaces.NewFlagController(flagInteractor)

	fmt.Println("Feature flag!")

	for {
		fmt.Print("> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		cmd := scanner.Text()

		switch strings.Split(cmd, " ")[0] {
		case "get":
			if err := handleGetFlags(cmd, flagController); err != nil {
				fmt.Println(err)
			}
		case "create-env":
			if err := handleCreateEnvironment(cmd, environmentController); err != nil {
				fmt.Println(err)
			}
		case "update-flag-value":
			if err := handleUpdateFlagValue(cmd, flagController); err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Println("Unknown command")
		}

	}
}

func handleUpdateFlagValue(cmd string, fc *interfaces.FlagController) error {
	value, key, flagId := "", "", ""

	var re = regexp.MustCompile(`(?m)(--key|--value|--flagId)\s+([aA-zZ|0-9|-]+)`)

	for _, match := range re.FindAllString(cmd, -1) {
		if strings.Index(match, "--key") == 0 {
			key = strings.Split(match, " ")[1]
		} else if strings.Index(match, "--value") == 0 {
			value = strings.Split(match, " ")[1]
		} else if strings.Index(match, "--flagId") == 0 {
			flagId = strings.Split(match, " ")[1]
		}
	}

	if key == "" || value == "" || flagId == "" {
		return errors.New("please specify a key, a value and a flag id")
	}

	flagIdNum, _ := strconv.Atoi(flagId)
	valueBool, _ := strconv.ParseBool(value)

	return fc.UpdateFlagValue(key, flagIdNum, valueBool)
}

func handleCreateEnvironment(cmd string, ec *interfaces.EnvironmentController) error {
	name, projectId := "", ""

	var re = regexp.MustCompile(`(?m)(--name|--projectId)\s+(\w+)`)

	for _, match := range re.FindAllString(cmd, -1) {
		if strings.Index(match, "--name") == 0 {
			name = strings.Split(match, " ")[1]
		} else if strings.Index(match, "--projectId") == 0 {
			projectId = strings.Split(match, " ")[1]
		}
	}

	if name == "" || projectId == "" {
		return errors.New("please specify a name and a project id")
	}

	err := ec.CreateEnvironment(name, projectId)
	if err != nil {
		return err
	}

	return nil
}

func handleGetFlags(cmd string, fc *interfaces.FlagController) error {
	args := strings.Split(cmd, "--key ")[1:]

	if len(args) != 1 {
		return errors.New("please specify a key")
	}

	key := args[0]
	flags, err := fc.GetAllFlagsByEnvironmentKey(key)
	if err != nil {
		return err
	}

	if len(flags) == 0 {
		return errors.New("no flags found")
	}

	for _, flag := range flags {
		fmt.Printf("Name: %s, Enabled: %t, Env: %v\n", flag.Name, flag.Enabled, flag.Environment)
	}

	return nil
}
