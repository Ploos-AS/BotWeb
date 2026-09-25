package pbmp
import("bufio";"encoding/json";"fmt";"net";"time")
type Client struct{Endpoint string;Timeout time.Duration}
type response struct{PBMP int `json:"pbmp"`;Type string `json:"type"`;ID string `json:"id"`;OK bool `json:"ok"`;Result json.RawMessage `json:"result"`;Error any `json:"error"`}
func(c Client)Call(method string)(json.RawMessage,error){d:=c.Timeout;if d==0{d=2*time.Second};x,err:=net.DialTimeout("unix",c.Endpoint,d);if err!=nil{return nil,err};defer x.Close();_ = x.SetDeadline(time.Now().Add(d));req:=map[string]any{"pbmp":1,"type":"request","id":"botweb","method":method,"params":map[string]any{}};if err=json.NewEncoder(x).Encode(req);err!=nil{return nil,err};var r response;if err=json.NewDecoder(bufio.NewReader(x)).Decode(&r);err!=nil{return nil,err};if r.PBMP!=1||r.Type!="response"||r.ID!="botweb"{return nil,fmt.Errorf("invalid PBMP response")};if !r.OK{return nil,fmt.Errorf("PBMP request failed")};return r.Result,nil}
