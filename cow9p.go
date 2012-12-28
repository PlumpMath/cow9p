package main

import (
	"code.google.com/p/go9p/p/clnt"
	"code.google.com/p/go9p/p/srv"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	srcAddr = flag.String("s", "<none>", "Source filesystem")
	dstAddr = flag.String("d", "<none>", "Destination filesystem")
	listenAddr = flag.String("addr", "tcp!:5640", "Listen address")
)

type CowFS struct {
	src, dst *clnt.Clnt
}

func (fs *CowFS) Serve(net, addr string) {
	reqin := make(chan *srv.Req)
	sv := srv.Srv{ Reqin: reqin }
	sv.StartNetListener(net, addr)
	for {
		req := <-reqin
		switch req.Tc.Type {
			// TODO: actually handle messages
		}
	}
}

func Mount(netSrc, addrSrc, netDst, addrDst string) (*CowFS, error) {
	src, err := clnt.Mount(netSrc, addrSrc, /* TODO: need to provide aname/user. */)
	if(err != nil) {
		return nil, fmt.Errorf("Error connecting to source fs : %s\n", err)
	}
	dst, err := clnt.Mount(netDst, addrDst, /* TODO: need to provide aname/user. */)
	if(err != nil) {
		src.Unmount()
		return nil, fmt.Errorf("Error connecting to destination fs : %s\n", err)
	}
	return &CowFS{src: src, dst: dst}, nil
}

func mergeErrs(errs ...error) error {
	err := fmt.Errorf("")
	haveErr := false
	for _, v := range errs {
		if v != nil {
			err = fmt.Errorf("%s\n", v)
			haveErr = true
		}
	}
	if haveErr {
		return fmt.Errorf("Errors : %s\n", err)
	}
	return nil
}

func splitNetAddr(addr string) (netPart, addrPart string, err error) {
	strs := strings.SplitN(addr, "!", 1)
	if len(strs) != 2 {
		return "", "", fmt.Errorf("Invalid filesystem address : %s\n", addr)
	}
	return strs[0], strs[1], nil
}

func main () {
	flag.Parse()
	netSrc, addrSrc, errSrc := splitNetAddr(*srcAddr)
	netDst, addrDst, errDst := splitNetAddr(*dstAddr)
	netListen, addrListen, errListen := splitNetAddr(*listenAddr)

	if err := mergeErrs(errSrc, errDst, errListen); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		fs, err := Mount(netSrc, addrSrc, netDst, addrDst)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		} else {
			fs.Serve(netListen, addrListen)
		}
	}
}
