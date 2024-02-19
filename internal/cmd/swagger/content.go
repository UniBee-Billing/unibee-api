package swagger

const (
	LatestSwaggerUIPageContent = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <meta name="description" content="SwaggerUI"/>
  <title>UniBee API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@latest/swagger-ui.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@latest/swagger-ui-bundle.js" crossorigin></script>
<script>
	window.onload = () => {
		window.ui = SwaggerUIBundle({
			url:    'api.json',
			dom_id: '#swagger-ui',
			deepLinking: true,
			plugins: [SwaggerUIBundle.plugins.DownloadUrl],
			filter: true,
			tagsSorter: 'alpha',
			tryItOutEnabled: true,
			queryConfigEnabled: true, // keeps the selected ?urls.primaryName=...
		});
	};
</script>
</body>
</html>
`
	V3SwaggerUIPageContent = `
<html>
  <head>
    <script src="https://unpkg.com/swagger-ui-dist@3/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@3/swagger-ui-standalone-preset.js"></script>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@3/swagger-ui.css" />
    <title>UniBee API</title>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script defer>
      window.onload = function () {
        const ui = SwaggerUIBundle({
          urls: [
            {
              name: "UniBee Open API Spec",
              url: "api.json",
            },
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
        window.ui = ui;
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
