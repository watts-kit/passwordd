#!/usr/bin/env python
# pylint: disable=invalid-name
import socket
import json
import sys
import os
# apt-get install python-configargparse
import argparse

def parseOptions():
    '''Parse the commandline options'''
    path_of_executable = sys.argv[0]
    folder_of_executable = os.path.split(path_of_executable)[0]

    parser = argparse.ArgumentParser(
        description='''Synchronisation of password daemons''')
    parser.add_argument('--pullhost', '--pull',
                        help='hostname to pull passwords from',
                        default='127.0.0.1')
    parser.add_argument('--pullport', '--port',
                        help='port to pull passwords from',
                        default=6969, type=int)
    parser.add_argument('--pushhost', '--push',
                        help='hostname to push passwords to',
                        default='127.0.0.1')
    parser.add_argument('--verbose', '-v', action='count')
    parser.add_argument('--pushport',
                        help='port to push passwords to',
                        default=6969, type=int)
    parser.add_argument('key', action='append', nargs='+',
                        help='List of keys to operate on')
    return parser.parse_args()

def connect(host_port):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(host_port)
    return s

def receive_and_close(sock):
    BUFFER_SIZE = 1024
    data = sock.recv(BUFFER_SIZE)
    sock.close()
    response = json.loads(data)
    return response

def get_secret(host_port, key):
    if args.verbose > 0:
        print ("connecting to: %s" % str(host_port))
    sock = connect(host_port)
    message = json.dumps({"action":"get", "key":key})
    sock.send(message)
    if args.verbose > 2:
        print ("sending message: '%s'" % message)
    response = receive_and_close(sock)
    if args.verbose > 2:
        print ("Response: '%s'" % response)
    if response["result"] == "ok":
        return response["value"]
    return ""

def set_secret(host_port, key, secret):
    sock = connect(host_port)
    message = json.dumps({"action":"set", "key":key, "value":secret})
    sock.send(message)
    response = receive_and_close(sock)
    if response["result"] == "ok":
        return "secret set"
    return "secret NOT set"

def overwrite_secret(host_port, key, secret):
    sock = connect(host_port)
    message = json.dumps({"action":"overwrite", "key":key, "value":secret})
    sock.send(message)
    if args.verbose > 2:
        print ("sending message: '%s'" % message)
    response = receive_and_close(sock)
    if response["result"] == "ok":
        return "secret set"
    return "secret NOT set"


args = parseOptions()
def main():
    if args.verbose > 0:
        print("keys: %s" % str(args.key))

    source = (args.pullhost, args.pullport)
    destin = (args.pushhost, args.pushport)
    if args.verbose > 1:
        print('Source: "%s"'%str(source))
        print('Destin: "%s"'%str(destin))

    for key in args.key[0]:
        # check if key is available in destination
        dst_secret = get_secret(destin, key)
        if args.verbose > 2:
            print ("got dst_secret: %s" % dst_secret)
        # if dst_secret: # i.e. if overwriting is not intended 
        #     print('Secret "%s" is already available in destination at %s:%d' %
        #           (key, args.pullhost, args.pullport))
        #     # -=> skip to next
        #     continue
        if args.verbose > 1:
            print ('    Got dest:%s' % dst_secret)
        # check if key is available in source
        src_secret = get_secret(source, key)
        if args.verbose > 2:
            print ('    Got src: %s' % src_secret)
        if not src_secret:
            print('Cannot get secret for "%s" from %s:%d' %
                  (key, args.pullhost, args.pullport))
            # -=> warning and skip to next
            continue
        if args.verbose > 2:
            print ('    Got src: %s' % src_secret)
        # get in source // happened above
        # set in destination
        result = set_secret(destin, key, src_secret)
        if args.verbose > 0:
            print ("setting new secret. Result: '%s' " % result)
        if result == 'secret NOT set':
            if args.verbose > 0:
                print ("setting didn't work, overwriting")
            result = overwrite_secret(destin, key, src_secret)
            if args.verbose > 0:
                print ("overwriting old secret. Result: '%s' " % result)

        if args.verbose > 0:
            print('copied secret for "%s"' % key)

if __name__ == "__main__":
    main()
