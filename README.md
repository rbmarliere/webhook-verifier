# usage

```
Usage of ./webhook-verifier:
  -port string
    	port to listen on (default "32777")
  -verbose
    	verbose logging
```

# nginx configuration

/etc/nginx/sites-available/your-project-config

```nginx
location /pull {
  proxy_set_header Secret "your_webhook_secret";
  proxy_set_header Project-Root "/path/to/project/root";
  proxy_pass http://127.0.0.1:7777;
}
```

# start server as systemd user daemon

$HOME/.config/systemd/user/webhook-verifier.service

```
[Unit]
Description=Webhook verifier

[Service]
Type=simple
ExecStart=/full/path/to/webhook-verifier

[Install]
WantedBy=default.target
```

`systemctl --user enable webhook-verifier`
`systemctl --user start webhook-verifier`
