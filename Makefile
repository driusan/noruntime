all:
	cd runtime && go build -o ../runtime.a .
	cd ..
	go tool asm -o main.o main_$(GOOS)_$(GOARCH).s
	go tool link -g -L . -w -E main main.o


clean:
	rm -f a.out *.o runtime.a
