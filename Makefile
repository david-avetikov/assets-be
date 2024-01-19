reimport:
	go mod tidy

checkVulnerabilities:
	test govulncheck || go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

staticAnalysis:
	test gocritic || go install -v github.com/go-critic/go-critic/cmd/gocritic@latest
	gocritic check ./...

generate:
	make generateJson
	make generateSwagger

generateJson:
	test -f $(GOPATH)/bin/easyjson || go get github.com/mailru/easyjson && go install github.com/mailru/easyjson/...@latest
#common
	cd common/model && $(GOPATH)/bin/easyjson -all pageable.go
	cd common/model && $(GOPATH)/bin/easyjson -all sort.go
	cd common/model && $(GOPATH)/bin/easyjson -all token.go
#asset
	cd modules/asset/model && $(GOPATH)/bin/easyjson -all asset.go
	cd modules/asset/model && $(GOPATH)/bin/easyjson -all cash_flow.go
#attachment
	cd modules/attachment/model && $(GOPATH)/bin/easyjson -all attachment.go
#authorization
	cd modules/authorization/model && $(GOPATH)/bin/easyjson -all authority.go
	cd modules/authorization/model && $(GOPATH)/bin/easyjson -all authorization_response.go
	cd modules/authorization/model && $(GOPATH)/bin/easyjson -all role.go
	cd modules/authorization/model && $(GOPATH)/bin/easyjson -all user.go
#country
	cd modules/country/model && $(GOPATH)/bin/easyjson -all city.go
	cd modules/country/model && $(GOPATH)/bin/easyjson -all country.go
	cd modules/country/model && $(GOPATH)/bin/easyjson -all currency.go

generateSwagger:
	test -f $(GOPATH)/bin/swag || go get -u github.com/swaggo/swag/cmd/swag
	$(GOPATH)/bin/swag init --parseDependency --parseInternal -g cmd/application/main.go

.PHONY: checkVulnerabilities staticAnalysis reimport generate generateJson generateSwagger
