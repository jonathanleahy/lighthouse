IMAGE=loadimpact/k6

get-health:
	@docker run --rm --network host -i ${IMAGE} run --vus 10 --duration 10s  - <get-health.js

all: get-health
