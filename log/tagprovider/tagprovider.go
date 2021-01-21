package tagprovider

type Provider func() (key, value string)
