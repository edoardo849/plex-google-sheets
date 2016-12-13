Plex to Google Spreadsheet exporter
===

Export your Plex library to a Google's cloud Spreadsheet. DO NOT USE THIS... yet. This project is under heavy development at the moment.

# Prerequisites

1. Use [this wizard](https://console.developers.google.com/flows/enableapi?apiid=sheets.googleapis.com&pli=1) to create or select a project in the Google Developers Console and automatically turn on the API. Click Continue, then Go to credentials.
1. On the Add credentials to your project page, click the Cancel button.
1. At the top of the page, select the OAuth consent screen tab. Select an Email address, enter a Product name if not already set, and click the Save button.
1. Select the Credentials tab, click the Create credentials button and select OAuth client ID.
1. Select the application type Other, enter the name "Google Sheets API Quickstart", and click the Create button.
1. Click OK to dismiss the resulting dialog.
Click the file_download (Download JSON) button to the right of the client ID.
1. Move this file to your working directory and rename it client_secret.json

# Features

Import your Plex library list into Google Sheets.
