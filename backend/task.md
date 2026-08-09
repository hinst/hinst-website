Looking at commit ee5afcf519af910a616ed593296d87e4c991e83a, earlier I used Google API to ping my URLs through search indexing.
Now we want to use Google API again in a similar way, but this time instead of pinging URL, we only check if it is currently indexed.

Code example:
```python
from google.oauth2 import service_account
from googleapiclient.discovery import build

# Load credentials from the JSON file downloaded from Google Cloud
crendentials_path = 'path/to/your/service-account-key.json'
scopes = ['https://googleapis.com']
creds = service_account.Credentials.from_service_account_file(crendentials_path, scopes=scopes)

# Build the Search Console service object
service = build('searchconsole', 'v1', credentials=creds)

# Define request parameters
request_body = {
    "inspectionUrl": "https://example.com",
    "siteUrl": "https://example.com"
}

# Execute the inspection request
try:
    response = service.urlInspection().index().inspect(body=request_body).execute()

    # Parse the indexing status
    inspection_result = response.get('inspectionResult', {})
    index_status_result = inspection_result.get('indexStatusResult', {})
    verdict = index_status_result.get('verdict')

    print(f"Indexing Status: {verdict}")
    # Outputs: "COVERED" (Indexed), "NEUTRAL" (Not Indexed / Excluded), etc.

except Exception as e:
    print(f"An error occurred: {e}")

```

But we use go, so we should add URL checker call into `server\webAppAdmin.go` and return it as field record.GoogleSearchIndexingStatus = true or false.
