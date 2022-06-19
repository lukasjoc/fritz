from fritz.typing import EXIT_SUCCESS, Path
from fritz.api import Router
import argparse
import sys
import requests

def current_ip() -> str:
    res = requests.get("https://httpbin.org/ip")
    return res.json().get("origin")

def main(args) -> None:
    if args.reboot:
        router = Router(args.config)
        connection = router.connect()
        print(f"Rebooting {connection.modelname}@{router.connector.connection_host}...",)
        connection.reboot()
        print("Get a 🍵 and wait until the router is reachable again.")
    elif args.reconnect:
        router = Router(args.config)
        connection = router.connect()
        print(f"Current Public IP: {current_ip()}")
        print(f"Reconnecting...")
        connection.reconnect()
        print(f"Current Public IP: {current_ip()}")
    elif args.info:
        router = Router(args.config)
        connection = router.connect()
        connector  = router.connector
        print("Model Details: ")
        print()
        print(f"Model:   {connection.modelname}")
        print(f"Version: {connection.system_version}")
        print()
        print("Connection Details: ")
        print(f"Host:    {connector.connection_host}")
        print(f"Port:    {connector.connection_port}")
        print(f"TLS:     {connector.connection_use_tls}")

    sys.exit(EXIT_SUCCESS)

if __name__ == "__main__":
    def str2bool(v):
        if isinstance(v, bool): return v
        if v.lower() in ('yes', 'true', 't', 'y', '1'): return True
        elif v.lower() in ('no', 'false', 'f', 'n', '0'): return False
        else:
            raise argparse.ArgumentTypeError('Boolean value expected.')

    parser = argparse.ArgumentParser(description="FRITZ!Box (Simple) Interface")

    parser.add_argument("--info", type=str2bool, nargs='?',
                        const=True, default=False,
                        help="connect and show info about the router")

    parser.add_argument("--reboot", type=str2bool, nargs='?',
                        const=True, default=False,
                        help="connect and reboot the router")

    parser.add_argument("--reconnect", type=str2bool, nargs='?',
                        const=True, default=False,
                        help="connect and reconnect with new IP from ISP")


    parser.add_argument("--config", type=Path, default="config.yml",
                        help="path to the configuration file")


    parser.add_argument("--version", action="version", version="%(prog)s 0.1.0",
        help="print the current version",
    )

    args = parser.parse_args()

    main(args)
