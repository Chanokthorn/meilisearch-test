start-ms-ui:
	docker run -d --restart=always --name="meilisearch-ui" -p 24900:24900 riccoxie/meilisearch-ui:latest