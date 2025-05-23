NAME	=	stone_analysis

RM		= 	rm -rf

all: build

build:
	go build -buildvcs=false -v -o $(NAME) ./cmd/stone-analysis

test:
	go test -v ./...

clean:
	$(RM) $(NAME)

fclean: clean

re: fclean all
