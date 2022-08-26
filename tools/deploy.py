import os
import re
import sys
import hashlib
import requests

PUBLIC_GPG = os.environ.get("GPG_PUBLIC_KEY").replace("\n", "\\n")
PRIVATE_GPG = os.environ.get("GPG_PRIVATE_KEY").replace("\n", "\\n")
TF_TOKEN = os.environ.get("TF_TOKEN")
KEY_ID = os.environ.get("KEY_ID")
BASE_URL = "https://app.terraform.io/api"
WORKSPACE = "reprise-digital"
PROVIDER = "looker"

def get_sha256(path: str) -> str:
    with open(path,"rb") as f:
        bytes = f.read()
        readable_hash = hashlib.sha256(bytes).hexdigest()
        return readable_hash

class File:
    def __init__(self, name: str, version: str= None, platform: str= None, arch: str= None, upload_link: str = None) -> None:
        self.name = name
        self.version = version
        self.platform = platform
        self.arch = arch
        self.upload_link = upload_link
        self.path = os.path.abspath(f"dist/{name}")
        self.sha256 = get_sha256(self.path)

def get_headers() -> dict:
    return {
        "Authorization": f"Bearer {TF_TOKEN}",
        "Content-Type": "application/vnd.api+json"
    }

def add_platform_endpoint(file: File) -> str:
    headers = get_headers()
    payload = {
        "data": {
            "type": "registry-provider-version-platforms",
            "attributes": {
            "os": file.platform,
            "arch": file.arch,
            "shasum": file.sha256,
            "filename": file.name
            }
        }
    }
    response = requests.post(BASE_URL+f"/v2/organizations/{WORKSPACE}/registry-providers/private/{WORKSPACE}/{PROVIDER}/versions/{file.version}/platforms", headers=headers, json = payload)
    if 200 <= response.status_code < 300:
        result = response.json()
        upload_link = result["data"]["links"]["provider-binary-upload"]
        return upload_link
    print(response.content)
    print("Error creating version, check if version already exists.")
    sys.exit(1)

def add_version_endpoint(version: str) -> tuple:
    headers = get_headers()
    payload = {
        "data": {
            "type": "registry-provider-versions",
            "attributes": {
            "version": version,
            "key-id": KEY_ID,
            "protocols": ["5.0"]
            }
        }
    }
    response = requests.post(BASE_URL+f"/v2/organizations/{WORKSPACE}/registry-providers/private/{WORKSPACE}/{PROVIDER}/versions", headers=headers, json=payload)
    if 200 <= response.status_code < 300:
        result = response.json()
        shasums = result["data"]["links"]["shasums-upload"]
        shasums_sig = result["data"]["links"]["shasums-sig-upload"]
        return shasums, shasums_sig
    print(response.content)
    print("Error creating version, check if version already exists.")
    sys.exit(1)

def upload_file(file: File):
    with open(file.path, "rb") as file_data:
        response = requests.put(file.upload_link, files = {"upload_file": file_data})
        if response.status_code >= 300:
            print(response.content)
            print(f"Error uploading {file.path} to {file.upload_link}.")
            sys.exit(1)

def main():
    files: list[File] = []
    for file in os.listdir("dist"):
        try:
            if file[-4:] == ".zip":
                try:
                    _, version, platform, arch = file.split("_")
                    files.append(File(file, version, platform, arch[:-4], None))
                except:
                    pass     
            elif file[-10:] == "SHA256SUMS":
                sha_file = File(name = file)
            elif file[-14:] == "SHA256SUMS.sig":
                shasig_file = File(name = file)
        except IndexError:
            pass
    sha_file.upload_link, shasig_file.upload_link = add_version_endpoint(version)
    upload_file(sha_file) 
    upload_file(shasig_file)
    for file in files:
        file.upload_link = add_platform_endpoint(file)
        upload_file(file)

if __name__ == "__main__":
    main()
