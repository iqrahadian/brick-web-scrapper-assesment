# brick-web-scrapper-assesment

### base architecture
![brick](https://github.com/iqrahadian/brick-web-scrapper-assesment/assets/13548762/c974d681-125e-418a-8765-62a6543f8e56)

### To Execute
- go mod vendor
- go run *.go

### ToDo
1. Handling url that start with ta.tokopedia.com for detail product page, failed to retrieve page because of security issue, resulting in rating & description will be empty
2. implement config to run with postgres, currently using sqlite for simplicity
3. implement config/env

### notes
1. const*.go file is used for testing html scrapper