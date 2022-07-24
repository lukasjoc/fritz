import argparse
import sys
import requests
from fritzconnection import FritzConnection
from yaml import load as yamlload, FullLoader as yamlFullLoader

FRITZ_TCP_PORT = 49000
FRITZ_TLS_PORT = 49443

Config = dict[any]
Path = str
Port = str or int
EXIT_SUCCESS = 0


class Connector:
    def __init__(self, config_file: Path) -> None:
        self.config_file = config_file

    @staticmethod
    def read_config(config_file: Path) -> Config:
        with open(config_file, "r") as config:
            data = yamlload(config, Loader=yamlFullLoader)
        return data

    @property
    def connection_host(self) -> str:
        return self.read_config(self.config_file).get('host')

    @property
    def connection_use_tls(self) -> str:
        return self.read_config(self.config_file).get('use_tls') or False

    @property
    def connection_port(self) -> Port:
        config = self.read_config(self.config_file)
        port = config.get("port")
        use_tls = config.get("tls")
        if port is None and use_tls:
            return FRITZ_TLS_PORT
        elif port is None:
            return FRITZ_TCP_PORT
        else:
            return port


class Router:
    def __init__(self, config_file: Path, connector: Connector = None) -> None:
        self.config_file = config_file
        self.router = None
        connector = connector or Connector(config_file)
        self.connector = connector

    def connect(self, pool_connections: int = 3, pool_maxsize: int = 3) -> FritzConnection:
        config = self.connector.read_config(self.config_file)
        print("Connected to FRITZ!Box")
        self.router = FritzConnection(
            address=config.get("host"),
            user=config.get("user"),
            password=config.get("password"),
            use_tls=config.get("use_tls") or False,
            pool_connections=pool_connections,
            pool_maxsize=pool_maxsize,
        )
        return self.router


def current_ip() -> str:
    res = requests.get("https://httpbin.org/ip")
    return res.json().get("origin")


def main(args) -> None:
    if args.reboot:
        router = Router(args.config)
        connection = router.connect()
        print(
            f"Rebooting {connection.modelname}@{router.connector.connection_host}...",)
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
        connector = router.connector
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
        if isinstance(v, bool):
            return v
        if v.lower() in ('yes', 'true', 't', 'y', '1'):
            return True
        elif v.lower() in ('no', 'false', 'f', 'n', '0'):
            return False
        else:
            raise argparse.ArgumentTypeError('Boolean value expected.')

    parser = argparse.ArgumentParser(
        description="FRITZ!Box (Simple) Interface")

    parser.add_argument(
        "--info", type=str2bool, nargs='?',
        const=True, default=False,
        help="connect and show info about the router"
    )

    parser.add_argument(
        "--reboot", type=str2bool, nargs='?',
        const=True, default=False,
        help="connect and reboot the router"
    )

    parser.add_argument(
        "--reconnect", type=str2bool, nargs='?',
        const=True, default=False,
        help="connect and reconnect with new IP from ISP"
    )

    parser.add_argument(
        "--config", type=Path, default="config.yml",
        help="path to the configuration file"
    )

    parser.add_argument(
        "--version", action="version", version="%(prog)s 0.2.0",
        help="print the current version",
    )

    args = parser.parse_args()

    main(args)
