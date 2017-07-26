package udptest

import (
	"net"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TCP", func() {
	const (
		resolveNetwork = "tcp"
		network        = "tcp"
	)

	It("tests", func() {
		dataChan := make(chan []byte)
		laddr, err := net.ResolveTCPAddr(resolveNetwork, "0.0.0.0:0")
		Expect(err).ToNot(HaveOccurred())
		ln, err := net.ListenTCP(network, laddr)
		Expect(err).ToNot(HaveOccurred())
		defer ln.Close()

		go func() {
			defer GinkgoRecover()
			data := make([]byte, 100)
			conn, err := ln.Accept()
			Expect(err).ToNot(HaveOccurred())
			n, err := conn.Read(data)
			Expect(err).ToNot(HaveOccurred())
			data = data[:n]
			dataChan <- data
		}()

		addrString := "localhost:" + strconv.Itoa(ln.Addr().(*net.TCPAddr).Port)
		addr, err := net.ResolveTCPAddr(resolveNetwork, addrString)
		Expect(err).NotTo(HaveOccurred())
		Expect(addr.Port).To(Equal(ln.Addr().(*net.TCPAddr).Port))

		data := make([]byte, 77)
		for i := 0; i < 10; i++ {
			conn, err := net.DialTCP(network, nil, addr)
			Expect(err).ToNot(HaveOccurred())
			defer conn.Close()
			_, err = conn.Write(data)
			Expect(err).ToNot(HaveOccurred())
			time.Sleep(time.Millisecond)
		}
		Eventually(dataChan).Should(Receive(Equal(data)))
	})
})
