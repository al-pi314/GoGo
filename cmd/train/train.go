package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/al-pi314/gogo"
	"github.com/al-pi314/gogo/game"
	"github.com/al-pi314/gogo/population"
	"github.com/spf13/viper"
)

func loadConfig() *gogo.Config {
	viper.SetEnvPrefix("X")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	config := gogo.Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(errors.Wrap(err, "viper could not unmarshall enviorment variables"))
	}

	rand.Seed(config.RandomSeed)
	return &config
}

func confirmDirectory(dirPath string) {
	if _, err := os.Stat(dirPath); !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("...output %s directory already exists! Directory files will be removed (5s to cancel)\n", dirPath)
		time.Sleep(5 * time.Second)

		if err := os.RemoveAll(dirPath); err != nil {
			log.Fatal("failed to empty directory")
		}
	}
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to created population directory"))
	}
}

func isArgSet(arg *string) bool {
	return arg != nil && *arg != ""
}

func main() {
	fmt.Println("...starting training")
	config := loadConfig()

	// display config
	b, _ := json.Marshal(config)
	fmt.Println(string(b))

	populationFile := flag.String("population", "", "path to population.json file containing a population")
	outputDirectory := flag.String("output", "", "path to output directory for training")
	flag.Parse()

	// select output directory
	if isArgSet(outputDirectory) {
		config.OutputDirectory = *outputDirectory
		fmt.Println("...using output directory provided in command line")
	}
	fmt.Printf("...output directory is set to %s\n", config.OutputDirectory)

	// confirm output directory
	confirmDirectory(config.OutputDirectory)
	fmt.Println("...output directory confirmed")

	// create population
	currPopulation := population.NewPopulation(config)
	fmt.Println("...population created")
	if isArgSet(populationFile) {
		currPopulation.LoadFromFile(populationFile)
		fmt.Println("...population overwritten from file")
	}
	currPopulation.OutputDirectory = config.OutputDirectory

	// test save population
	currPopulation.Save()
	game.NewGame(game.Game{
		SaveFileName: fmt.Sprintf("%s/games/dummy_game.json", config.OutputDirectory),
	}).Save()
	fmt.Println("...test save completed (check output directory!)")

	// train population
	currPopulation.Train(population.TrainingSettings{
		Rounds:    config.Rounds,
		Groups:    config.Groups,
		KeepBestN: config.KeepBestN,

		SaveInterval:     config.SaveInterval,
		SaveGameInterval: config.SaveGameInterval,
	})

	fmt.Println("...training completed")

	// save population
	currPopulation.Save()
	fmt.Println("...last population saved - finished!")
}
