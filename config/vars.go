package config

// если нужно по умолчанию имя,
// используется в Config как имя файла конфиг и в repo//DbSelf как имя БД
const Name = "kituxlsx"

var Mode = "development"

// This should preferably be set at build time via build scripts
// Set during build (adjust module path): go build -ldflags "-X 'kituxlsx/config.ExeVersion=v1.0.0'"
const ExeVersion string = "0.0.1"
