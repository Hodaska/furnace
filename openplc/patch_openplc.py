import sys

path = '/workdir/webserver/openplc.py'
content = open(path).read()

old = (
    "    def _rpc(self, msg, timeout=1000):\n"
    "        data = \"\"\n"
    "        if not self.runtime_status == \"Running\":\n"
    "            return data\n"
    "        try:\n"
    "            s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)\n"
    "            s.connect(('localhost', 43628))\n"
    "            s.send(f'{msg}\\n'.encode('utf-8'))\n"
    "            data = s.recv(timeout).decode('utf-8')\n"
    "            s.close()\n"
    "            self.runtime_status = \"Running\"\n"
    "        except socket.error as serr:\n"
    "            print(f'Socket error during {msg}, is the runtime active?')\n"
    "            self.runtime_status = \"Stopped\"\n"
    "        return data"
)

new = (
    "    def _rpc(self, msg, timeout=1000):\n"
    "        data = \"\"\n"
    "        if not self.runtime_status == \"Running\":\n"
    "            return data\n"
    "        import time as _t\n"
    "        for _attempt in range(10):\n"
    "            try:\n"
    "                s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)\n"
    "                s.connect(('localhost', 43628))\n"
    "                s.send(f'{msg}\\n'.encode('utf-8'))\n"
    "                data = s.recv(timeout).decode('utf-8')\n"
    "                s.close()\n"
    "                self.runtime_status = \"Running\"\n"
    "                return data\n"
    "            except socket.error:\n"
    "                if _attempt < 9:\n"
    "                    _t.sleep(1)\n"
    "                else:\n"
    "                    print(f'Socket error during {msg}, is the runtime active?')\n"
    "                    self.runtime_status = \"Stopped\"\n"
    "        return data"
)

if old not in content:
    print("ERROR: patch target not found in openplc.py", file=sys.stderr)
    sys.exit(1)

open(path, 'w').write(content.replace(old, new))
print("openplc.py patched OK - _rpc now retries 10x")
