package gogo

type Config struct {
	// BOARD
	Dymension  int `mapstructure:"dymension"`
	SquareSize int `mapstructure:"square_size"`
	BorderSize int `mapstructure:"border_size"`

	// CORE
	RandomSeed int64 `mapstructure:"random_seed"`

	// NEURAL NETWORK
	Activation   string `mapstructure:"activation"`
	HiddenLayers []int  `mapstructure:"hidden_layers"`

	// TRAINING
	PopulationSize    int     `mapstructure:"population_size"`
	MutationRate      float64 `mapstructure:"mutation_rate"`
	StabilizationRate float64 `mapstructure:"stabilization_rate"`
	Matches           int     `mapstructure:"matches"`
	AgentsFile        string  `mapstructure:"agents_file"`
}
