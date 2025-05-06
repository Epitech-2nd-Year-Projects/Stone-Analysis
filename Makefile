NAME	=	stone_analysis

RM		= 	rm -rf

all: $(NAME)

$(NAME):
	go build -buildvcs=false -o $(NAME) ./cmd/stone-analysis

clean:
	$(RM) $(NAME)

fclean: clean

re: fclean all

test:
	go test ./...
