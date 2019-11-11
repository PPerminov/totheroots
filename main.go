package main

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

/* address must be:
1. Protocol is ALWAYS https
2. address blocks divided by dots
3. Last two blocks is domain cai.dev
4. third from the end - cluster name
5. if first block can be converted to integer then I suppose that it is a port to connect
6. Example: 6000.jenkins.doc-cluster.cai.dev will redirect to port 6000 of internal kubernetes service "jenkins" in default namespace
7. Example: jenkins-hcydhs-jhucydjs.pod.jenkins-ns.doc-cluster.cai.dev will redirect to port 80 of pod jenkins-hcydhs-jhucydjs in jenkins-ns namespace.
8. Full address template: port.name.type.namespace.clustername.cai.dev Only clustername and cai.dev are mandatory fields. If everything else is ommited then it will try to connect to 80.default.svc.default.clustername.cai.dev
9. Address recognition going left -> right. So if you want to point a name then you can not to point anything else. But for pod and namespace you need to input name and name.pod respectively
*/

type Destination struct {
	resultedOne []string
	port        int
	logger      *logrus.Logger
	uuid        uuid.UUID
	path        string
}

func (d *Destination) getTarget(req *http.Request) {
	var start, port, index int
	var internalPath []string
	var item string
	var err error

	d.port = 80
	d.resultedOne = []string{"default", "svc", "default"}
	d.d(d)
	host := strings.Split(req.Host, ":")
	d.d(host)
	internalPath = strings.Split(host[0], ".")
	port, err = strconv.Atoi(internalPath[0])
	if err != nil {
		start = 0
	} else {
		d.port = port
		start = 1
	}
	internalPath = internalPath[start : len(internalPath)-3]

	for index, item = range internalPath {
		d.resultedOne[index] = item
	}
	d.resultedOne[1], d.resultedOne[2] = d.resultedOne[2], d.resultedOne[1]
	d.path = "http://" + strings.Join(d.resultedOne, ".") + ":" + strconv.Itoa(d.port)
}

func (d *Destination) serveReverseProxy(res http.ResponseWriter, req *http.Request) {
	d.uuid = uuid.NewV4()
	d.d(d.uuid)
	d.getTarget(req)
	d.d(d)
	url2, err := url.Parse(d.path)
	d.d(url2)
	if err != nil {
		d.logger.Error(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url2)
	req.URL.Host = url2.Host
	req.URL.Scheme = url2.Scheme
	proxy.ServeHTTP(res, req)
}

func (m *mainStruct) handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	runner := Destination{}
	runner.logger = m.logger
	runner.serveReverseProxy(res, req)
}

func (d *Destination) d(a interface{}) {
	d.logger.Debug(d.uuid, a)
}

type mainStruct struct {
	logger *logrus.Logger
}

func main() {
	var runner mainStruct
	var listenTo string = "0.0.0.0:8080"
	runner.logger = logrus.New()
	if os.Getenv("DEBUG") != "" {
		runner.logger.Level = logrus.DebugLevel
	}
	// temporary
	runner.logger.Level = logrus.DebugLevel
	http.HandleFunc("/", runner.handleRequestAndRedirect)
	if os.Getenv("CAIWAY_LISTEN") != "" {
		listenTo = os.Getenv("CAIWAY_LISTEN")
	}
	if err := http.ListenAndServe(listenTo, nil); err != nil {
		runner.logger.Panic(err)
	}
}
