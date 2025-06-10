# Overview

This is a one-off script to migrate our CSS grafana dashboard folders from AWS into Openshift. It will move all the old folders in the old instance to the new one, and prefix folders with "CSS App - " to prevent foldername conflicts.

## Setup

Set the variables at the top of the file:

**OLD_GRAFANA_URL**: Root URL of the old instance.
**OLD_API_KEY**: A service account token with admin permission to the old instance.
**NEW_GRAFANA_URL**: Root URL of the new instance.
**NEW_API_KEY**:  A service account token with admin permission to the new instance.
**NEW_CSS_DS_ID**: Datasource ID for the CSS postgres connection in the new instance, can be found in panel JSON.
**OLD_CSS_DS_ID**: Datasource ID for the CSS postgres connection in the old instance, can be found in panel JSON.

## Usage:

Run `node index.js`.
