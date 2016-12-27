go test -c \
	&& ./vision.test.exe -test.benchtime=10s -test.cpuprofile=prof.out -test.bench=. \
	&& go tool pprof -weblist=. vision.test.exe prof.out
#rm vision.test.exe prof.out
