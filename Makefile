
all: bnf buildgo 

bnf: src/lang1/lang1.bnf
	( cd src/lang1/; gocc -a lang1.bnf )

buildgo:
	go build
