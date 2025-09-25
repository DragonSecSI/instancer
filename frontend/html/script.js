var config = {};

async function api_fetch(url, options = {}) {
  const token = getToken();

  // Merge headers, ensuring Authorization is included for API calls
  const headers = {
    ...options.headers,
  };

  if (token && url.startsWith("/api/")) {
    headers.Authorization = token;
  }

  // Merge all options with our headers
  const fetchOptions = {
    ...options,
    headers,
  };

  return fetch(url, fetchOptions);
}

async function show_error(message, response = null) {
  let fullMessage = message;

  if (response && !response.ok) {
    try {
      const error = await response.json();
      if (error.error) {
        fullMessage = `${message} - ${error.error}`;
      } else {
        fullMessage = `${message} (${response.status} ${response.statusText})`;
      }
    } catch (e) {
      // If we can't parse the response, just use the status text
      fullMessage = `${message} (${response.status} ${response.statusText})`;
    }
  }

  let e = document.getElementById("error");
  e.textContent = `⚠️ Error: ${fullMessage}`;
  e.classList.remove("hidden");
}

async function refresh() {
  let e = document.getElementById("error");
  try {
    let token = getToken();
    if (!token) {
      return;
    }

    let challenges = await loadChallenges();
    let instances = await loadInstances();

    render(challenges, instances);
    e.classList.add("hidden");
  } catch (error) {
    console.error(error);
    show_error(error.message);
  }
}

async function getConfig() {
  try {
    const response = await fetch("/api/v1/meta/conf");

    if (!response.ok) {
      console.error(`Failed to load config: ${response.status} ${response.statusText}`);
      throw new Error(`Failed to load config (${response.status})`);
    }

    config = await response.json();
  } catch (error) {
    console.error("Failed to load config:", error);
    throw error;
  }
}

function getToken() {
  const urlParams = new URLSearchParams(window.location.search);
  let token = urlParams.get("token");
  if (token) {
    localStorage.setItem("token", token);
  }

  token = localStorage.getItem("token");
  if (!token) {
    throw new Error("No token found");
  }

  return token;
}

async function loadChallenges() {
  const response = await api_fetch("/api/v1/challenge/");

  if (!response.ok) {
    console.error(`Failed to load challenges: ${response.status} ${response.statusText}`);
    throw new Error(`Failed to load challenges (${response.status})`);
  }

  let challenges = await response.json();
  if (!Array.isArray(challenges)) {
    console.error("Invalid challenges format");
    throw new Error("Invalid challenges format");
  }

  return challenges;
}

async function loadInstances() {
  const response = await api_fetch("/api/v1/instance/");

  if (!response.ok) {
    console.error(`Failed to load instances: ${response.status} ${response.statusText}`);
    throw new Error(`Failed to load instances (${response.status})`);
  }

  let instances = await response.json();
  if (!Array.isArray(instances)) {
    console.error("Invalid instances format");
    throw new Error("Invalid instances format");
  }

  return instances;
}

function setButtonState(button, disabled) {
  button.classList.toggle("disabled", disabled);
}

function render(challenges, instances) {
  let challs = document.getElementById("results-challenges");
  challs.innerHTML = "";

  let categories = {};
  challenges.forEach((challenge) => {
    if (challenge.category == null || challenge.category == "") {
      challenge.category = "other";
    }
    if (!categories[challenge.category]) {
      categories[challenge.category] = [];
    }
    categories[challenge.category].push(challenge);
  });

  let sortedCategories = Object.keys(categories).sort((a, b) => {
    if (a === "other") return 1; // Move "Other" to the end
    if (b === "other") return -1;
    return a.localeCompare(b);
  });

  let challsmap = {};
  sortedCategories.forEach((category) => {
    let categoryDiv = document.createElement("div");
    categoryDiv.className = "result-category";

    let h2 = document.createElement("h3");
    h2.textContent = "> " + category;
    categoryDiv.appendChild(h2);

    categories[category].forEach((challenge) => {
      challsmap[challenge.id] = challenge;

      let result = document.createElement("div");
      result.className = "result";

      let resultId = document.createElement("div");
      resultId.className = "result-id";
      resultId.textContent = "ID:" + challenge.id;
      result.appendChild(resultId);

      let resultName = document.createElement("div");
      resultName.className = "result-name";
      resultName.textContent = challenge.name;
      result.appendChild(resultName);

      let resultControl = document.createElement("div");
      resultControl.className = "result-control result-right";

      let resultControlStart = document.createElement("span");
      resultControlStart.className = "result-control-entry result-clickable";
      resultControlStart.textContent = "▶";

      resultControlStart.onclick = async () => {
        resultControlStart.classList.add("running");
        try {
          await startInstance(challenge.id);
        } finally {
          resultControlStart.classList.remove("running");
        }
      };

      resultControl.appendChild(resultControlStart);
      result.appendChild(resultControl);
      categoryDiv.appendChild(result);
    });

    challs.appendChild(categoryDiv);
  });

  let insts = document.getElementById("results-instances");
  insts.innerHTML = "";

  instances.forEach((instance) => {
    let result = document.createElement("div");
    result.className = "result";

    let resultId = document.createElement("div");
    resultId.className = "result-id";
    resultId.textContent = "ID:" + instance.id;
    result.appendChild(resultId);

    let resultName = document.createElement("div");
    resultName.className = "result-name result-clickable";
    resultName.textContent = getFQDN(instance.name, instance.type) + " 📋";
    resultName.onclick = () => {
      let conn = getConnectionString(instance.name, instance.type);
      navigator.clipboard
        .writeText(conn)
        .then(() => {
          let message = document.getElementById("message");
          message.textContent = `Connection string copied to clipboard: ${conn}`;
        })
        .catch((err) => {
          console.error("Failed to copy connection string to clipboard:", err);
        });
    };
    result.appendChild(resultName);

    let resultChallenge = document.createElement("div");
    resultChallenge.className = "result-challenge";
    resultChallenge.textContent = challsmap[instance.challenge_id]
      ? "(" + challsmap[instance.challenge_id].name + ")"
      : "(Unknown Challenge)";
    result.appendChild(resultChallenge);

    if (instance.active) {
      let resultStartTime = document.createElement("div");
      resultStartTime.className = "result-starttime";
      resultStartTime.textContent = "[" + getFuzzyDuration(instance.created_at, instance.duration) + "]";
      result.appendChild(resultStartTime);
    }

    let resultControl = document.createElement("div");
    resultControl.className = "result-control result-right";

    if (instance.active) {
      let resultControlExtend = document.createElement("span");
      resultControlExtend.className = "result-control-entry result-clickable";
      resultControlExtend.textContent = "🕓";
      resultControlExtend.onclick = async (e) => {
        resultControlExtend.classList.add("running");
        try {
          await extendInstance(instance.id);
        } finally {
          resultControlExtend.classList.remove("running");
        }
      };
      resultControl.appendChild(resultControlExtend);

      let resultControlDelete = document.createElement("span");
      resultControlDelete.className = "result-control-entry result-clickable";
      resultControlDelete.textContent = "🗑️";
      resultControlDelete.onclick = async (e) => {
        if (confirm("Are you sure you want to delete this instance?")) {
          resultControlDelete.classList.add("running");
          try {
            await deleteInstance(instance.id);
          } finally {
            resultControlDelete.classList.remove("running");
          }
        }
      };
      resultControl.appendChild(resultControlDelete);
    }

    let resultControlStatus = document.createElement("span");
    resultControlStatus.className = "result-control-entry";
    resultControlStatus.textContent = instance.active ? "🟢" : "🔴";
    resultControl.appendChild(resultControlStatus);

    result.appendChild(resultControl);
    insts.appendChild(result);
  });
}

async function startInstance(challengeId) {
  const response = await api_fetch(`/api/v1/instance/new/${challengeId}`);

  if (!response.ok) {
    console.error(`Failed to start challenge: ${response.status} ${response.statusText}`);
    await show_error("Failed to start challenge", response);
  } else {
    await refresh();
  }
}

async function extendInstance(instanceId) {
  const response = await api_fetch(`/api/v1/instance/extend/${instanceId}`, {
    method: "POST",
  });

  if (!response.ok) {
    console.error(`Failed to extend instance: ${response.status} ${response.statusText}`);
    await show_error("Failed to extend instance", response);
  } else {
    await refresh();
  }
}

async function deleteInstance(instanceId) {
  const response = await api_fetch(`/api/v1/instance/${instanceId}`, {
    method: "DELETE",
  });

  if (response.status !== 204) {
    console.error(`Failed to delete instance: ${response.status} ${response.statusText}`);
    await show_error("Failed to delete instance", response);
  } else {
    await refresh();
  }
}

function getFQDN(name, type) {
  if (type === 0) {
    return `${name}${config.web_suffix}`;
  } else if (type === 1) {
    return `${name}${config.socket_suffix}`;
  } else {
    return name;
  }
}

function getConnectionString(name, type) {
  if (type === 0) {
    return `https://${name}${config.web_suffix}`;
  } else if (type === 1) {
    return `ncat --ssl ${name}${config.socket_suffix} ${config.socket_port}`;
  } else {
    return name;
  }
}

function getFuzzyDuration(timestring, duration) {
  let t = new Date(timestring.replace(" ", "T") + "Z");
  if (isNaN(t.getTime())) {
    console.error("Invalid date format:", timestring);
    return "Invalid date";
  }

  let diff = t.getTime() - new Date().getTime();
  diff = Math.floor(diff / 1000);

  let seconds = diff + duration;
  if (seconds < 0) {
    return "Cleanup imminent";
  }

  if (seconds < 60) {
    return `${seconds}s remaining`;
  } else if (seconds < 3600) {
    return `${Math.floor(seconds / 60)}m remaining`;
  }

  return `${Math.floor(seconds / 3600)}h remaining`;
}

document.addEventListener("DOMContentLoaded", async () => {
  try {
    await getConfig();
    await refresh();
    setInterval(refresh, 60000);
  } catch (error) {
    console.error("Failed to initialize application:", error);
    await show_error(`Failed to initialize application - ${error.message}`);
  }

  let message = document.getElementById("message");
  message.textContent = "Welcome to the instancer!";

  const chall = new URLSearchParams(window.location.search).get("chall");
  if (chall) {
    startInstance(Number(chall));
  }
});
