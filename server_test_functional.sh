#!/bin/bash

read -r -d '' headers <<EOF
Secret: It's a Secret to Everybody
Project-Root: /tmp/test_repo
X-Hub-Signature-256: sha256=757107ea0eb2509fc211221cce984b8a37570b6d7586c22c46f4379c8b043e17
EOF

read -r -d '' payload <<EOF
Hello, World!
EOF

curl -X POST -H "$headers" -d "Hello, World!" localhost:8080

