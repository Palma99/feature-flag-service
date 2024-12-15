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

	"github.com/Palma99/feature-flag-service/config"
	"github.com/Palma99/feature-flag-service/internals/application/services"
	usecase "github.com/Palma99/feature-flag-service/internals/application/usecase"
	infrastructure "github.com/Palma99/feature-flag-service/internals/infrastructure/repository"
	interfaces "github.com/Palma99/feature-flag-service/internals/interfaces/cli"
)

var db *sql.DB

var applicationConfig *config.Config

var loggedUserToken string = ""

func init() {

	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	applicationConfig = config

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		applicationConfig.DB.User,
		applicationConfig.DB.Password,
		applicationConfig.DB.Host,
		applicationConfig.DB.Port,
		applicationConfig.DB.Database,
	)

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	db = dbConn
}

func main() {
	keyService := services.NewKeyService()
	passwordService := services.NewPasswordService()
	jwtService := services.NewJWTService(applicationConfig.JwtSecret)

	userRepository := infrastructure.NewPgUserRepository(db)
	authInteractor := usecase.NewAuthInteractor(userRepository, passwordService, jwtService)
	authController := interfaces.NewUserController(authInteractor)

	projectRepository := infrastructure.NewPgProjectRepository(db)
	projectInteractor := usecase.NewProjectInteractor(projectRepository, keyService)
	projectController := interfaces.NewProjectController(projectInteractor)

	environmentRepository := infrastructure.NewPgEnvironmentRepository(db)
	environmentInteractor := usecase.NewEnvironmentInteractor(environmentRepository, projectRepository, keyService)
	environmentController := interfaces.NewEnvironmentInteractor(environmentInteractor)

	flagRepository := infrastructure.NewPgFlagRepository(db)
	flagInteractor := usecase.NewFlagInteractor(flagRepository, environmentRepository, keyService, projectRepository)

	flagController := interfaces.NewFlagController(flagInteractor)

	fmt.Println("Feature flag!")

	for {
		fmt.Print("> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		cmd := strings.TrimSpace(scanner.Text())

		loggedUserId, authError := authInteractor.ValidateToken(loggedUserToken)

		switch strings.Split(cmd, " ")[0] {
		// case "get":
		// 	if authError != nil {
		// 		fmt.Println("You are not logged in")
		// 		continue
		// 	}
		// 	if err := handleGetFlags(cmd, flagController); err != nil {
		// 		fmt.Println(err)
		// 	}
		case "create-project":
			if authError != nil {
				fmt.Println("You are not logged in")
				continue
			}
			if err := handleCreateProject(cmd, projectController, loggedUserId); err != nil {
				fmt.Println(err)
			}
		case "create-env":
			if authError != nil {
				fmt.Println("You are not logged in")
				continue
			}
			if err := handleCreateEnvironment(cmd, environmentController, loggedUserId); err != nil {
				fmt.Println(err)
			}
		case "update-flag-value":
			if authError != nil {
				fmt.Println("You are not logged in")
				continue
			}
			if err := handleUpdateFlagValue(cmd, flagController); err != nil {
				fmt.Println(err)
			}
		case "auth":
			if err := handleAuthCommand(cmd, authController); err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Println("Unknown command")
		}

	}
}

func handleAuthCommand(cmd string, uc *interfaces.AuthController) error {
	action := strings.Split(cmd, " ")[1]

	switch action {
	case "create-user":
		if err := handleCreateUser(cmd, uc); err != nil {
			return err
		}
	case "get-token":
		if err := handleGetToken(cmd, uc); err != nil {
			return err
		}
	default:
		return errors.New("unknown action")
	}

	return nil
}

func handleGetToken(cmd string, uc *interfaces.AuthController) error {
	email, password := "", ""

	var re = regexp.MustCompile(`(?m)(--email|--password)\s+([\w|@|.]+)`)

	for _, match := range re.FindAllString(cmd, -1) {
		if strings.Index(match, "--email") == 0 {
			email = strings.Split(match, " ")[1]
		} else if strings.Index(match, "--password") == 0 {
			password = strings.Split(match, " ")[1]
		}
	}

	if email == "" || password == "" {
		return errors.New("please specify an email and a password")
	}

	userLoginDTO := usecase.UserLoginDTO{
		Email:    email,
		Password: password,
	}

	token, err := uc.GetToken(userLoginDTO)
	if err != nil {
		return err
	}

	loggedUserToken = *token

	fmt.Println(*token)

	fmt.Println("Authorized.")
	return nil
}

func handleCreateUser(cmd string, uc *interfaces.AuthController) error {
	email, nickname, password := "", "", ""

	var re = regexp.MustCompile(`(?m)(--email|--nickname|--password)\s+([\w|@|.]+)`)

	for _, match := range re.FindAllString(cmd, -1) {
		if strings.Index(match, "--email") == 0 {
			email = strings.Split(match, " ")[1]
		} else if strings.Index(match, "--nickname") == 0 {
			nickname = strings.Split(match, " ")[1]
		} else if strings.Index(match, "--password") == 0 {
			password = strings.Split(match, " ")[1]
		}
	}

	if email == "" || nickname == "" || password == "" {
		return errors.New("please specify an email, a nickname and a password")
	}

	userDTO := usecase.CreateUserDTO{
		Email:    email,
		Nickname: nickname,
		Password: password,
	}

	return uc.CreateUser(userDTO)
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

	// return fc.UpdateFlagValue(key, flagIdNum, valueBool)
	fmt.Println("Flag updated.", key, flagIdNum, valueBool)
	return nil
}

func handleCreateProject(cmd string, pc *interfaces.ProjectController, userId int) error {
	name := ""

	var re = regexp.MustCompile(`(?m)(--name)\s+(\w+)`)

	for _, match := range re.FindAllString(cmd, -1) {
		if strings.Index(match, "--name") == 0 {
			name = strings.Split(match, " ")[1]
		}
	}

	if name == "" {
		return errors.New("please specify a name")
	}

	err := pc.CreateProject(name, userId)

	if err != nil {
		return err
	}

	fmt.Println("Project created.")

	return nil
}

func handleCreateEnvironment(cmd string, ec *interfaces.EnvironmentController, userId int) error {
	name := ""
	var projectId int64

	var re = regexp.MustCompile(`(?m)(--name|--projectId)\s+(\w+)`)

	for _, match := range re.FindAllString(cmd, -1) {
		if strings.Index(match, "--name") == 0 {
			name = strings.Split(match, " ")[1]
		} else if strings.Index(match, "--projectId") == 0 {
			projectId, _ = strconv.ParseInt(strings.Split(match, " ")[1], 10, 64)
		}
	}

	if name == "" || projectId == 0 {
		return errors.New("please specify a name and a project id")
	}

	err := ec.CreateEnvironment(name, projectId, userId)
	if err != nil {
		return err
	}

	return nil
}

// func handleGetFlags(cmd string, fc *interfaces.FlagController) error {
// 	args := strings.Split(cmd, "--key ")[1:]

// 	if len(args) != 1 {
// 		return errors.New("please specify a key")
// 	}

// 	key := args[0]
// 	flags, err := fc.GetAllFlagsByEnvironmentKey(key)
// 	if err != nil {
// 		return err
// 	}

// 	if len(flags) == 0 {
// 		return errors.New("no flags found")
// 	}

// 	for _, flag := range flags {
// 		fmt.Printf("Name: %s, Enabled: %t, Env: %v\n", flag.Name, flag.Enabled, flag.Environment)
// 	}

// 	return nil
// }
