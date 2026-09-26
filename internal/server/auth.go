package server

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"context"
	"net/http"
	"strings"
	"sync"
	"time"
)

type AuthConfig struct { ViewerToken string; OperatorToken string; TrustProxy bool }
type authContextKey string
const authRoleKey authContextKey="role"
const authMethodKey authContextKey="method"
func withAuthContext(r *http.Request,role,method string)*http.Request{ctx:=context.WithValue(r.Context(),authRoleKey,role);ctx=context.WithValue(ctx,authMethodKey,method);return r.WithContext(ctx)}
func authenticatedRole(r *http.Request)string{v,_:=r.Context().Value(authRoleKey).(string);return v}
func authenticatedMethod(r *http.Request)string{v,_:=r.Context().Value(authMethodKey).(string);return v}
type session struct { role string; csrf string; created time.Time; expires time.Time }
type sessionStore struct { mu sync.Mutex; values map[string]session }

func newSessionStore()*sessionStore{return &sessionStore{values:map[string]session{}}}
func(c AuthConfig) enabled()bool{return c.ViewerToken!=""||c.OperatorToken!=""}
func(c AuthConfig) role(token string)string{if c.OperatorToken!=""&&secureEqual(token,c.OperatorToken){return "operator"};if c.ViewerToken!=""&&secureEqual(token,c.ViewerToken){return "viewer"};return ""}
func secureEqual(a,b string)bool{if len(a)!=len(b){return false};return subtle.ConstantTimeCompare([]byte(a),[]byte(b))==1}
func bearerToken(r *http.Request)string{v:=r.Header.Get("Authorization");if !strings.HasPrefix(v,"Bearer "){return ""};return strings.TrimSpace(strings.TrimPrefix(v,"Bearer "))}
func randomToken()(string,error){b:=make([]byte,32);if _,e:=rand.Read(b);e!=nil{return "",e};return base64.RawURLEncoding.EncodeToString(b),nil}
func(s *sessionStore)create(role string)(string,error){id,e:=randomToken();if e!=nil{return "",e};csrf,e:=randomToken();if e!=nil{return "",e};s.mu.Lock();now:=time.Now();s.values[id]=session{role:role,csrf:csrf,created:now,expires:now.Add(12*time.Hour)};s.mu.Unlock();return id,nil}
func(s *sessionStore)get(id string)(session,bool){s.mu.Lock();defer s.mu.Unlock();v,ok:=s.values[id];if !ok{return session{},false};if time.Now().After(v.expires){delete(s.values,id);return session{},false};return v,true}
func(s *sessionStore)role(id string)string{v,ok:=s.get(id);if !ok{return ""};return v.role}
func(s *sessionStore)delete(id string){s.mu.Lock();delete(s.values,id);s.mu.Unlock()}
func requestAuth(c AuthConfig,s *sessionStore,r *http.Request)(string,string){if role:=c.role(bearerToken(r));role!=""{return role,"bearer"};if ck,e:=r.Cookie("botweb_session");e==nil{if role:=s.role(ck.Value);role!=""{return role,"session"}};return "",""}
func requestRole(c AuthConfig,s *sessionStore,r *http.Request)string{role,_:=requestAuth(c,s,r);return role}
func sessionWriteOK(c AuthConfig,s *sessionStore,r *http.Request)bool{if c.role(bearerToken(r))!=""{return true};ck,e:=r.Cookie("botweb_session");if e!=nil{return false};v,ok:=s.get(ck.Value);return ok&&secureEqual(r.Header.Get("X-CSRF-Token"),v.csrf)}
func requestHTTPS(c AuthConfig,r *http.Request)bool{if r.TLS!=nil{return true};return c.TrustProxy&&strings.EqualFold(strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-Proto"),",")[0]),"https")}
func authMiddleware(c AuthConfig,s *sessionStore,next http.Handler)http.Handler{if !c.enabled(){return next};return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){if r.URL.Path=="/healthz"||r.URL.Path=="/login"||r.URL.Path=="/app.css"{next.ServeHTTP(w,r);return};role,method:=requestAuth(c,s,r);if role==""{if strings.HasPrefix(r.URL.Path,"/api/"){w.Header().Set("WWW-Authenticate","Bearer");http.Error(w,"authentication required",http.StatusUnauthorized);return};http.Redirect(w,r,"/login",http.StatusSeeOther);return};r=withAuthContext(r,role,method);write:=r.Method!=http.MethodGet&&r.Method!=http.MethodHead;if write&&r.URL.Path=="/logout"{if method=="session"&&!sessionWriteOK(c,s,r){http.Error(w,"CSRF validation failed",http.StatusForbidden);return};next.ServeHTTP(w,r);return};if write&&role!="operator"{http.Error(w,"operator role required",http.StatusForbidden);return};if write&&!sessionWriteOK(c,s,r){http.Error(w,"CSRF validation failed",http.StatusForbidden);return};next.ServeHTTP(w,r)})}

func(s *sessionStore)rotate(id string)(string,error){v,ok:=s.get(id);if !ok{return "",http.ErrNoCookie};next,e:=randomToken();if e!=nil{return "",e};csrf,e:=randomToken();if e!=nil{return "",e};s.mu.Lock();delete(s.values,id);v.csrf=csrf;v.expires=time.Now().Add(12*time.Hour);s.values[next]=v;s.mu.Unlock();return next,nil}
