#!/usr/bin/env bash
set -euo pipefail

# ---------------------------------------------------------------------------
# Grafana Dashboard Exporter
#
# Usage:
#   ./script.sh -u <grafana-url> -t <service-account-token> [-o <output-dir>]
#
# Options:
#   -u  Grafana instance URL  (e.g. https://grafana.example.com)
#   -t  Grafana service account token
#   -o  Output directory (default: ./grafana-dashboards-export)
#   -h  Show this help message
# ---------------------------------------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# --- defaults ---------------------------------------------------------------
OUTPUT_DIR="${SCRIPT_DIR}/grafana-dashboards-export"
GRAFANA_URL=""
TOKEN=""

# --- helpers ----------------------------------------------------------------
usage() {
  grep '^#' "$0" | grep -v '#!/' | sed 's/^# \{0,2\}//'
  exit 0
}

err() {
  echo "[ERROR] $*" >&2
  exit 1
}

warn() {
  echo "[WARN]  $*" >&2
}

info() {
  echo "[INFO]  $*"
}

# --- arg parsing ------------------------------------------------------------
while getopts ":u:t:o:h" opt; do
  case "${opt}" in
    u) GRAFANA_URL="${OPTARG}" ;;
    t) TOKEN="${OPTARG}" ;;
    o) OUTPUT_DIR="${OPTARG}" ;;
    h) usage ;;
    :) err "Option -${OPTARG} requires an argument." ;;
    \?) err "Unknown option: -${OPTARG}" ;;
  esac
done

# --- validation -------------------------------------------------------------
[[ -z "${GRAFANA_URL}" ]] && err "Grafana instance URL is required. Use -u <url>"
[[ -z "${TOKEN}" ]]       && err "Service account token is required. Use -t <token>"

# Strip trailing slash for consistency
GRAFANA_URL="${GRAFANA_URL%/}"

# Validate URL format
if [[ ! "${GRAFANA_URL}" =~ ^https?:// ]]; then
  err "Invalid URL '${GRAFANA_URL}': must start with http:// or https://"
fi

# Validate token is non-whitespace
if [[ "${TOKEN}" =~ ^[[:space:]]*$ ]]; then
  err "Token must not be empty or whitespace only."
fi

# Ensure required tools are available
for cmd in curl jq; do
  command -v "${cmd}" &>/dev/null || err "'${cmd}' is required but not found in PATH."
done

# ---------------------------------------------------------------------------
# Verify connectivity and token validity by calling the /api/org endpoint
# ---------------------------------------------------------------------------
info "Verifying connection to ${GRAFANA_URL} ..."

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  "${GRAFANA_URL}/api/org")

case "${HTTP_CODE}" in
  200) info "Authentication successful." ;;
  401) err "Authentication failed (HTTP 401). Check that the token is valid." ;;
  403) err "Access denied (HTTP 403). The token may lack the required permissions." ;;
  000) err "Could not connect to '${GRAFANA_URL}'. Check the URL and network access." ;;
  *)   err "Unexpected response from Grafana API: HTTP ${HTTP_CODE}" ;;
esac

# ---------------------------------------------------------------------------
# Helper: Grafana API GET
# ---------------------------------------------------------------------------
grafana_get() {
  local path="$1"
  curl -s -f \
    -H "Authorization: Bearer ${TOKEN}" \
    -H "Content-Type: application/json" \
    "${GRAFANA_URL}${path}"
}

# ---------------------------------------------------------------------------
# Fetch all folders (General folder has uid="" / no entry in /api/folders)
# ---------------------------------------------------------------------------
info "Fetching folder list ..."
FOLDERS_JSON=$(grafana_get "/api/folders?limit=1000")

# Build a uid->title map using a temp file (bash 3 compatible)
# Format: <uid>TAB<title>
FOLDER_MAP_FILE=$(mktemp)
trap 'rm -f "${FOLDER_MAP_FILE}"' EXIT

# Seed with the implicit General folder (uid may be empty string or "general")
printf 'general\tGeneral\n' >> "${FOLDER_MAP_FILE}"
printf '\tGeneral\n'        >> "${FOLDER_MAP_FILE}"

echo "${FOLDERS_JSON}" | jq -r '.[] | [.uid, .title] | @tsv' >> "${FOLDER_MAP_FILE}"

folder_lookup() {
  local key="$1"
  awk -F'\t' -v k="${key}" '$1 == k { print $2; exit }' "${FOLDER_MAP_FILE}"
}

FOLDER_COUNT=$(wc -l < "${FOLDER_MAP_FILE}" | tr -d ' ')
info "Found ${FOLDER_COUNT} folder entry(ies) (including General)."

# ---------------------------------------------------------------------------
# Fetch all dashboards (search API returns metadata only)
# ---------------------------------------------------------------------------
info "Fetching dashboard list ..."
DASHBOARDS_JSON=$(grafana_get "/api/search?type=dash-db&limit=5000")
TOTAL=$(echo "${DASHBOARDS_JSON}" | jq 'length')
info "Found ${TOTAL} dashboard(s)."

if [[ "${TOTAL}" -eq 0 ]]; then
  warn "No dashboards found. Nothing to export."
  exit 0
fi

mkdir -p "${OUTPUT_DIR}"
info "Exporting dashboards to: ${OUTPUT_DIR}"

EXPORTED=0
FAILED=0

while IFS= read -r dash; do
  uid=$(echo "${dash}"         | jq -r '.uid')
  title=$(echo "${dash}"       | jq -r '.title')
  folder_uid=$(echo "${dash}"  | jq -r '.folderUid // "general"')
  folder_title=$(folder_lookup "${folder_uid}")
  folder_title="${folder_title:-${folder_uid}}"

  # Sanitise folder and dashboard titles for use as filesystem paths
  safe_folder=$(echo "${folder_title}" | tr '/:*?"<>|\\' '_')
  safe_title=$(echo "${title}"         | tr '/:*?"<>|\\' '_')

  target_dir="${OUTPUT_DIR}/${safe_folder}"
  mkdir -p "${target_dir}"

  target_file="${target_dir}/${safe_title}.json"

  # Fetch the full dashboard JSON
  dash_json=$(grafana_get "/api/dashboards/uid/${uid}" 2>/dev/null) || {
    warn "Failed to fetch dashboard '${title}' (uid=${uid}). Skipping."
    (( FAILED++ )) || true
    continue
  }

  # Extract only the dashboard payload (strip meta wrapper)
  echo "${dash_json}" | jq '.dashboard' > "${target_file}"
  info "  Exported: ${safe_folder}/${safe_title}.json"
  (( EXPORTED++ )) || true

done < <(echo "${DASHBOARDS_JSON}" | jq -c '.[]')

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
echo ""
info "Export complete."
info "  Exported : ${EXPORTED}"
[[ "${FAILED}" -gt 0 ]] && warn "  Failed   : ${FAILED}" || info "  Failed   : ${FAILED}"
info "  Location : ${OUTPUT_DIR}"
