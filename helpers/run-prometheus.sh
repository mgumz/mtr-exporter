../mtr-exporter -bind :8080 -schedule "@every 10s" -- -c 2 -n example.com &
../mtr-exporter -bind :8081 -schedule "@every 10s" -- -c 2 -n golang.org &
../mtr-exporter -bind :8082 -schedule "@every 10s" -- -c 2 -n prometheus.io &
../mtr-exporter -bind :8083 -schedule "@every 10s" -- -c 2 -n www.bitwizard.nl &

prometheus --config.file=prometheus.yml
