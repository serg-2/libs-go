package commonlib

type DatabaseConfig struct {
	Username string `json:"user"`
	Password string `json:"password"`
	Database string `json:"dbname"`
	Host     string `json:"host"`
}
