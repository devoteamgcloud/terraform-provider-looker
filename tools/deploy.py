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

def get_headers() -> dict:
    return {
        "Authorization": f"Bearer {TF_TOKEN}",
        "Content-Type": "application/vnd.api+json"
    }

def add_platform_endpoint(filename: str, shasum: str, version: str) -> str:
    headers = get_headers()
    payload = {
        "data": {
            "type": "registry-provider-version-platforms",
            "attributes": {
            "os": "linux",
            "arch": "amd64",
            "shasum": shasum,
            "filename": filename
            }
        }
    }
    response = requests.post(BASE_URL+f"/v2/organizations/{WORKSPACE}/registry-providers/private/{WORKSPACE}/{PROVIDER}/versions/{version}/platforms", headers=headers, json = payload)
    if 200 <= response.status_code < 300:
        result = response.json()
        upload_link = result["data"]["links"]["provider-binary-upload"]
        return upload_link
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
    print("Error creating version, check if version already exists.")
    sys.exit(1)

def upload_shasums(path_sha: str, path_sha_sig: str, link_sha: str, link_sha_sig: str) -> None:
    with open(path_sha, "rb") as file:
        response = requests.post(link_sha, files = {"form_field_name": file})
        if response.status_code >= 300:
            print(f"Error uploading {path_sha} to {link_sha}.")
            sys.exit(1)

    with open(path_sha_sig, "rb") as file:
        response = requests.post(link_sha_sig, files = {"form_field_name": file})
        if response.status_code >= 300:
            print(f"Error uploading {path_sha_sig} to {link_sha_sig}.")
            sys.exit(1)

def upload_file(file_path: str, link: str):
    with open(file_path, "rb") as file:
        response = requests.post(link, files = {"form_field_name": file})
        if response.status_code >= 300:
            print(f"Error uploading {file_path} to {link}.")
            sys.exit(1)
    

def find_index(term: str, string: str) -> tuple:
    myarray = []
    for m in re.finditer(term, string):
        myarray.append(m.start())
    return myarray[0]+1, myarray[1]

def get_sha256(path: str) -> str:
    with open(path,"rb") as f:
        bytes = f.read() # read entire file as bytes
        readable_hash = hashlib.sha256(bytes).hexdigest()
        return readable_hash

def main():
    sig_file = ""
    sha_file = ""
    release_file = ""
    for file in os.listdir("dist"):
        try:
            if file[-14:] == "SHA256SUMS.sig":
                sig_file = file
        except IndexError:
            pass
        try:
            if file[-10:] == "SHA256SUMS":
                sha_file = file
        except IndexError:
            pass
        try:
            if file[-15] == "linux_amd64.zip":
                release_file = file
        except IndexError:
            pass
        

    index_start, index_end = find_index("_", sha_file)
    version = sig_file[index_start:index_end]
    sha_file = os.path.abspath(f"dist/{sha_file}")
    sig_file = os.path.abspath(f"dist/{sig_file}")
    sha_link, sig_link = add_version_endpoint(version)
    upload_shasums(sha_file, sig_file, sha_link, sig_link)
    release_sha256 = get_sha256(release_file)
    upload_link = add_platform_endpoint(release_file, release_sha256, version)
    release_file = os.path.abspath(f"dist/{release_file}")
    upload_file(release_file, upload_link)



if __name__ == "__main__":
    main()
