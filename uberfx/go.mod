module github.com/tuantran1810/go-di-template/uberfx

go 1.24.2

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/golang-migrate/migrate/v4 v4.18.3
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/spf13/cobra v1.9.1
	github.com/tuantran1810/go-di-template/libs v0.0.0-00010101000000-000000000000
	go.uber.org/fx v1.24.0
	go.uber.org/zap v1.27.0
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gorm v1.30.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)

replace github.com/tuantran1810/go-di-template/libs => ../libs
