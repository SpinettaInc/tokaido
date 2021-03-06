package proxy

import (
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/utils"
)

const proxyNetwork = "tokaido_proxy"

// generateNginxConf ...
func generateNginxConf(projectName, domain, proxyPassDomain string) []byte {
	return []byte(`server {
  listen          ` + constants.ProxyPort + ` ssl;
  server_name     ` + projectName + `-toktestdb.` + domain + ` ` + projectName + `.` + domain + `;
  server_tokens   off;

  ssl_certificate           /tokaido/proxy/config/client/tls/wildcard.crt;
  ssl_certificate_key       /tokaido/proxy/config/client/tls/wildcard.key;

  error_page 502 /tokaido-errors/502.html;
  error_page 503 /tokaido-errors/503.html;
  error_page 504 /tokaido-errors/504.html;

  location ^~ /tokaido-errors/ {
    root /tokaido/proxy/config/nginx/errors/;
  }

  location / {
    proxy_pass             ` + proxyPassDomain + `;
    proxy_set_header       Host              $host:$server_port;
    proxy_intercept_errors on;
  }
}
`)
}

// restartNginx sends a HUP signal to nginx
func restartNginx() {
	utils.DebugString("restarting the proxy container's nginx process")
	composeStdout("kill", "-s", "HUP", "proxy")
}
