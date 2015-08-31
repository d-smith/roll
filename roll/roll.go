package roll

type Core struct{}

type CoreConfig struct{}

func NewCore(config *CoreConfig) *Core {
	return &Core{}
}
