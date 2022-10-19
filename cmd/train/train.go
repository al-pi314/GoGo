package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/al-pi314/gogo"
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

func confirmFile(filePath string) *os.File {
	if _, err := os.Stat(filePath); !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("...output %s file already exists! File will be overwritten (3s to cancel).\n", filePath)
		time.Sleep(3 * time.Second)
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		return nil
	}
	return file
}

func isArgSet(arg *string) bool {
	return arg != nil && *arg != ""
}

func main() {
	fmt.Println("...starting training")
	config := loadConfig()

	populationFile := flag.String("population", "", "path to population.json file containing a population")
	outputFile := flag.String("output", "", "path to output file after population training is done")
	flag.Parse()

	// select output file
	if isArgSet(outputFile) {
		config.OutputFile = *outputFile
		fmt.Println("...using output file provided in command line")
	}
	fmt.Printf("...output file is set to %s\n", config.OutputFile)

	// confirm output file
	file := confirmFile(config.OutputFile)
	if file == nil {
		log.Fatal("provided output file cannot be accessed or created!")
	}
	fmt.Println("...output file confirmed")

	// create population
	population := NewPopulation(config)
	fmt.Println("...population created")
	if isArgSet(populationFile) {
		population.LoadFromFile(populationFile)
		fmt.Println("...population overwritten from file")
	}
	population.OutputFileName = config.OutputFile

	// test save population
	population.Save()
	fmt.Println("...test save completed (check output file!)")

	// train population
	population.Train(config.Matches, config.SaveInterval, file)
	fmt.Println("...training completed")

	// save population
	population.Save()
	file.Close()
	fmt.Println("...last population saved - finished!")
}
