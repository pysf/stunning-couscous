package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/manifoldco/promptui"
	"github.com/pysf/stunning-couscous/internal/bulkgen"
	"github.com/pysf/stunning-couscous/internal/partner"
)

// 52.51999140,
// 13.40497255

func main() {
	logger := log.New(os.Stdout, "=> ", log.Ldate|log.Ltime)
	cliArgs := seederArgs{}

	latPrompt := promptui.Prompt{
		Label: "Let's find some partners around a custom location. What latitude are you considering?",
		Validate: func(input string) error {
			if _, err := strconv.ParseFloat(input, 64); err != nil {
				return fmt.Errorf("invalid latitude: %w", err)
			}
			return nil
		},
		Default: "52.51999140",
	}

	lat, err := latPrompt.Run()
	if err != nil {
		log.Fatalf("CLI failed %s\n", err)
	}
	cliArgs.Latitude = lat

	lngPrompt := promptui.Prompt{
		Label: "And what about longitude?",
		Validate: func(input string) error {
			if _, err := strconv.ParseFloat(input, 64); err != nil {
				return fmt.Errorf("invalid longitude: %w", err)
			}
			return nil
		},
		Default: "13.40497255",
	}

	lng, err := lngPrompt.Run()
	if err != nil {
		log.Fatalf("CLI failed %s\n", err)
	}
	cliArgs.Longitude = lng

	sizePrompt := promptui.Prompt{
		Label: "How many partners do you want?",
		Validate: func(input string) error {
			if _, err := strconv.ParseInt(input, 10, 32); err != nil {
				return fmt.Errorf("invalid size: %w", err)
			}
			return nil
		},
		Default: "100",
	}

	size, err := sizePrompt.Run()
	if err != nil {
		log.Fatalf("CLI failed %s\n", err)
	}
	cliArgs.Size = size

	databaseAddressPrompt := promptui.Prompt{
		Label: "What is the Postgre IP address?",
		Validate: func(input string) error {
			if ip := net.ParseIP(input); ip == nil {
				return fmt.Errorf("invalid ip")
			}
			return nil
		},
		Default: "127.0.0.1",
	}

	databeseIP, err := databaseAddressPrompt.Run()
	if err != nil {
		log.Fatalf("CLI failed %s\n", err)
	}
	cliArgs.IP = databeseIP

	databasePortPrompt := promptui.Prompt{
		Label: "What is the Postgre port number?",
		Validate: func(input string) error {
			if _, err := strconv.ParseInt(input, 10, 32); err != nil {
				return fmt.Errorf("invalid Port: %w", err)
			}
			return nil
		},
		Default: "5432",
	}

	databesePort, err := databasePortPrompt.Run()
	if err != nil {
		log.Fatalf("CLI failed %s\n", err)
	}
	cliArgs.Port = databesePort

	databaseNamePrompt := promptui.Prompt{
		Label: "Can you suggest a name for the database? (Please be aware that changing this will impact the main app!)",
		Validate: func(input string) error {
			if len(input) == 0 {
				return fmt.Errorf("invalid db name")
			}
			return nil
		},
		Default: "stunning-couscous",
	}

	databeseName, err := databaseNamePrompt.Run()
	if err != nil {
		log.Fatalf("CLI failed %s\n", err)
	}
	cliArgs.DB = databeseName

	usernamePrompt := promptui.Prompt{
		Label: "May I have the database credentials, Username?",
		Validate: func(input string) error {
			if len(input) == 0 {
				return fmt.Errorf("username can not be empty")
			}
			return nil
		},
		Default: "postgres",
	}

	username, err := usernamePrompt.Run()
	if err != nil {
		log.Fatalf("CLI failed %s\n", err)
	}
	cliArgs.Username = username

	passwordPrompt := promptui.Prompt{
		Label:   "Also, what is the password?",
		Default: "gHteuivwdvkew4wt",
	}

	password, err := passwordPrompt.Run()
	if err != nil {
		fmt.Printf("CLI failed %v\n", err)
		return
	}
	cliArgs.Password = password

	if err := validator.New().Struct(cliArgs); err != nil {
		logger.Fatalf("Validation: Err= %s", err)
	}

	os.Setenv("POSTGRESQL_HOST", cliArgs.IP)
	os.Setenv("POSTGRESQL_PORT", cliArgs.Port)
	os.Setenv("POSTGRESQL_DATABASE", cliArgs.DB)
	os.Setenv("POSTGRESQL_USERNAME", cliArgs.Username)
	os.Setenv("POSTGRESQL_PASSWORD", cliArgs.Password)

	l := partner.Location{}
	if err := partner.FillLocation(cliArgs.Latitude, cliArgs.Longitude, &l); err != nil {
		logger.Fatalf("FillLocation: Err= %s", err)
	}

	s, err := strconv.Atoi(cliArgs.Size)
	if err != nil {
		logger.Fatalf("Atoi: Err= %s", err)
	}

	logger.Printf("Creating random locations around (%v,%v) \n", cliArgs.Latitude, cliArgs.Longitude)
	locations := bulkgen.GenerateRandomLocations(l, s)

	logger.Printf("%v random location generated \n", len(locations))
	partners := bulkgen.GeneratePartner(locations)

	logger.Printf("%v partners generated \n", len(partners))
	partnerRepo, err := partner.NewPartnerRepo()
	if err != nil {
		log.Fatal(err)
	}

	if err = partnerRepo.BulkImport(partners); err != nil {
		log.Fatal(err)
	}
	logger.Printf("%v partners are imported \n", len(partners))
	logger.Println("Done!")
}

type seederArgs struct {
	Latitude  string `validator:"required,latitude"`
	Longitude string `validator:"required,longitude"`
	Size      string `validator:"required,number"`
	IP        string `validator:"required,ip4_addr"`
	Port      string `validator:"required,number"`
	DB        string `validator:"required,alphanumunicode"`
	Username  string `validator:"required"`
	Password  string `validator:"omitempty"`
}
