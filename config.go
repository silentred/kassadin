package kassadin

const (
	ModeDev  = "dev"
	ModeProd = "prod"

	ProvidorFile  = "file"
	ProvidorRedis = "redis"

	RotateByDay  = "day"
	RotateBySize = "size"
)

type AppConfig struct {
	Name string
	Mode string
	Port int

	Sess  SessionConfig
	Log   LogConfig
	Mysql MysqlConfig
	Redis RedisConfig
}

type SessionConfig struct {
	Providor  string
	StorePath string
	Enable    bool
}

type LogConfig struct {
	Name        string
	Providor    string
	RotateMode  string
	RotateLimit string
	LogPath     string
}

type MysqlConfig struct {
}

type RedisConfig struct {
}
