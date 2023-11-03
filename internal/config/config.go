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
		RabbitMQ: RabbitMQ{
			Ip:       GetString("rabbitmq.ip"),
			Port:     GetInt("rabbitmq.port"),
			Username: GetString("rabbitmq.username"),
			Password: GetString("rabbitmq.password"),
		},
		QiNiu: QiNiu{
			AccessKey: GetString("qiniu.access_key"),
			SecretKey: GetString("qiniu.secret_key"),
			Bucket:    GetString("qiniu.bucket"),
			Domain:    GetString("qiniu.domain"),
		},
	}
}

type config struct {
	App      App
	JWT      JWT
	Server   Server
	Mysql    Mysql
	Redis    Redis
	RabbitMQ RabbitMQ
	QiNiu    QiNiu
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

// RabbitMQ config.
type RabbitMQ struct {
	Ip       string `mapstructure:"ip"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// QiNiu config.
type QiNiu struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Domain    string `mapstructure:"domain"`
}
