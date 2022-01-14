# go build -gcflags '-N -l' test_qsort.go
go build test_qsort.go
./test_qsort > example3.xml
bash test_show_graph.sh

