package kassadin

const (
	ModeDev  = "dev"
	ModeProd = "prod"

	ProvidorFile  = "file"
	ProvidorRedis = "redis"

	RotateByDay  = "day"
	RotateBySize = "size"
)

type appConfig struct {
	Name string
	Mode string
	Port int

	Sess  sessionConfig
	Log   logConfig
	Mysql mysqlConfig
	Redis redisConfig
}

type sessionConfig struct {
	Providor  string
	StorePath string
	Enable    bool
}

type logConfig struct {
	Name        string
	Providor    string
	RotateMode  string
	RotateLimit string
	LogPath     string
}

type mysqlConfig struct {
}

type redisConfig struct {
}
