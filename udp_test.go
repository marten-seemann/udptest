package udptest

import (
	"net"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const network = "udp"

var _ = Describe("UDP", func() {
	It("tests", func() {
		dataChan := make(chan []byte)
		laddr, err := net.ResolveUDPAddr(network, "0.0.0.0:0")
		Expect(err).ToNot(HaveOccurred())
		ln, err := net.ListenUDP(network, laddr)
		Expect(err).ToNot(HaveOccurred())
		defer ln.Close()

		go func() {
			defer GinkgoRecover()
			data := make([]byte, 100)
			n, _, err := ln.ReadFrom(data)
			Expect(err).ToNot(HaveOccurred())
			data = data[:n]
			dataChan <- data
		}()

		addrString := "localhost:" + strconv.Itoa(ln.LocalAddr().(*net.UDPAddr).Port)
		addr, err := net.ResolveUDPAddr(network, addrString)
		Expect(err).NotTo(HaveOccurred())
		Expect(addr.Port).To(Equal(ln.LocalAddr().(*net.UDPAddr).Port))

		data := make([]byte, 77)
		for i := 0; i < 10; i++ {
			conn, err := net.DialUDP(network, nil, addr)
			Expect(err).ToNot(HaveOccurred())
			defer conn.Close()
			_, err = conn.Write(data)
			Expect(err).ToNot(HaveOccurred())
			time.Sleep(time.Millisecond)
		}
		Eventually(dataChan).Should(Receive(Equal(data)))
	})
})
