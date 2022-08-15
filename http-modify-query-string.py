"""Modify HTTP query parameters."""
from mitmproxy import http

def request(flow: http.HTTPFlow) -> None:
    try:
        with open("/tmp/external.port") as f:
            content = f.read()

            if content.isdigit() and 1024 < int(content) < 65535 :
                if "port" in flow.request.query and flow.request.query["port"].isdigit() :
                    flow.request.query["port"] = content
                    print("successfully modify port to", content)
                else :
                    print("url doesn't contain port parameters or port is null, skipping")
                    return
            else :
                print("invalid port:", content)
                return
    except OSError:
        ptinr("OSError:", OSError)
