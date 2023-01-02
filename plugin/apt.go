package plugin

import (
	"git.zabbix.com/ap/plugin-support/conf"
	"git.zabbix.com/ap/plugin-support/metric"
	"git.zabbix.com/ap/plugin-support/plugin"
	"github.com/go-co-op/gocron"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	pluginName  = "APT"
	keyUpdates  = "apt.updates"
	keySecurity = "apt.security"
)

type Options struct {
	plugin.SystemOptions `conf:"optional,name=System"`

	Interval int `conf:"optional,range=1:1440,default=1"`
}

type Plugin struct {
	plugin.Base
	updates   int
	security  int
	scheduler *gocron.Scheduler
	options   Options
}

var Impl Plugin

func (p *Plugin) Export(key string, _ []string, _ plugin.ContextProvider) (result interface{}, err error) {
	switch key {
	case keyUpdates:
		return p.updates, nil
	case keySecurity:
		return p.security, nil
	default:
		return nil, plugin.UnsupportedMetricError
	}
}

var updateMetrics = func(p *Plugin) {
	p.Debugf("updateMetrics")

	commands := map[string]string{
		keyUpdates:  "apt-get -s upgrade | grep 'upgraded,.*newly installed,' | cut -d ' ' -f1",
		keySecurity: "apt-get -s upgrade | grep 'standard security updates' | cut -d ' ' -f1",
	}

	for key, cmd := range commands {
		out, err := exec.Command("bash", "-c", cmd).Output()

		if err != nil {
			p.Errf("cannot execute %s: %s", key, err)
		}

		packages, err := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 32)

		switch key {
		case keyUpdates:
			p.updates = int(packages)
			break
		case keySecurity:
			p.security = int(packages)
			break
		}
	}
}

func (p *Plugin) Start() {
	_, _ = p.scheduler.Every(p.options.Interval).Minutes().StartImmediately().Do(updateMetrics, p)
	p.scheduler.StartAsync()
}

func (p *Plugin) Stop() {
	p.scheduler.Stop()
}

func (p *Plugin) Configure(_ *plugin.GlobalOptions, options interface{}) {
	if err := conf.Unmarshal(options, &p.options); err != nil {
		p.Errf("cannot unmarshal configuration options: %s", err)
	}
}

func (p *Plugin) Validate(options interface{}) error {
	var opts Options

	return conf.Unmarshal(options, &opts)
}

var metrics = metric.MetricSet{
	keyUpdates:  metric.New("Available Updates", []*metric.Param{}, false),
	keySecurity: metric.New("Security Updates", []*metric.Param{}, false),
}

func init() {
	Impl.updates = 0
	Impl.security = 0
	Impl.scheduler = gocron.NewScheduler(time.UTC)
	Impl.scheduler.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	plugin.RegisterMetrics(&Impl, pluginName, metrics.List()...)
}
