package internal

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	box "github.com/lukasjoc/fritz/internal/go-fritzbox"
	"gopkg.in/yaml.v2"
)

const (
	FritzTcpPort      = 49000
	DefaultConfigPath = ".config/fritz/config.yml"
)

type Config struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Fritz struct {
	Client *box.Client
	config *Config
}

func NewFritz() (*Fritz, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(filepath.Join(usr.HomeDir, DefaultConfigPath))
	if err != nil {
		return nil, err
	}
	var config Config
	if err = yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	configBaseURL, err := url.Parse(config.Host)
	if err != nil {
		return nil, err
	}
	client := box.NewClient(nil)
	client.BaseURL = configBaseURL
	return &Fritz{client, &config}, nil
}

func (f *Fritz) Connect() error {
	if err := f.Client.Auth(f.config.Username, f.config.Password); err != nil {
		return err
	}
	return nil
}

func (f *Fritz) Info() error {
	if err := f.Connect(); err != nil {
		return err
	}
	defer f.Client.Session.Close()
	req, _ := f.Client.NewRequest("POST", "data.lua", url.Values{
		"xhr":         {"1"},
		"sid":         {f.Client.Session.Sid},
		"lang":        {"de"},
		"page":        {"overview"},
		"xhrId":       {"all"},
		"useajax":     {"1"},
		"no_sidrenew": {""},
	})
	info := struct {
		Data any `json:"data"`
	}{}
	fmt.Println("Getting info...")
	_, err := f.Client.Do(req, &info)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(info.Data, " ", "  ")
	if err != nil {
		return nil
	}
	fmt.Println(string(data))
	return nil
}
func (f *Fritz) Reconnect() error {
	if err := f.Connect(); err != nil {
		return err
	}
	defer f.Client.Session.Close()

	req, _ := f.Client.NewRequest(
		"GET",
		"internet/inetstat_monitor.lua?myXhr=1&action=disconnect&useajax=1&xhr=1&t1695669022799=nocache",
		url.Values{})
	req1, _ := f.Client.NewRequest(
		"GET",
		"internet/inetstat_monitor.lua?myXhr=1&action=connect&useajax=1&xhr=1&t1695669022799=nocache",
		url.Values{})

	fmt.Println("Disconnecting...")
	if _, err := f.Client.Do(req, nil); err != nil {
		return err
	}

	fmt.Println("Connecting...")
	if _, err := f.Client.Do(req1, nil); err != nil {
		return err
	}
	fmt.Println("OK. This can take up to 30s to take full effect..")
	return nil
}

func (f *Fritz) Reboot() error {
	if err := f.Connect(); err != nil {
		return err
	}
	defer f.Client.Session.Close()
	req, _ := f.Client.NewRequest(
		"POST",
		"data.lua",
		url.Values{
			"sid":    {f.Client.Session.Sid},
			"xhr":    {"1"},
			"page":   {"reboot"},
			"reboot": {"0"},
		})
	if _, err := f.Client.Do(req, nil); err != nil {
		return err
	}
	req1, _ := f.Client.NewRequest(
		"POST",
		"reboot.lua",
		url.Values{
			"ajax":        {"1"},
			"sid":         {f.Client.Session.Sid},
			"no_sidrenew": {"1"},
			"xhr":         {"1"},
			"useajax":     {"1"},
		})
	_, err := f.Client.Do(req1, nil)
	if err != nil {
		return err
	}
	fmt.Println("Rebooting... This can take a while. (5-8 minutes)")
	return nil
}
