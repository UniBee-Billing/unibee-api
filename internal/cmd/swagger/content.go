package swagger

const (
	RedoclyContent = `
<!DOCTYPE html>
<html>
<head>
<title>API Reference</title>
<meta charset="utf-8"/>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
body {
margin:  0;
padding: 0;
}
</style>
</head>
<body>
<redoc spec-url="api.json" show-object-schema-examples="true"></redoc>
<script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"> </script>
</body>
</html>
`
	LatestSwaggerUIPageContent = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<meta name="description" content="SwaggerUI"/>
	<script src="https://unpkg.com/swagger-ui-dist@3/swagger-ui-standalone-preset.js"></script>
	<title>UniBee API</title>
	<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@latest/swagger-ui.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@latest/swagger-ui-bundle.js" crossorigin></script>
<script>
	window.onload = () => {
		window.ui = SwaggerUIBundle({
			//url:    'api.json',
			urls: [
				//{
				// name: "UniBee Api Spec",
				// url: "api.json",
				//},
				{
				 name: "UniBee Merchant Portal Api Spec (Merchant Open API) ",
				 url: "api.sdk.generator.json",
				},
				{
				 name: "UniBee User Portal Api Spec (Web Component) ",
				 url: "api.user.portal.generator.json",
				},
			],
			dom_id: '#swagger-ui',
			deepLinking: true,
			presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
			plugins: [SwaggerUIBundle.plugins.DownloadUrl],
			layout: "StandaloneLayout",
			filter: true,
			tagsSorter: 'alpha',
			tryItOutEnabled: true,
			queryConfigEnabled: true, // keeps the selected ?urls.primaryName=...
		});
	};
</script>
</body>
<style>
  .swagger-ui .topbar .download-url-wrapper input[type="text"] {
	border: 2px solid #77889a;
  }
  .swagger-ui .topbar .download-url-wrapper .download-url-button {
	background: #77889a;
  }
  .swagger-ui img {
	display: none;
  }
  .swagger-ui .topbar {
	background-color: #ededed;
	border-bottom: 2px solid #c1c1c1;
  }
  .swagger-ui .topbar .download-url-wrapper .select-label {
	color: #3b4151;
  }
</style>
</html>
`
	V3SwaggerUIPageContent = `
<html>
<head>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<meta name="description" content="SwaggerUI"/>
	<script src="https://unpkg.com/swagger-ui-dist@3/swagger-ui-standalone-preset.js"></script>
	<title>UniBee API</title>
	<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@3/swagger-ui.css" />
</head>
<script src="https://unpkg.com/swagger-ui-dist@3/swagger-ui-bundle.js"></script>
<body>
<div id="swagger-ui"></div>
<script>
	window.onload = function () {
		window.ui = SwaggerUIBundle({
			urls: [
				{
				 name: "UniBee Api Spec",
				 url: "api.json",
				},
				//{
				//  name: "UniBee Merchant Portal Api Spec (Merchant Open API) ",
				//  url: "api.sdk.generator.json",
				//},
				//{
				//  name: "UniBee User Portal Api Spec (Web Component) ",
				//  url: "api.user.portal.generator.json",
				//},
			],
			dom_id: "#swagger-ui",
			deepLinking: true,
			presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
			plugins: [SwaggerUIBundle.plugins.DownloadUrl],
			layout: "StandaloneLayout",
			filter: true,
			tagsSorter: 'alpha',
			tryItOutEnabled: true,
			queryConfigEnabled: true, // keeps the selected ?urls.primaryName=...
		});
	};
</script>
<style>
  .swagger-ui .topbar .download-url-wrapper input[type="text"] {
	border: 2px solid #77889a;
  }
  .swagger-ui .topbar .download-url-wrapper .download-url-button {
	background: #77889a;
  }
  .swagger-ui img {
	display: none;
  }
  .swagger-ui .topbar {
	background-color: #ededed;
	border-bottom: 2px solid #c1c1c1;
  }
  .swagger-ui .topbar .download-url-wrapper .select-label {
	color: #3b4151;
  }
</style>
</body>
</html>
`
)
