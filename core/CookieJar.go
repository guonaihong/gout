package core

import (
	"errors"
	"fmt"
	"github.com/go-basic/uuid"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type TCookieJar struct {
	mu sync.Mutex
	db *gorm.DB
}
type Entery struct {
	gorm.Model_UUID
	JarKey string `yaml:"jar_key"`
	Name       string `yaml:"name"`
	Value      string `yaml:"value"`
	Domain     string `yaml:"domain"`
	Path       string `yaml:"path"`
	SameSite   string `yaml:"same_site"`
	Secure     bool `yaml:"secure"`
	HttpOnly   bool `yaml:"http_only"`
	Persistent bool `yaml:"persistent"`
	HostOnly   bool `yaml:"host_only"`
	Expires    time.Time `yaml:"expires"`
	//Creation   time.Time `yaml:"creation"`
	//LastAccess time.Time `yaml:"last_access"`
}

func (e *Entery) BeforeCreate(tx *gorm.DB) error {
	e.ID=uuid.New()
	return nil
}



func (jar *TCookieJar) InitCookieJarWithSqlite(db *gorm.DB) (err error) {
	jar.db=db
	jar.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Entery{})
	return
}
func (jar TCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie)  {
	for _, cookie := range cookies {
		e,err:= jar.newEntry(cookie)
		if err!=nil{
			fmt.Println(err.Error())
		}
		e.JarKey=fmt.Sprintf("%s://%s",u.Scheme,u.Host)
		var oldCookie Entery
		jar.db.Delete(&oldCookie,Entery{
			JarKey: e.JarKey,
			Name: e.Name,
		})
		jar.db.Save(&e)
	}
}

func (jar TCookieJar) Cookies(u *url.URL) (cookies []*http.Cookie)  {
	var enteries []Entery
	jar.db.Find(&enteries,Entery{
		JarKey: fmt.Sprintf("%s://%s",u.Scheme,u.Host),
	})
	for _,entery:=range enteries{
		if !strings.HasPrefix(u.Path,entery.Path){
			continue
		}
		fmt.Println(u.Host)
		cookies = append(cookies, &http.Cookie{
			Name: entery.Name,
			Value: entery.Value,
			Domain:entery.Domain})
	}
	return
}


var endOfTime = time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)
func (jar *TCookieJar) newEntry(c *http.Cookie) (e Entery, err error) {
	now:=time.Now()
	e.Name = c.Name
	if c.Path == "" || c.Path[0] != '/' {
		e.Path = "/"
	} else {
		e.Path = c.Path
	}
	e.Domain=c.Domain
	if err != nil {
		return e, err
	}
	if c.MaxAge < 0 {
		return e, nil
	} else if c.MaxAge > 0 {
		e.Expires = now.Add(time.Duration(c.MaxAge) * time.Second)
		e.Persistent = true
	} else {
		if c.Expires.IsZero() {
			e.Expires = endOfTime
			e.Persistent = false
		} else {
			if !c.Expires.After(now) {
				return e, nil
			}
			e.Expires = c.Expires
			e.Persistent = true
		}
	}

	e.Value = c.Value
	e.Secure = c.Secure
	e.HttpOnly = c.HttpOnly

	switch c.SameSite {
	case http.SameSiteDefaultMode:
		e.SameSite = "SameSite"
	case http.SameSiteStrictMode:
		e.SameSite = "SameSite=Strict"
	case http.SameSiteLaxMode:
		e.SameSite = "SameSite=Lax"
	}

	return e, nil
}

var (
	errIllegalDomain   = errors.New("MyCookieJar: illegal cookie domain attribute")
	errMalformedDomain = errors.New("MyCookieJar: malformed cookie domain attribute")
	errNoHostname      = errors.New("MyCookieJar: no host name available (IP only)")
)

