package common


//  utils
import (
	"strconv"
	"crypto/md5"
	"encoding/hex"
	"net"
	"strings"
	"fmt"
)

const (
	OP_TYPE_NONE = iota
	OP_TYPE_ADD
	OP_TYPE_START
	OP_TYPE_STOP
	OP_TYPE_RESTART
	OP_TYPE_REWRITE
	OP_TYPE_MAX
)

const (
	JOB_STATE_NONE = iota
	JOB_STATE_TODO
	JOB_STATE_DOING
	JOB_STATE_TORESTART
	JOB_STATE_TOSTOP
	JOB_STATE_STOP
	JOB_STATE_DONE
)

func PanicIf(err error) {

	if err != nil {
		fmt.Println("error happend")
		//panic(err)
	}
}

func ParseInt(value string) int {

	if value == "" {
		return 0
	}

	val, _ := strconv.Atoi(value)
	return val
}

func IntString(value int) string {

	return strconv.Itoa(value)
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GetLocalAddr() string {
	conn, err := net.Dial("udp", "localhost:80")
	if err != nil {
		return ""
	}

	defer conn.Close()

	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

func GetIPFromUrl(url string) string {

	url = strings.Split(url, "//")[1]
	url = strings.Split(url, "/")[0]
	if strings.Contains(url, ":") {

		return strings.Split(url, ":")[0]
	}

	return url

}