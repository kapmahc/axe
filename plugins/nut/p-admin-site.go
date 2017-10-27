package nut

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
	"github.com/spf13/viper"
)

func (p *AdminPlugin) getSiteStatus(l string, c *gin.Context) (interface{}, error) {
	data := gin.H{}

	var err error
	if data["os"], err = p._osStatus(); err != nil {
		return nil, err
	}
	if data["network"], err = p._networkStatus(); err != nil {
		return nil, err
	}
	data["jobber"] = p.Jobber.Status()
	data["routes"] = p._routes()

	if data["redis"], err = p._redisStatus(); err != nil {
		return nil, err
	}
	if data["postgresql"], err = p._dbStatus(); err != nil {
		return nil, err
	}

	return data, nil
}

func (p *AdminPlugin) _osStatus() (gin.H, error) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	hn, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	hu, err := user.Current()
	if err != nil {
		return nil, err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	var ifo syscall.Sysinfo_t
	if err := syscall.Sysinfo(&ifo); err != nil {
		return nil, err
	}
	return gin.H{
		"app author":           fmt.Sprintf("%s <%s>", web.AuthorName, web.AuthorEmail),
		"app licence":          web.Copyright,
		"app version":          fmt.Sprintf("%s(%s) - %s", web.Version, web.BuildTime, viper.GetString("env")),
		"app root":             pwd,
		"who-am-i":             fmt.Sprintf("%s@%s", hu.Username, hn),
		"go version":           runtime.Version(),
		"go root":              runtime.GOROOT(),
		"go runtime":           runtime.NumGoroutine(),
		"go last gc":           time.Unix(0, int64(mem.LastGC)).Format(time.ANSIC),
		"os cpu":               runtime.NumCPU(),
		"os ram(free/total)":   fmt.Sprintf("%dM/%dM", ifo.Freeram/1024/1024, ifo.Totalram/1024/1024),
		"os swap(free/total)":  fmt.Sprintf("%dM/%dM", ifo.Freeswap/1024/1024, ifo.Totalswap/1024/1024),
		"go memory(alloc/sys)": fmt.Sprintf("%dM/%dM", mem.Alloc/1024/1024, mem.Sys/1024/1024),
		"os time":              time.Now().Format(time.ANSIC),
		"os arch":              fmt.Sprintf("%s(%s)", runtime.GOOS, runtime.GOARCH),
		"os uptime":            (time.Duration(ifo.Uptime) * time.Second).String(),
		"os loads":             ifo.Loads,
		"os procs":             ifo.Procs,
	}, nil
}
func (p *AdminPlugin) _networkStatus() (gin.H, error) {
	sts := gin.H{}
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, v := range ifs {
		ips := []string{v.HardwareAddr.String()}
		adrs, err := v.Addrs()
		if err != nil {
			return nil, err
		}
		for _, adr := range adrs {
			ips = append(ips, adr.String())
		}
		sts[v.Name] = ips
	}
	return sts, nil
}

func (p *AdminPlugin) _dbStatus() (gin.H, error) {
	args := viper.GetStringMap("postgresql")

	var ver string
	if _, err := p.DB.Query(pg.Scan(&ver), "SELECT VERSION()"); err != nil {
		return nil, err
	}
	// http://blog.javachen.com/2014/04/07/some-metrics-in-postgresql.html
	var size string
	if _, err := p.DB.Query(pg.Scan(&size), "select pg_size_pretty(pg_database_size('postgres'))"); err != nil {
		return nil, err
	}

	// sts, err:=p.DB.Exec("select pid,current_timestamp - least(query_start,xact_start) AS runtime,substr(query,1,25) AS current_query from pg_stat_activity where not pid=pg_backend_pid()")
	// if err!=nil{
	// 	return nil, err
	// }
	// val[fmt.Sprintf("pid-%d", pid)] = fmt.Sprintf("%s (%v)", ts.Format("15:04:05.999999"), qry)

	return gin.H{
		"url": fmt.Sprintf(
			"postgres://%s@%s:%d/%s",
			args["user"],
			args["host"],
			args["port"],
			args["dbname"],
		),
		"version": ver,
		"size":    size,
	}, nil
}

func (p *AdminPlugin) _routes() []gin.H {
	var items []gin.H
	for _, r := range p.Router.Routes() {
		items = append(items, gin.H{"method": r.Method, "path": r.Path})
	}
	return items
}

func (p *AdminPlugin) _redisStatus() ([]string, error) {
	c := p.Redis.Get()
	defer c.Close()
	info, err := redis.String(c.Do("INFO"))
	if err != nil {
		return nil, err
	}
	return strings.Split(info, "\n"), nil
}