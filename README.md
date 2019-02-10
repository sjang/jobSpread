# jobSpread
jobSpread는 하나의 요청을 여러 작업으로 분산하여 처리하기 위한 예제 프로젝트입니다.


## 구성
echo framework를 이용한 웹서버
분산예시로 redis connection pool 생성
RESTFul 예제로 음악의 곡정보를 응답하는 GET /track 구현 포함


## 특징
분산처리를 위한 goroutine worker를 미리 확보함

## 빌드
go build server.go


## 실행
./server -f config.json

