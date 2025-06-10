const OLD_GRAFANA_URL ="";
const OLD_API_KEY = "";
const NEW_GRAFANA_URL = "";
const NEW_API_KEY = "";
const NEW_CSS_DS_ID="";
const OLD_CSS_DS_ID="";

const fetchDahsboardUIDs = () => {
  return fetch(`${OLD_GRAFANA_URL}/api/search?type=dash-db`, {
    headers: {
      Authorization: `Bearer ${OLD_API_KEY}`,
    },
  })
    .then((res) => res.json())
    .then((dashboards) => dashboards.map((dash) => dash.uid));
};

const fetchDashboardJSON = (uid) => {
  return fetch(`${OLD_GRAFANA_URL}/api/dashboards/uid/${uid}`, {
    headers: {
      Authorization: `Bearer ${OLD_API_KEY}`,
    },
  })
    .then((res) => res.json())
    .then((res) => {
      res.dashboard.id = null;
      res.dashboard.uid = null;
      return res;
    });
};

const fetchFolders = () => {
  return fetch(`${OLD_GRAFANA_URL}/api/folders`, {
    headers: {
      Authorization: `Bearer ${OLD_API_KEY}`,
    },
  }).then((res) => res.json());
};

const getDashboardUidsForFolder = (folder) => {
  return fetch(
    `${OLD_GRAFANA_URL}/api/search?folderIds=${folder.id}&type=dash-db`,
    {
      headers: {
        Authorization: `Bearer ${OLD_API_KEY}`,
      },
    }
  )
    .then((res) => res.json())
    .then((dashboards) =>
      dashboards.map((dash) => {
        return {
          folderId: folder.id,
          folderName: folder.title,
          dashboardId: dash.uid,
        };
      })
    );
};

const createDashboard = (json, folderId) => {
  return fetch(`${NEW_GRAFANA_URL}/api/dashboards/db`, {
    headers: {
      Authorization: `Bearer ${NEW_API_KEY}`,
      "content-type": "application/json",
    },
    method: "POST",
    body: JSON.stringify({
      dashboard: json.dashboard,
      folderId,
      overwrite: false,
    }),
  }).then((res) => res.json());
};

const createFolder = (name) => {
  return fetch(`${NEW_GRAFANA_URL}/api/folders`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${NEW_API_KEY}`,
      "content-type": "application/json",
    },
    body: JSON.stringify({
      title: name,
    }),
  })
    .then((res) => res.json())
    .then((res) => res.id);
};

const main = async () => {
  const folders = await fetchFolders();
  folders.forEach(async (folder) => {
    const newFolderId = await createFolder(`CSS App - ${folder.title}`);
    const dashboardsUIDs = await getDashboardUidsForFolder(folder);
    for (let dashboardUID of dashboardsUIDs) {
        await fetchDashboardJSON(dashboardUID.dashboardId)
            .then(res => {
                res.dashboard.panels.forEach(panel => {
                    panel.datasource.uid = NEW_CSS_DS_ID
                })
                return res;
            })
            .then((json) => createDashboard(json, newFolderId))
            .then(res => console.log(res))
    }
  });
};

main();
