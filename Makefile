NAME	=	stone_analysis

RM		= 	rm -rf

all: $(NAME)

$(NAME):
	go build -buildvcs=false -v -o $(NAME) ./cmd/stone-analysis

test:
	go test -v ./...

clean:
	$(RM) $(NAME)

fclean: clean

re: fclean all

test:
	go test ./...
