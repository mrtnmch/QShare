# QShare

PoC of sharing data from an iOS (or Android) device to a browser without an account or a database, based on RSA and QR codes. Only text is supported for now.

1. Download and install the iOS extension: https://www.icloud.com/shortcuts/db03f99f0a594b4896357424422226b6
2. Go to https://mxmx.app/qshare
3. Share any text via the system dialog and pick "QShare". Once the camera app is open, scan the QR.
4. The content gets passed into the browser. If it's a link, the browser is redirected right away.

## How it works

A unique code and private/public key pair are generated during the server startup (different every time). When a QR image is requested, the code is signed/encrypted with the public key (including sha256 & random salt), encoded into base64 and turned into QR (note: it's not stored anywhere).

When the QR gets presented to a user, the encrypted code is also passed in the HTTP response headers and browser starts polling the result with the code. The uploaded content also contains this code - therefore the content and browser session can be tied to each other, without any account or database whatsoever.

The request looks like this:
```
curl --location --request POST 'https://mxmx.app/qshare/api/upload' \
--header 'Content-Type: application/json' \
--data-raw '{
	"code": "<base64 code from QR>",
	"content": "Text or link (gets automatically redirected)"
}'
```

Because the QR code has been encrypted via the public key, it can be decrypted with and only with the private key. Therefore, during the content upload, the code gets decrypted and checked if it equals the code generated during the startup. This way, only valid, server-generated codes are accepted.
