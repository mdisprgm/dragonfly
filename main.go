package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
)

func main() {
	/*로그*/
	log := logrus.New() //포인터 반환함
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel
	/*로그 끝*/

	//채팅 로그, 아마도..?
	chat.Global.Subscribe(chat.StdoutSubscriber{})

	//구성 불러오기
	config, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}

	//서버
	srv := server.New(&config, log)     //서버 새로 만듦, 구성과 로거(포인터)를 넘겨받음
	srv.CloseOnProgramEnd()             //서버를 안전하게 종료하게 만듦
	if err := srv.Start(); err != nil { //시작 실패하면 프로세스 종료
		log.Fatalln(err)
	}

	//loop
	for {
		if _, err := srv.Accept(); err != nil {
			return
		}
	}
}

// readConfig reads the configuration from the config.toml file, or creates the file if it does not yet exist.
// config.toml에서 구성 불러옴, 파일 없으면 새로 생성
func readConfig() (server.Config, error) {
	c := server.DefaultConfig()                               //(파일 없을 때 새로 만들 파일에 들어갈?) 기본 구성
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) { //존재하지 않으면
		data, err := toml.Marshal(c) //Marshaling == Encoding, JS에서 JSON.stringify와 같은 역할로 추정
		if err != nil {
			return c, fmt.Errorf("failed encoding default config: %v", err) //인코딩 실패 에러
		}
		if err := ioutil.WriteFile("config.toml", data, 0644); err != nil { //파일(구성 파일) 쓰기 실패 에러
			return c, fmt.Errorf("failed creating config: %v", err)
		}
		return c, nil //정상
	}
	data, err := ioutil.ReadFile("config.toml") //구성 읽기
	if err != nil {                             //읽기 실패하면
		return c, fmt.Errorf("error reading config: %v", err) //에러
	}
	if err := toml.Unmarshal(data, &c); err != nil { //디코딩
		return c, fmt.Errorf("error decoding config: %v", err) //디코딩 실패 시 에러
	}
	return c, nil //정상
}
