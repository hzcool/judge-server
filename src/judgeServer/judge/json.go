package judge

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

//将v转化成json字符串
func encodeJson(v interface{}) ([]byte, error){
    var buf bytes.Buffer
    encoder := json.NewEncoder(&buf)
    encoder.SetEscapeHTML(false)
    if err := encoder.Encode(v); err != nil{
        return nil, err
    }
    return buf.Bytes(), nil
}

//解析请求body中的json数据
func decodeJson(body io.Reader) (map[string]interface{},error) {
	var res  map[string]interface{}
	err := json.NewDecoder(body).Decode(&res)
	if err != nil {
		return nil,err
	}
	return res,err
}

func responseJson(w http.ResponseWriter, v interface{}) {
	bytes,_ := encodeJson(v)
	io.WriteString(w, string(bytes))
}




