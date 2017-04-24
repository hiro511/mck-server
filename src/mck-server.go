package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/hiro511/mck"
)

const inputDirName = "inputs"
const doneDirName = "done"
const maxCount = 1000

var basePath string

type server struct {
	InputFile string
	Count     int
}

func (s *server) FetchJobs(ctx context.Context, request *pb.JobRequest) (*pb.Job, error) {
	fmt.Println("FetchJobs request is accepted.")
	if s.InputFile == "" {
		s.InputFile = findInputFile()
	}
	return &pb.Job{"job_name", []byte(""), 3, ""}, nil
}

func (s *server) DownloadMCK(ctx context.Context, request *pb.MCKRequest) (*pb.MolComKit, error) {
	return &pb.MolComKit{[]byte("hoge")}, nil
}

func newServer() *server {
	s := new(server)
	s.InputFile = ""
	s.Count = 0
	return s
}

func main() {
	inputFile, err := findInputFile()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(inputFile)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMckServer(grpcServer, newServer())
	err = grpcServer.Serve(lis)
	if err != nil {
		fmt.Println("hogehoge")
	}
}

func init() {
	var err error
	basePath, err = os.Getwd()
	if err != nil {
		// return errors.New("absolute path couldn't be gotten")
	}
}

func findInputFile() (string, error) {
	// inputDir := filepath.Join(basePath, inputDirName)
	_, err := os.Stat(inputDirName)
	if err != nil {
		return "", errors.New("input directory doesn't exist")
	}

	fileInfos, err := ioutil.ReadDir(inputDirName)
	if err != nil {
		return "", errors.New("couldn't read inputs directory")
	}
	for _, fileInfo := range fileInfos {
		return filepath.Join(inputDirName, fileInfo.Name()), nil
	}

	return "", errors.New("couldn't read inputs directory")
}

func moveFile(filePath, toDir string) error {

	return nil
}
