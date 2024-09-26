package swagger

const (
	swaggerUIDocURLPlaceHolder = `{SwaggerUIDocUrl}`
	swaggerUITemplate          = `
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
		<redoc spec-url="{SwaggerUIDocUrl}" show-object-schema-examples="true"></redoc>
		<script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"> </script>
	</body>
</html>
`
)
