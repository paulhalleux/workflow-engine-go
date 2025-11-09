@echo off
echo Cleaning generated Swagger output...

cd "../ui/packages/wf-engine-api"
del /s /q "src/docs"
rmdir /s /q "src/.openapi-generator"
rmdir /s /q "src/docs"
del /s "src/\.openapi-generator-ignore"
