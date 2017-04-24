package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/hiro511/mck"
)

const inputDirName = "inputs"
const doneDirName = "done"
const resultDirName = "results"
const maxCount = 30

var basePath string
var mckName string

type server struct {
	InputFileName string
	InputFile     string
	Count         int32
	MCKName       string
}

func (s *server) FetchJobs(ctx context.Context, request *pb.JobRequest) (*pb.Job, error) {
	fmt.Println("FetchJobs is called.")
	if s.InputFile == "" || s.Count >= maxCount {
		var err error
		if s.InputFile != "" {
			moveToDone(s.InputFileName)
		}
		s.InputFile, err = findInputFile()
		if err != nil {
			return nil, err
		}
		s.InputFileName = filepath.Base(s.InputFile)
		s.Count = 0
	}

	rest := maxCount - s.Count
	count := request.NumRequest
	if rest < request.NumRequest {
		count = rest
	}
	s.Count += count

	inputFileData, err := ioutil.ReadFile(s.InputFile)
	if err != nil {
		return nil, err
	}

	return &pb.Job{s.InputFileName, inputFileData, count, mckName}, nil
}

func (s *server) DownloadMCK(ctx context.Context, request *pb.MCKRequest) (*pb.MolComKit, error) {
	fmt.Println("downloadMCK is called.")
	mck, err := ioutil.ReadFile(request.Name)
	log.Printf("read file: %v", request.Name)
	if err != nil {
		return nil, errors.New("could not find the specified molcomkit")
	}
	return &pb.MolComKit{mck}, nil
}

func (s *server) SendResult(ctx context.Context, result *pb.JobResult) (*pb.Empty, error) {
	_, err := os.Stat(resultDirName)
	if err != nil {
		os.Mkdir(resultDirName, 0644)
	}
	resultFile := filepath.Join(resultDirName, result.Name)
	var fp *os.File
	fp, err = os.OpenFile(resultFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	writer := bufio.NewWriter(fp)
	_, err = writer.WriteString(result.Result)
	if err != nil {
		return nil, err
	}
	writer.Flush()

	return &pb.Empty{}, nil
}

func newServer() *server {
	s := new(server)
	s.InputFile = ""
	s.Count = 0
	return s
}

func main() {
	flag.Parse()
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
	flag.StringVar(&mckName, "name", "MolComKit_v2.4.3", "specify the name of MolComKit")
}

func findInputFile() (string, error) {
	// inputDir := filepath.Join(basePath, inputDirName)
	_, err := os.Stat(inputDirName)
	if err != nil {
		return "", err
	}

	fileInfos, err := ioutil.ReadDir(inputDirName)
	if err != nil {
		return "", err
	}
	for _, fileInfo := range fileInfos {
		if filepath.Ext(fileInfo.Name()) != ".txt" {
			continue
		}
		return filepath.Join(inputDirName, fileInfo.Name()), nil
	}

	return "", errors.New("inputs directory is empty")
}

func moveToDone(fileName string) error {
	err := os.Rename(filepath.Join(inputDirName, fileName), filepath.Join(doneDirName, fileName))
	if err != nil {
		return err
	}

	return nil
}
