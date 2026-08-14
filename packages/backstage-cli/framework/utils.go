package framework

import "github.com/sqids/sqids-go"

// sqidAlphabet 乱序后的 0-9a-z 字库，作为常量不再变动
const sqidAlphabet = "2r8a59o7v1mipbufxd0eyc3hnk6gqltzs4jw"

var sqidInstance *sqids.Sqids

func init() {
	s, _ := sqids.New(sqids.Options{
		Alphabet:  sqidAlphabet,
		MinLength: 7,
	})
	sqidInstance = s
}

// EncodeId 将数字 ID 编码为 sqid 字符串，前缀加 "s"
func EncodeId(id int64) string {
	result, _ := sqidInstance.Encode([]uint64{uint64(id)})
	return "s" + result
}

// DecodeId 将 sqid 字符串解码为数字 ID，去掉第一个字母
func DecodeId(s string) int64 {
	if len(s) < 2 {
		return 0
	}
	result := sqidInstance.Decode(s[1:])
	if len(result) == 0 {
		return 0
	}
	return int64(result[0])
}
