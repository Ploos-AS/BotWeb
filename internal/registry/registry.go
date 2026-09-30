package registry
import("encoding/json";"os")
type Bot struct{ID string `json:"id"`;Name string `json:"name"`;Profile string `json:"profile,omitempty"`;Transport string `json:"transport"`;Endpoint string `json:"endpoint"`}
type Registry struct{Bots []Bot `json:"bots"`}
func Load(path string)(Registry,error){var r Registry;b,err:=os.ReadFile(path);if err!=nil{return r,err};err=json.Unmarshal(b,&r);if err!=nil{return r,err};for i:=range r.Bots{if r.Bots[i].Profile==""{r.Bots[i].Profile="bot-m0"}};return r,nil}
func(r Registry)Find(id string)(Bot,bool){for _,b:=range r.Bots{if b.ID==id{return b,true}};return Bot{},false}
