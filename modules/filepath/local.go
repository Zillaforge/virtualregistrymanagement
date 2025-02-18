package filepath

import (
	"fmt"
	"net/url"

	"github.com/spf13/viper"
)

// Local is <scheme>://<path> => <root>/<prefix>/<path>/<suffix>
type Local struct {
	name   string
	scheme string
	root   string
	prefix string
	suffix string
}

func UnmarshalLocalVolume(svcCfg *viper.Viper) *Local {
	return &Local{
		name:   svcCfg.GetString("name"),
		scheme: svcCfg.GetString("scheme"),
		root:   svcCfg.GetString("root"),
		prefix: svcCfg.GetString("prefix"),
		suffix: svcCfg.GetString("suffix"),
	}
}

func (l *Local) Path(uri *url.URL, attr attribute) string {
	r, _ := url.JoinPath(l.root, l.Prefix(attr.Prefix...), uri.Host, uri.Path, l.Suffix(attr.Suffix...))
	p, _ := url.PathUnescape(r)
	return p
}

func (l *Local) Prefix(args ...interface{}) string {
	return fmt.Sprintf(l.prefix, args...)
}

func (l *Local) Suffix(args ...interface{}) string {
	return fmt.Sprintf(l.suffix, args...)
}
