package server
import("os";"path/filepath";"net/http";"net/http/httptest";"strings";"testing")
func testServer(t *testing.T)*Server{t.Helper();p:=filepath.Join(t.TempDir(),"bots.json");if err:=os.WriteFile(p,[]byte(`{"bots":[{"id":"x","name":"Example","transport":"unix","endpoint":"/secret/path.sock"}]}`),0600);err!=nil{t.Fatal(err)};s,e:=New(p);if e!=nil{t.Fatal(e)};return s}
func TestIndexAndHeaders(t *testing.T){r:=httptest.NewRequest("GET","/",nil);w:=httptest.NewRecorder();testServer(t).Handler().ServeHTTP(w,r);if w.Code!=200{t.Fatal(w.Code)};if !strings.Contains(w.Body.String(),"PBMP fleet console"){t.Fatal("missing UI")};if w.Header().Get("Content-Security-Policy")==""{t.Fatal("missing CSP")}}
func TestRegistryDoesNotLeakEndpoint(t *testing.T){r:=httptest.NewRequest("GET","/api/v1/bots",nil);w:=httptest.NewRecorder();testServer(t).Handler().ServeHTTP(w,r);if strings.Contains(w.Body.String(),"secret/path"){t.Fatal("endpoint leaked")}}

func TestBotRouteServesAppShell(t *testing.T){r:=httptest.NewRequest("GET","/bots/x",nil);w:=httptest.NewRecorder();testServer(t).Handler().ServeHTTP(w,r);if w.Code!=200{t.Fatal(w.Code)};if !strings.Contains(w.Body.String(),"PBMP fleet console"){t.Fatal("missing app shell")}}
func TestCapabilityNavigationIsGeneric(t *testing.T){r:=httptest.NewRequest("GET","/app.js",nil);w:=httptest.NewRecorder();testServer(t).Handler().ServeHTTP(w,r);body:=w.Body.String();if !strings.Contains(body,"channels.list")||!strings.Contains(body,"modules.list"){t.Fatal("missing capability navigation")};if strings.Contains(body,"implementation==='luca'")||strings.Contains(body,"implementation==='engo'"){t.Fatal("implementation-specific UI branch")}}

func TestWriteRequiresJSON(t *testing.T){r:=httptest.NewRequest("POST","/api/v1/bots/x/channels/join",strings.NewReader("{}"));w:=httptest.NewRecorder();testServer(t).Handler().ServeHTTP(w,r);if w.Code!=http.StatusForbidden{t.Fatalf("code=%d",w.Code)}}
func TestWriteRejectsCrossOrigin(t *testing.T){r:=httptest.NewRequest("POST","http://botweb.local/api/v1/bots/x/channels/join",strings.NewReader("{}"));r.Host="botweb.local";r.Header.Set("Content-Type","application/json");r.Header.Set("Origin","https://evil.example");w:=httptest.NewRecorder();testServer(t).Handler().ServeHTTP(w,r);if w.Code!=http.StatusForbidden{t.Fatalf("code=%d",w.Code)}}


func TestBotAIUIIsCapabilityDriven(t *testing.T){
	r:=httptest.NewRequest("GET","/app.js",nil);w:=httptest.NewRecorder();testServer(t).Handler().ServeHTTP(w,r);body:=w.Body.String()
	if !strings.Contains(body,"botai.status")||!strings.Contains(body,"/botai"){t.Fatal("missing BotAI capability view")}
	if strings.Contains(body,"implementation==='engo'"){t.Fatal("Engo-specific BotAI UI branch")}
}


func TestAuthViewerCanRead(t *testing.T){s:=testServer(t);s.auth=AuthConfig{ViewerToken:"view",OperatorToken:"operate"};r:=httptest.NewRequest("GET","/",nil);r.Header.Set("Authorization","Bearer view");w:=httptest.NewRecorder();s.Handler().ServeHTTP(w,r);if w.Code!=http.StatusOK{t.Fatalf("code=%d",w.Code)}}
func TestAuthRejectsMissingToken(t *testing.T){s:=testServer(t);s.auth=AuthConfig{ViewerToken:"view"};r:=httptest.NewRequest("GET","/",nil);w:=httptest.NewRecorder();s.Handler().ServeHTTP(w,r);if w.Code!=http.StatusUnauthorized{t.Fatalf("code=%d",w.Code)}}
func TestViewerCannotWrite(t *testing.T){s:=testServer(t);s.auth=AuthConfig{ViewerToken:"view",OperatorToken:"operate"};r:=httptest.NewRequest("POST","/api/v1/bots/x/channels/join",strings.NewReader("{}"));r.Header.Set("Authorization","Bearer view");r.Header.Set("Content-Type","application/json");w:=httptest.NewRecorder();s.Handler().ServeHTTP(w,r);if w.Code!=http.StatusForbidden{t.Fatalf("code=%d",w.Code)}}
func TestHealthDoesNotRequireAuth(t *testing.T){s:=testServer(t);s.auth=AuthConfig{ViewerToken:"view"};r:=httptest.NewRequest("GET","/healthz",nil);w:=httptest.NewRecorder();s.Handler().ServeHTTP(w,r);if w.Code!=http.StatusOK{t.Fatalf("code=%d",w.Code)}}


func TestBrowserAuthRedirectsToLogin(t *testing.T){s:=testServer(t);s.auth=AuthConfig{ViewerToken:"view"};r:=httptest.NewRequest("GET","/",nil);w:=httptest.NewRecorder();s.Handler().ServeHTTP(w,r);if w.Code!=http.StatusSeeOther||w.Header().Get("Location")!="/login"{t.Fatalf("code=%d location=%q",w.Code,w.Header().Get("Location"))}}
func TestBrowserLoginCreatesSession(t *testing.T){s:=testServer(t);s.auth=AuthConfig{ViewerToken:"view"};r:=httptest.NewRequest("POST","/login",strings.NewReader("token=view"));r.Header.Set("Content-Type","application/x-www-form-urlencoded");w:=httptest.NewRecorder();s.Handler().ServeHTTP(w,r);if w.Code!=http.StatusSeeOther{t.Fatalf("code=%d",w.Code)};res:=w.Result();var ck *http.Cookie;for _,c:=range res.Cookies(){if c.Name=="botweb_session"{ck=c}};if ck==nil||!ck.HttpOnly||ck.SameSite!=http.SameSiteStrictMode{t.Fatal("missing hardened session cookie")};q:=httptest.NewRequest("GET","/",nil);q.AddCookie(ck);v:=httptest.NewRecorder();s.Handler().ServeHTTP(v,q);if v.Code!=http.StatusOK{t.Fatalf("session code=%d",v.Code)}}
func TestBrowserLoginRejectsBadToken(t *testing.T){s:=testServer(t);s.auth=AuthConfig{ViewerToken:"view"};r:=httptest.NewRequest("POST","/login",strings.NewReader("token=wrong"));r.Header.Set("Content-Type","application/x-www-form-urlencoded");w:=httptest.NewRecorder();s.Handler().ServeHTTP(w,r);if w.Code!=http.StatusUnauthorized{t.Fatalf("code=%d",w.Code)}}
