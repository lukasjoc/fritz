from fritzconnection import FritzConnection
from .typing import Config, Path, Port
from yaml import load as yamlload, FullLoader as yamlFullLoader

FRITZ_TCP_PORT = 49000
FRITZ_TLS_PORT = 49443

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
        if port is None and use_tls: return FRITZ_TLS_PORT
        elif port is None: return FRITZ_TCP_PORT
        else: return port

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
