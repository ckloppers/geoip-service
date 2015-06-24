package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gocraft/web"
	"github.com/oschwald/geoip2-golang"
)

type Context struct {
	db *geoip2.Reader
}

func (ctx *Context) OpenMaxMindDB(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	db, err := geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx.db = db
	next(rw, req)
}

func (ctx *Context) LookUpCountryForIp(rw web.ResponseWriter, req *web.Request) {

	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP(req.PathParams["ipstring"])
	record, err := ctx.db.Country(ip)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf(rw, "IP: ", req.PathParams["ipstring"], "ISO country code: %v\n", record.Country.IsoCode)
	fmt.Fprint(rw, "IP: ", req.PathParams["ipstring"], "  ISO country code: ", record.Country.IsoCode)

}

func main() {

	router := web.New(Context{}).
		Middleware(web.LoggerMiddleware).
		Middleware((*Context).OpenMaxMindDB).
		Get("/ccfip/:ipstring", (*Context).LookUpCountryForIp)

	http.ListenAndServe("localhost:3000", router) // Start the server!
}
