window.onload = function () {
	//<editor-fold desc="Changeable Configuration Block">

	// the following lines will be replaced by docker/configurator, when it runs in a docker-container
	window.ui = SwaggerUIBundle({
		url: "/api/openapi.yaml",
		dom_id: '#swagger-ui',
		deepLinking: true,
		presets: [
			SwaggerUIBundle.presets.apis,
			SwaggerUIStandalonePreset
		],
		plugins: [
			SwaggerUIBundle.plugins.DownloadUrl
		],
		layout: "StandaloneLayout",

		requestInterceptor: function (req) {
			// Only modify actual API requests, not the initial spec loading
			if (!req.url.endsWith('.yaml') && !req.url.endsWith('.json')) {
				req.url = req.url.replace(window.location.origin, "http://localhost:8080");
			}
			return req;
		}
	});

	//</editor-fold>
};
