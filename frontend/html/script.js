var config = {};

function refreshTimed() {
	refresh();
	setTimeout(refreshTimed, 60000);
}

function refresh() {
	let e = document.getElementById('error');
	try {
		let token = getToken();
		if (!token) {
			return;
		}

		let challenges = loadChallenges(token);
		let instances = loadInstances(token);

		render(challenges, instances);

		e.classList.add('hidden');
	} catch (error) {
		console.error(error);
		e.textContent = `⚠️ Error: ${error.message}`;
		e.classList.remove('hidden');
		return;
	}
}

function getConfig() {
	let xhr = new XMLHttpRequest();
	xhr.open('GET', '/api/v1/meta/conf', false);
	xhr.send();

	if (xhr.status !== 200) {
		console.error(`Failed to load config: ${xhr.status} ${xhr.statusText}`);
		throw new Error(`Failed to load config (${xhr.status})`);
	}

	config = JSON.parse(xhr.responseText);
}

function getToken() {
	const urlParams = new URLSearchParams(window.location.search);
	let token = urlParams.get('token');
	if (token) {
		localStorage.setItem('token', token);
	}
	
	token = localStorage.getItem('token');
	if (!token) {
		throw new Error('No token found');
	}

	return token;
}

function loadChallenges(token) {
	let xhr = new XMLHttpRequest();
	xhr.open('GET', '/api/v1/challenge/', false);
	xhr.setRequestHeader('Authorization', token);
	xhr.send();

	if (xhr.status !== 200) {
		console.error(`Failed to load challenges: ${xhr.status} ${xhr.statusText}`);
		throw new Error(`Failed to load challenges (${xhr.status})`);
	}

	let challenges = JSON.parse(xhr.responseText);
	if (!Array.isArray(challenges)) {
		console.error('Invalid challenges format');
		throw new Error('Invalid challenges format');
	}

	return challenges;
}

function loadInstances(token) {
	let xhr = new XMLHttpRequest();
	xhr.open('GET', '/api/v1/instance/', false);
	xhr.setRequestHeader('Authorization', token);
	xhr.send();

	if (xhr.status !== 200) {
		console.error(`Failed to load instances: ${xhr.status} ${xhr.statusText}`);
		throw new Error(`Failed to load instances (${xhr.status})`);
	}

	let instances = JSON.parse(xhr.responseText);
	if (!Array.isArray(instances)) {
		console.error('Invalid instances format');
		throw new Error('Invalid instances format');
	}

	return instances;
}

function render(challenges, instances) {
	let challs = document.getElementById('results-challenges');
	challs.innerHTML = '';

	let categories = {};
	challenges.forEach(challenge => {
		if (challenge.category == null || challenge.category == '') {
			challenge.category = 'other';
		}
		if (!categories[challenge.category]) {
			categories[challenge.category] = [];
		}
		categories[challenge.category].push(challenge);
	});

	let sortedCategories = Object.keys(categories).sort((a, b) => {
		if (a === 'other') return 1; // Move "Other" to the end
		if (b === 'other') return -1;
		return a.localeCompare(b);
	});

	let challsmap = {};
	sortedCategories.forEach(category => {
		let categoryDiv = document.createElement('div');
		categoryDiv.className = 'result-category';

		let h2 = document.createElement('h3');
		h2.textContent = "> " + category;
		categoryDiv.appendChild(h2);

		categories[category].forEach(challenge => {
			challsmap[challenge.id] = challenge;

			let result = document.createElement('div');
			result.className = 'result';

			let resultId = document.createElement('div');
			resultId.className = 'result-id';
			resultId.textContent = "ID:"+challenge.id;
			result.appendChild(resultId);

			let resultName = document.createElement('div');
			resultName.className = 'result-name';
			resultName.textContent = challenge.name;
			result.appendChild(resultName);

			let resultControl = document.createElement('div');
			resultControl.className = 'result-control result-right';
			resultControl.onclick = () => {
				startInstance(challenge.id);
			}

			let resultControlStart = document.createElement('span');
			resultControlStart.className = 'result-control-entry result-clickable';
			resultControlStart.textContent = '▶';
			resultControl.appendChild(resultControlStart);

			result.appendChild(resultControl);
			categoryDiv.appendChild(result);
		});
	
		challs.appendChild(categoryDiv);
	});

  let insts = document.getElementById('results-instances');
	insts.innerHTML = '';

	instances.forEach(instance => {
		let result = document.createElement('div');
		result.className = 'result';

		let resultId = document.createElement('div');
		resultId.className = 'result-id';
		resultId.textContent = "ID:"+instance.id;
		result.appendChild(resultId);

		let resultName = document.createElement('div');
		resultName.className = 'result-name result-clickable';
		resultName.textContent = getFQDN(instance.name, instance.type) + ' 📋';
		resultName.onclick = () => {
			let conn = getConnectionString(instance.name, instance.type);
			navigator.clipboard.writeText(conn).then(() => {
				let message = document.getElementById('message');
				message.textContent = `Connection string copied to clipboard: ${conn}`;
			}).catch(err => {
				console.error('Failed to copy connection string to clipboard:', err);
			});
		}
		result.appendChild(resultName);

		let resultChallenge = document.createElement('div');
		resultChallenge.className = 'result-challenge';
		resultChallenge.textContent = challsmap[instance.challenge_id] ? "("+challsmap[instance.challenge_id].name+")" : '(Unknown Challenge)';
		result.appendChild(resultChallenge);


		if (instance.active) {
			let resultStartTime = document.createElement('div');
			resultStartTime.className = 'result-starttime';
			resultStartTime.textContent = "[" + getFuzzyDuration(instance.created_at, instance.duration) + "]";
			result.appendChild(resultStartTime);
		}

		let resultControl = document.createElement('div');
		resultControl.className = 'result-control result-right';

		if (instance.active) {
			let resultControlExtend = document.createElement('span');
			resultControlExtend.className = 'result-control-entry result-clickable';
			resultControlExtend.textContent = '🕓';
			resultControlExtend.onclick = (e) => {
				extendInstance(instance.id);
			};
			resultControl.appendChild(resultControlExtend);

			let resultControlDelete = document.createElement('span');
			resultControlDelete.className = 'result-control-entry result-clickable';
			resultControlDelete.textContent = '🗑️';
			resultControlDelete.onclick = (e) => {
			  if (confirm('Are you sure you want to delete this instance?')) {
					deleteInstance(instance.id);
				}
			};
			resultControl.appendChild(resultControlDelete);
		}

		let resultControlStatus = document.createElement('span');
		resultControlStatus.className = 'result-control-entry';
		resultControlStatus.textContent = instance.active ? '🟢' : '🔴';
		resultControl.appendChild(resultControlStatus);

		result.appendChild(resultControl);
		insts.appendChild(result);
	});
}

function startInstance(challengeId) {
	let token = getToken();
  if (!token) {
		let e = document.getElementById('error');
		e.textContent = '⚠️ Error: No token found';
		e.classList.remove('hidden');
		return;
	}

	let xhr = new XMLHttpRequest();
	xhr.open('GET', `/api/v1/instance/new/${challengeId}`, false);
	xhr.setRequestHeader('Authorization', token);
	xhr.send();

	if (xhr.status !== 200) {
		console.error(`Failed to start challenge: ${xhr.status} ${xhr.statusText}`);
		let e = document.getElementById('error');
		if (xhr.status === 409) {
			let error = JSON.parse(xhr.responseText);
			e.textContent = `⚠️ Error: ${error.error} (${xhr.status})`;
		} else if (xhr.status === 404) {
			e.textContent = '⚠️ Error: Challenge not found';
		} else {
			let error = xhr.responseText ? JSON.parse(xhr.responseText) : {"error": "Unknown error"};
			e.textContent = `⚠️ Error: Failed to start challenge (${xhr.status}) ${error.error}`;
		}
		e.classList.remove('hidden');
	} else {
    refresh();
	}
}

function extendInstance(instanceId) {
	let token = getToken();
	if (!token) {
		let e = document.getElementById('error');
		e.textContent = '⚠️ Error: No token found';
		e.classList.remove('hidden');
		return;
	}

	let xhr = new XMLHttpRequest();
	xhr.open('POST', `/api/v1/instance/extend/${instanceId}`, false);
	xhr.setRequestHeader('Authorization', token);
	xhr.send();

	if (xhr.status !== 200) {
		let error = xhr.responseText ? JSON.parse(xhr.responseText) : {"error": "Unknown error"};
		console.error(`Failed to extend instance: ${xhr.status} ${xhr.statusText}`);
		let e = document.getElementById('error');
		e.textContent = `⚠️ Error: Failed to extend instance (${xhr.status}) ${error.error}`;
		e.classList.remove('hidden');
	} else {
		refresh();
	}
}

function deleteInstance(instanceId) {
	let token = getToken();
	if (!token) {
		let e = document.getElementById('error');
		e.textContent = '⚠️ Error: No token found';
		e.classList.remove('hidden');
		return;
	}

	let xhr = new XMLHttpRequest();
	xhr.open('DELETE', `/api/v1/instance/${instanceId}`, false);
	xhr.setRequestHeader('Authorization', token);
	xhr.send();

	if (xhr.status !== 204) {
		let error = xhr.responseText ? JSON.parse(xhr.responseText) : {"error": "Unknown error"};
		console.error(`Failed to delete instance: ${xhr.status} ${xhr.statusText}`);
		let e = document.getElementById('error');
		e.textContent = `⚠️ Error: Failed to delete instance (${xhr.status}) ${error.error}`;
		e.classList.remove('hidden');
	} else {
		refresh();
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
	let t = new Date(timestring.replace(' ', 'T') + 'Z');
	if (isNaN(t.getTime())) {
		console.error('Invalid date format:', timestring);
		return 'Invalid date';
	}

	let diff = t.getTime() - new Date().getTime();
	diff = Math.floor(diff / 1000);

	let seconds = diff + duration;
	if (seconds < 0) {
		return 'Cleanup imminent';
	}
	
	if (seconds < 60) {
		return `${seconds}s remaining`;
	} else if (seconds < 3600) {
		return `${Math.floor(seconds / 60)}m remaining`;
	}

	return `${Math.floor(seconds / 3600)}h remaining`;
}

document.addEventListener('DOMContentLoaded', () => {
	getConfig();
	refreshTimed();
});
