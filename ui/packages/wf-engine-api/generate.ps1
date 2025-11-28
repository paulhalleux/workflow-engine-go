# cleanup previous generation
Remove-Item -Path ./src/* -Recurse -Force

# generate new API client
docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli generate -i /local/docs/swagger.json -g typescript-fetch -o /local/src --additional-properties=stringEnums=true

# cleanup unnecessary files
Remove-Item -Path ./src/.openapi-generator -Recurse -Force
Remove-Item -Path ./src/.openapi-generator-ignore -Force
Remove-Item -Path ./src/docs -Recurse -Force