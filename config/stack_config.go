package config

import (
	"github.com/coreos/kube-aws/coreos/userdatavalidation"
	"github.com/coreos/kube-aws/filereader/jsontemplate"
	"github.com/coreos/kube-aws/gzipcompressor"
)

type StackConfig struct {
	*Config
	StackTemplateOptions
	UserDataWorker        string
	UserDataController    string
	userDataEtcd          string
	ControllerSubnetIndex int
}

type CompressedStackConfig struct {
	*StackConfig
	UserDataEtcd string
}

func (c *StackConfig) ValidateUserData() error {
	err := userdatavalidation.Execute([]userdatavalidation.Entry{
		{Name: "UserDataWorker", Content: c.UserDataWorker},
		{Name: "UserDataController", Content: c.UserDataController},
		{Name: "UserDataEtcd", Content: c.userDataEtcd},
	})

	return err
}

func (c *StackConfig) Compress() (*CompressedStackConfig, error) {
	var err error
	var compressedEtcdUserData string

	if compressedEtcdUserData, err = gzipcompressor.CompressString(c.userDataEtcd); err != nil {
		return nil, err
	}

	var stackConfig CompressedStackConfig
	stackConfig.StackConfig = &(*c)
	stackConfig.UserDataEtcd = compressedEtcdUserData

	return &stackConfig, nil
}

func (c *CompressedStackConfig) RenderStackTemplateAsBytes() ([]byte, error) {
	bytes, err := jsontemplate.GetBytes(c.StackTemplateTmplFile, *c, c.PrettyPrint)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

func (c *CompressedStackConfig) RenderStackTemplateAsString() (string, error) {
	bytes, err := c.RenderStackTemplateAsBytes()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
