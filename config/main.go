package config

import (
	"os"
	"strconv"

	"github.com/thedevsir/frame-backend/services/utils"
)

var (
	Port             string
	Mode             string
	CorsAllowOrigins string
	RoutesBodyLimit  string

	DBAddress  string
	DBName     string
	DBUsername string
	DBPassword string
	DBSource   string

	AbuseIP         int
	AbuseIPUsername int

	SigningKey      string
	AdminSigningKey string

	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string

	EmailAppName string
	EmailFrom    string

	EmailThemeName      string
	EmailThemeLink      string
	EmailThemeLogo      string
	EmailThemeCopyright string

	EmailVerifyLink string
	EmailResetLink  string

	MinioEndpoint        string
	MinioBuckets         string
	MinioAccessKeyID     string
	MinioSecretAccessKey string

	AvatarPictureFormats string
	AvatarPictureMaxSize int64
)

func Composer(envPath string) (err error) {

	utils.Env(envPath)

	Port = os.Getenv("PORT")
	Mode = os.Getenv("MODE")
	CorsAllowOrigins = os.Getenv("CORS_ALLOW_ORIGINS")
	RoutesBodyLimit = os.Getenv("ROUTES_BODY_LIMIT")

	DBAddress = os.Getenv("DB_ADDRESS")
	DBName = os.Getenv("DB_NAME")
	DBUsername = os.Getenv("DB_USERNAME")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBSource = os.Getenv("DB_SOURCE")

	AbuseIP, err = strconv.Atoi(os.Getenv("ABUSE_IP"))
	if err != nil {
		panic(err)
	}

	AbuseIPUsername, err = strconv.Atoi(os.Getenv("ABUSE_IP_USERNAME"))
	if err != nil {
		panic(err)
	}

	SigningKey = os.Getenv("SIGNING_KEY")
	AdminSigningKey = os.Getenv("ADMIN_SIGNING_KEY")

	SMTPHost = os.Getenv("SMTP_HOST")
	SMTPPort, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		panic(err)
	}

	SMTPUsername = os.Getenv("SMTP_USERNAME")
	SMTPPassword = os.Getenv("SMTP_PASSWORD")

	EmailAppName = os.Getenv("EMAIL_APP_NAME")
	EmailFrom = os.Getenv("EMAIL_FROM")

	EmailThemeName = os.Getenv("EMAIL_THEME_NAME")
	EmailThemeLink = os.Getenv("EMAIL_THEME_LINK")
	EmailThemeLogo = os.Getenv("EMAIL_THEME_LOGO")
	EmailThemeCopyright = os.Getenv("EMAIL_THEME_COPYRIGHT")

	EmailVerifyLink = os.Getenv("EMAIL_VERIFY_LINK")
	EmailResetLink = os.Getenv("EMAIL_RESET_LINK")

	MinioEndpoint = os.Getenv("MINIO_ENDPOINT")
	MinioBuckets = os.Getenv("MINIO_BUCKETS")
	MinioAccessKeyID = os.Getenv("MINIO_ACCESS_KEY_ID")
	MinioSecretAccessKey = os.Getenv("MINIO_SECRET_ACCESS_KEY")

	AvatarPictureFormats = os.Getenv("AVATAR_PICTURE_FORMATS")
	AvatarPictureMaxSize, err = strconv.ParseInt(os.Getenv("AVATAR_PICTURE_MAX_SIZE"), 10, 64)
	if err != nil {
		panic(err)
	}

	return nil
}
