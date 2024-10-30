.PHONY:cover
cover1:
	go test -short -count=1 -coverprofile=coverage1.out ./pkg/storage/postgress
	go tool cover -html=coverage1.out
	rm coverage1.out
	
cover2:
	go test -short -count=1 -coverprofile=coverage2.out ./pkg/rss/
	go tool cover -html=coverage2.out
	rm coverage2.out
cover3:
	go test -short -count=1 -coverprofile=coverage3.out ./pkg/api/
	go tool cover -html=coverage3.out
	rm coverage3.out