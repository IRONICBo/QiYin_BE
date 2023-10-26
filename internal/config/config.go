package config

// Config global config instance.
var Config *config

// ConfigInit init config.
func ConfigInit(configPath string) {
	// init viper
	initViper(configPath)

	// init Configuration
	Config = &config{
		App: App{
			Version: GetString("app.version"),
			Debug:   GetBool("app.debug"),
			LogFile: GetString("app.log_file"),
		},
		JWT: JWT{
			Secret:     GetString("jwt.secret"),
			Issuer:     GetString("jwt.issuer"),
			ExpireDays: GetInt("jwt.expire_days"),
		},
		Server: Server{
			Ip:          GetString("server.ip"),
			Port:        GetInt("server.port"),
			MaxFileSize: GetIntOrDefault("server.max_file_size", 10),
		},
		Mysql: Mysql{
			Ip:           GetString("mysql.ip"),
			Port:         GetInt("mysql.port"),
			Username:     GetString("mysql.username"),
			Password:     GetString("mysql.password"),
			Database:     GetString("mysql.database"),
			MaxLifetime:  GetIntOrDefault("mysql.max_lifetime", 120),
			MaxOpenConns: GetIntOrDefault("mysql.max_open_conns", 100),
			MaxIdleConns: GetIntOrDefault("mysql.max_idle_conns", 20),
		},
		Redis: Redis{
			Ip:       GetString("redis.ip"),
			Port:     GetInt("redis.port"),
			Password: GetString("redis.password"),
			Database: GetInt("redis.database"),
		},
	}
}

type config struct {
	App    App
	JWT    JWT
	Server Server
	Mysql  Mysql
	Redis  Redis
}

// App config.
type App struct {
	Version string `mapstructure:"version"`
	Debug   bool   `mapstructure:"debug"`
	LogFile string `mapstructure:"log_file"`
}

// JWT config.
type JWT struct {
	Secret     string `mapstructure:"secret"`
	Issuer     string `mapstructure:"issuer"`
	ExpireDays int    `mapstructure:"expire_days"`
}

// Server config.
type Server struct {
	Ip          string `mapstructure:"ip"`
	Port        int    `mapstructure:"port"`
	MaxFileSize int    `mapstructure:"max_file_size"`
}

// Mysql config.
type Mysql struct {
	Ip           string `mapstructure:"ip"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	MaxLifetime  int    `mapstructure:"max_lifetime"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// Redis config.
type Redis struct {
	Ip       string `mapstructure:"ip"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}
