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
	Rounds            int     `mapstructure:"rounds"`
	Groups            int     `mapstructure:"groups"`
	KeepBestN         int     `mapstructure:"keep_best_n"`
	SaveInterval      int     `mapstructure:"save_interval"`
	SaveGameInterval  int     `mapstructure:"save_game_interval"`
	OutputDirectory   string  `mapstructure:"output_directory"`
}
