package udptest

import (
	"crypto/rand"
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UDP", func() {
	It("tests", func() {
		dataChan := make(chan []byte)
		serverAddrChan := make(chan *net.UDPAddr)
		go func() {
			defer GinkgoRecover()
			addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
			Expect(err).ToNot(HaveOccurred())
			ln, err := net.ListenUDP("udp", addr)
			Expect(err).ToNot(HaveOccurred())
			defer ln.Close()
			serverAddrChan <- ln.LocalAddr().(*net.UDPAddr)
			data := make([]byte, 100)
			n, _, err := ln.ReadFrom(data)
			Expect(err).ToNot(HaveOccurred())
			data = data[:n]
			dataChan <- data
		}()

		raddr := <-serverAddrChan
		conn, err := net.DialUDP("udp", nil, raddr)
		Expect(err).ToNot(HaveOccurred())
		data := make([]byte, 77)
		_, err = rand.Read(data)
		Expect(err).ToNot(HaveOccurred())
		_, err = conn.Write(data)
		Expect(err).ToNot(HaveOccurred())
		Eventually(dataChan).Should(Receive(Equal(data)))
	})
})
