package args

import (
	"argovue/aws"
	"flag"
	"fmt"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

type RcloneParams struct {
	Image    string
	Endpoint string
	Key      string
	Secret   string
	Session  string
	Region   string
	Bucket   string
}

// Args type
type Args struct {
	verboseLevel     string
	port             int
	bindAddr         string
	dir              string
	args             []string
	oidcProvider     string
	oidcClientID     string
	oidcClientSecret string
	oidcRedirectURL  string
	oidcScopes       string
	oidcUserId       string
	uiRootURL        string
	uiRootDomain     string
	k8sNamespace     string
	k8sRelease       string
	k8sService       string
	baseDomain       string
	sessionKey       string
	redisAddr        string
	tlsIssuer        string
	rcloneImage      string
	s3Endpoint       string
	s3Key            string
	s3Secret         string
	s3Region         string
	s3Bucket         string
	s3Principal      string
	s3PrincipalKey   string
}

// New type
func New() *Args {
	return new(Args).Parse()
}

func getEnvOrDefault(name, def string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}
	return def
}

// Parse parameters
func (a *Args) Parse() *Args {
	flag.StringVar(&a.verboseLevel, "verbose", getEnvOrDefault("VERBOSE", "info"), "Set verbosity level")
	flag.IntVar(&a.port, "port", 8080, "Listen port")
	flag.StringVar(&a.bindAddr, "bind", os.Getenv("BIND_ADDR"), "Bind address")
	flag.StringVar(&a.dir, "dir", getEnvOrDefault("UI_DIR", "ui/dist"), "Static files folder")
	flag.StringVar(&a.oidcProvider, "oidc-provider", os.Getenv("OIDC_PROVIDER"), "OIDC provider")
	flag.StringVar(&a.oidcClientID, "oidc-client-id", os.Getenv("OIDC_CLIENT_ID"), "OIDC client id")
	flag.StringVar(&a.oidcClientSecret, "oidc-client-secret", os.Getenv("OIDC_CLIENT_SECRET"), "OIDC client secret")
	flag.StringVar(&a.oidcRedirectURL, "oidc-redirect-url", os.Getenv("OIDC_REDIRECT_URL"), "OIDC redirect url")
	flag.StringVar(&a.oidcScopes, "oidc-scopes", getEnvOrDefault("OIDC_SCOPES", "groups"), "OIDC scopes")
	flag.StringVar(&a.oidcUserId, "oidc-user-id", getEnvOrDefault("OIDC_USER_ID", "preferred_username"), "OIDC user id field")
	flag.StringVar(&a.uiRootURL, "ui-root-url", getEnvOrDefault("UI_ROOT_URL", "http://localhost:8080/ui/#/"), "UI root url for redirects")
	flag.StringVar(&a.k8sNamespace, "k8s-namespace", getEnvOrDefault("K8S_NAMESPACE", "default"), "Kubernetes objects namespace")
	flag.StringVar(&a.k8sRelease, "k8s-release-name", getEnvOrDefault("K8S_RELEASE_NAME", "argovue"), "Release name (to look up config and filebrowser)")
	flag.StringVar(&a.k8sService, "k8s-service-name", getEnvOrDefault("K8S_SERVICE_NAME", "argovue"), "Argovue service name (to forward ingresses)")
	flag.StringVar(&a.baseDomain, "base-domain", getEnvOrDefault("BASE_DOMAIN", ""), "Base domain to base ingress on")
	flag.StringVar(&a.sessionKey, "session-key", getEnvOrDefault("SESSION_KEY", "0123456789abcdef"), "HTTP Session encryption key")
	flag.StringVar(&a.redisAddr, "redis-addr", getEnvOrDefault("REDIS_ADDR", "localhost:6379"), "Redis session store address")
	flag.StringVar(&a.tlsIssuer, "tls-issuer", getEnvOrDefault("TLS_ISSUER", "tpp-venafi-issuer"), "Ingress tls issuer annotation value")
	flag.StringVar(&a.rcloneImage, "rclone-image", getEnvOrDefault("RCLONE_IMAGE", ""), "Rclone image")
	flag.StringVar(&a.s3Endpoint, "s3-endpoint", getEnvOrDefault("S3_ENDPOINT", ""), "AWS S3 endpoint")
	flag.StringVar(&a.s3Key, "s3-key", getEnvOrDefault("S3_KEY", ""), "AWS S3 key")
	flag.StringVar(&a.s3Secret, "s3-secret", getEnvOrDefault("S3_SECRET", ""), "AWS S3 secret")
	flag.StringVar(&a.s3Region, "s3-region", getEnvOrDefault("S3_REGION", ""), "AWS S3 region")
	flag.StringVar(&a.s3Bucket, "s3-bucket", getEnvOrDefault("S3_BUCKET", ""), "AWS S3 bucket")
	flag.StringVar(&a.s3PrincipalKey, "s3-principal-key", getEnvOrDefault("S3_PRINCIPAL_KEY", ""), "AWS S3 principal key")
	flag.StringVar(&a.s3Principal, "s3-principal", getEnvOrDefault("S3_PRINCIPAL", ""), "AWS S3 principal")

	url, _ := url.Parse(a.uiRootURL)
	a.uiRootDomain = fmt.Sprintf("%s://%s", url.Scheme, url.Host)

	flag.Parse()
	a.args = flag.Args()
	return a
}

// Dir to serve files from
func (a *Args) Dir() string {
	return a.dir
}

// BindAddr to bind to Web Server
func (a *Args) BindAddr() string {
	return a.bindAddr
}

// Port to bind to Web Server
func (a *Args) Port() int {
	return a.port
}

// UIRootURL returns UI root url
func (a *Args) UIRootURL() string {
	return a.uiRootURL
}

func (a *Args) UIRootDomain() string {
	return a.uiRootDomain
}

func (a *Args) SessionKey() string {
	return a.sessionKey
}

func (a *Args) OIDC() (string, string, string, string, string) {
	return a.oidcProvider, a.oidcClientID, a.oidcClientSecret, a.oidcRedirectURL, a.oidcScopes
}

func (a *Args) OidcUserId() string {
	return a.oidcUserId
}

func (a *Args) Namespace() string {
	return a.k8sNamespace
}

func (a *Args) Release() string {
	return a.k8sRelease
}

func (a *Args) Service() string {
	return a.k8sService
}

func (a *Args) BaseDomain() string {
	return a.baseDomain
}

func (a *Args) Redis() string {
	return a.redisAddr
}

func (a *Args) TLSIssuer() string {
	return a.tlsIssuer
}

func (a *Args) RcloneParams() (p RcloneParams) {
	p.Image = a.rcloneImage
	p.Endpoint = a.s3Endpoint
	p.Region = a.s3Region
	p.Bucket = a.s3Bucket
	return
}

func (a *Args) AWS() *aws.AWS {
	return &aws.AWS{
		Endpoint:       a.s3Endpoint,
		Key:            a.s3Key,
		Secret:         a.s3Secret,
		Region:         a.s3Region,
		Bucket:         a.s3Bucket,
		PrincipalKey:   a.s3PrincipalKey,
		PrincipalValue: a.s3Principal,
	}
}

// LogLevel set loglevel
func (a *Args) LogLevel() *Args {
	switch a.verboseLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	return a
}
