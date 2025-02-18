package filepath

import (
	cnt "VirtualRegistryManagement/constants"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
)

var provider = map[string]filepathInterface{}

type filepathInterface interface {
	Path(uri *url.URL, attr attribute) string
	Prefix(args ...interface{}) string
	Suffix(args ...interface{}) string
}

type attribute struct {
	Prefix []interface{}
	Suffix []interface{}
}

func (attr attribute) path(path string) (string, error) {
	uri, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	if val, exist := provider[uri.Scheme]; !exist {
		return "", tkErr.New(cnt.FilepathTypeIsNotSupportedErr)
	} else {
		return val.Path(uri, attr), nil
	}
}

const (
	_localType = "local"
)

func Init(config interface{}) (err error) {
	for _, volume := range cast.ToSlice(config) {
		if volume == nil {
			continue
		}
		cfg := viper.New()
		cfg.MergeConfigMap(cast.ToStringMap(volume))

		if !cfg.IsSet("type") || !cfg.IsSet("scheme") {
			return tkErr.New(cnt.FilepathTypeAndSchemeIsRequiredErr)
		}

		typ := cfg.GetString("type")
		scheme := cfg.GetString("scheme")

		if _, exist := provider[scheme]; exist {
			return tkErr.New(cnt.FilepathSchemeIsRepeatedErr)
		}

		switch typ {
		case _localType:
			provider[scheme] = UnmarshalLocalVolume(cfg)
		default:
			return tkErr.New(cnt.FilepathTypeIsNotSupportedErr)
		}
	}
	return nil
}

func Path(path string, prefix, suffix []interface{}) (string, error) {
	attr := attribute{
		Prefix: prefix,
		Suffix: suffix,
	}
	return attr.path(path)
}

func Validate(path string, prefix, suffix []interface{}, defaultExtension, defaultFilename string) (string, error) {
	isDir := false
	if len(path) != 0 && path[len(path)-1] == '/' {
		isDir = true
	}
	p, _ := url.Parse(path)
	pp := filepath.Join(p.Host, p.Path)

	var ppp string
	if isDir {
		ppp = fmt.Sprintf("%s://%s/%s.%s", p.Scheme, pp, defaultFilename, defaultExtension)
	} else if filepath.Ext(pp) == "" {
		ppp = fmt.Sprintf("%s://%s.%s", p.Scheme, pp, defaultExtension)
	} else {
		ppp = fmt.Sprintf("%s://%s", p.Scheme, pp)
	}

	pppp, err := Path(ppp, prefix, suffix)
	if err != nil {
		return "", err
	}

	{
		if _, err := os.Stat(pppp); err == nil {
			return "", tkErr.New(cnt.FilepathFilepathIsExistErr)
		}
	}

	dir := filepath.Dir(pppp)
	file, err := os.Stat(dir)
	if err != nil {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return "", tkErr.New(cnt.FilepathFilepathIsExistErr)
		}
		file, _ = os.Stat(dir)
	}
	if !file.IsDir() {
		return "", tkErr.New(cnt.FilepathFilepathIsExistErr)
	}

	return ppp, nil
}
