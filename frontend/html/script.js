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

	let challsmap = {};
	challenges.forEach(challenge => {
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
		resultControl.className = 'result-control';
		resultControl.onclick = () => {
			startChallenge(challenge.id);
		}
		resultControl.textContent = '▶';
		result.appendChild(resultControl);

		challs.appendChild(result);
	});

  let insts = document.getElementById('results-instances');
	insts.innerHTML = '';

	instances.reverse();
	instances.forEach(instance => {
		let result = document.createElement('div');
		result.className = 'result';
		result.onclick = () => {
			let fqdn = instance.name + ".tls.vuln.si";
			navigator.clipboard.writeText(fqdn).then(() => {
				let message = document.getElementById('message');
				message.textContent = `FQDN copied to clipboard: ${fqdn}`;
			}).catch(err => {
				console.error('Failed to copy FQDN to clipboard:', err);
			});
		}

		let resultId = document.createElement('div');
		resultId.className = 'result-id';
		resultId.textContent = "ID:"+instance.id;
		result.appendChild(resultId);

		let resultName = document.createElement('div');
		resultName.className = 'result-name';
		resultName.textContent = instance.name + '.tls.vuln.si';
		result.appendChild(resultName);

		let resultChallenge = document.createElement('div');
		resultChallenge.className = 'result-challenge';
		resultChallenge.textContent = challsmap[instance.challenge_id] ? "("+challsmap[instance.challenge_id].name+")" : '(Unknown Challenge)';
		result.appendChild(resultChallenge);

		let resultStartTime = document.createElement('div');
		resultStartTime.className = 'result-starttime';
    resultStartTime.textContent = "[" + instance.created_at + "]";
		result.appendChild(resultStartTime);

		let resultControl = document.createElement('div');
		resultControl.className = 'result-control';
		resultControl.textContent = instance.active ? '🟢' : '🔴';
		result.appendChild(resultControl);

		insts.appendChild(result);
	});
}

function startChallenge(challengeId) {
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
			e.textContent = '⚠️ Error: Challenge with the same ID already started';
		} else if (xhr.status === 404) {
			e.textContent = '⚠️ Error: Challenge not found';
		} else {
			e.textContent = `⚠️ Error: Failed to start challenge (${xhr.status})`;
		}
		e.classList.remove('hidden');
	} else {
    refresh();
	}
}

document.addEventListener('DOMContentLoaded', () => {
	refreshTimed();
});
