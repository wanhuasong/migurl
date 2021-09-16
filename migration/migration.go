package migration

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/juju/errors"
	"github.com/wanhuasong/migurl/config"
	"github.com/wanhuasong/migurl/utils"
)

type Migration struct {
	InstanceConfig        *config.Config
	ApiBaseURL            string
	QiniuPublicDomain     string
	GenericStorageBaseURL string

	ContainerID string
	Client      *client.Client
}

func NewMigration(cfgFile, apiBaseURL, qiniuPublicDomain, genericStorageBaseURL string) (m *Migration, err error) {
	if cfgFile == "" {
		err = errors.New("Invalid config file")
		return
	}
	var cfg *config.Config
	cfg, err = config.LoadConfig(cfgFile)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	cli, err := client.NewEnvClient()
	if err != nil {
		err = errors.Trace(err)
		return
	}
	m = &Migration{
		InstanceConfig:        cfg,
		ApiBaseURL:            apiBaseURL,
		QiniuPublicDomain:     qiniuPublicDomain,
		GenericStorageBaseURL: genericStorageBaseURL,
		Client:                cli,
	}
	return
}

func (m *Migration) Do() error {
	if err := m.InitContainerID(); err != nil {
		return errors.Trace(err)
	}

	var baseURL string
	if m.ApiBaseURL != "" {
		baseURL = m.ApiBaseURL
	} else if m.QiniuPublicDomain != "" {
		baseURL = m.QiniuPublicDomain
	} else if m.GenericStorageBaseURL != "" {
		baseURL = m.GenericStorageBaseURL
	} else {
		return errors.New("Must specify apiBaseURL, qiniuPublicDomain or genericStorageBaseURL")
	}
	sqls := []string{
		fmt.Sprintf("UPDATE `user` SET `avatar`=REPLACE(avatar, '%s', '');", baseURL),
		fmt.Sprintf("UPDATE `team` SET `logo`=REPLACE(`logo`, '%s', '');", baseURL),
		fmt.Sprintf("UPDATE `organization` SET `logo`=REPLACE(`logo`, '%s', ''), `favicon`=REPLACE(`favicon`, '%s', '');", baseURL, baseURL),
	}
	for _, sql := range sqls {
		if err := m.execSql(sql); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func (m *Migration) execSql(sql string) error {
	command := []string{
		"mysql",
		"-h" + m.InstanceConfig.MySQLHost,
		"-P" + m.InstanceConfig.MySQLPort,
		"-u" + m.InstanceConfig.MySQLUserName,
		"-p" + m.InstanceConfig.MySQLPassword,
		"-D" + m.InstanceConfig.ProjectDBName,
		"-e" + sql,
	}
	err := utils.ExecInContainer(m.Client, m.ContainerID, command)
	return errors.Trace(err)
}

func (m *Migration) InitContainerID() (err error) {
	containers, err := utils.ListContainers(m.Client, false)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	for _, container := range containers {
		var count int64
		for _, port := range container.Ports {
			if port.PublicPort == m.InstanceConfig.Port || port.PublicPort == m.InstanceConfig.HTTPSPort {
				count++
			}
		}
		if count == 2 || count == 4 {
			m.ContainerID = container.ID
			return
		}
	}
	err = errors.New("Container not found")
	return
}
