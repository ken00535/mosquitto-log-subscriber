package client

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Host struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	CaPath   string `json:"caPath"`
	CertPath string `json:"certPath"`
	KeyPath  string `json:"keyPath"`
}

// NewClient new a mqtt client
func NewClient(host Host) (mqtt.Client, error) {
	certpool := x509.NewCertPool()
	// config.CAPath = filepath.Join(path, path.CaPath)
	// config.CertPath = filepath.Join(path, path.CertPath)
	// config.KeyPath = filepath.Join(path, path.KeyPath)
	pemCerts, err := ioutil.ReadFile(host.CaPath)
	if err != nil {
		panic(err)
	}
	ok := certpool.AppendCertsFromPEM(pemCerts)
	if !ok {
		panic("failed to parse root ca")
	}

	cert, err := tls.LoadX509KeyPair(host.CertPath, host.KeyPath)
	if err != nil {
		panic(err)
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}

	url := "tcps://" + host.IP + ":" + strconv.Itoa(host.Port)
	opts := mqtt.NewClientOptions().AddBroker(url).
		SetUsername(host.Username).SetPassword(host.Password)
	tlsConfig := tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}
	opts.SetTLSConfig(&tlsConfig)
	opts.SetClientID("log-subscriber")
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}
