package main

import (
	ninep "code.google.com/p/go9p/p"
	"code.google.com/p/go9p/p/clnt"
	"code.google.com/p/go9p/p/srv"
	"flag"
	"fmt"
	"os"
	"os/user"
	"strings"
)

var (
	srcAddr    = flag.String("s", "<none>", "Source filesystem")
	dstAddr    = flag.String("d", "<none>", "Destination filesystem")
	listenAddr = flag.String("addr", "tcp!:5640", "Listen address")
)

// clnt.Mount wants a user, let's wrap os/user's notion and give it the right
// methods.
type User struct {
	raw *user.User
}

func (u *User) Name() string {
	return u.raw.Username
}

func (u *User) Id() int {
	// FIXME: This is not a reasonable solution. I'm kinda
	// punting on the whole user/auth thing for the moment;
	// will have to deal with it eventually.
	return 100
}

func (u *User) Groups() []ninep.Group {
	// FIXME: see above
	return nil
}

func (u *User) IsMember(g ninep.Group) bool {
	// FIXME: see above
	return false
}

// The core data structure for cow9p. This essentially *is* the server.
type CowFS struct {
	src, dst *clnt.Clnt
}

// Starts the server loop. Note that this doesn't return until the server
// shuts down; If you want it to run in a separate goroutine you'll need
// to do that yourself.
//
// net and addr are the network and address on which to listen.
func (fs *CowFS) Serve(net, addr string) {
	reqin := make(chan *srv.Req)
	sv := srv.Srv{Reqin: reqin}
	sv.StartNetListener(net, addr)
	for {
		req := <-reqin
		switch req.Tc.Type {
		// TODO: actually handle messages
		}
	}
}

// Given the networks & addresses of the source and destination filesystems,
// this will connect to them and return a *CowFS ready to use them. (or an
// error)
func Mount(netSrc, addrSrc, netDst, addrDst string, user ninep.User) (*CowFS, error) {
	// the empty string here (and below) is the aname - which subtree to
	// connect to. right now we're just punting on providing a nice interface
	// to this. you get the whole tree.
	src, err := clnt.Mount(netSrc, addrSrc, "", user)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to source fs : %s\n", err)
	}
	dst, err := clnt.Mount(netDst, addrDst, "", user)
	if err != nil {
		src.Unmount()
		return nil, fmt.Errorf("Error connecting to destination fs : %s\n", err)
	}
	return &CowFS{src: src, dst: dst}, nil
}

// Merge several errors (including possible nil values) into one.
// The output error will be nil if and only if the input error
// is nil. This is handy for checking for the presence of any of
// several errors.
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

// splits a network address of the form net!addr into its component parts.
// This is needed because like most of the network related go packages,
// go9p expects these as separate arguments.
//
// returns an error if the string is not of the above form, otherwise
// err will be nil.
func splitNetAddr(addr string) (netPart, addrPart string, err error) {
	strs := strings.SplitN(addr, "!", 2)
	if len(strs) != 2 {
		return "", "", fmt.Errorf("Invalid filesystem address : %s\n", addr)
	}
	return strs[0], strs[1], nil
}

func main() {
	flag.Parse()

	netSrc, addrSrc, errSrc := splitNetAddr(*srcAddr)
	netDst, addrDst, errDst := splitNetAddr(*dstAddr)
	netListen, addrListen, errListen := splitNetAddr(*listenAddr)

	osUser, userErr := user.Current()

	if err := mergeErrs(errSrc, errDst, errListen, userErr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		fs, err := Mount(netSrc, addrSrc, netDst, addrDst, &User{osUser})
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		} else {
			fs.Serve(netListen, addrListen)
		}
	}
}
