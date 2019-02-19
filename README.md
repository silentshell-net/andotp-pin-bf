# andotp-pin-bf
Brute force andOTP PIN - guess why i needed to write this 

## Requirements:
Obtain the shared_prefs.xml file from your android device, you should need a rooted phone in order to do that. In the XML file extract the `pref_auth_credentials` and `pref_auth_salt` keys and pass it to the script.

## Usage:
Just pass both values to the script :
`./andotp-pin-bf -auth aaaaaaaaaaaaa== -salt aaaaaaaaaaaaa==`
