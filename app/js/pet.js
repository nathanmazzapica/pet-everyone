// Ignore queue is used to keep optimistic rendering from duplicating clicks
let optimisticIgnoreCount = 0;
let socket = null;

const countEl = document.getElementById('pet-count');
const personalCountEl = document.getElementById('user-pet-count');
const petContainer = document.getElementById('pet-container');
const chatMessages = document.getElementById('chat-message-container');
const chatInput = document.getElementById('chat-input');
const chatSend = document.getElementById('chat-send');

let globalCount = 1000;
let localCount = 0;

function getCount() {
	return globalCount;
}

function setCount(n) {
	globalCount = n;
	updateCountUI();
}

function getLocalCount() {
	return localCount;
}

function setLocalCount(n) {
	localCount = n;
	updateCountUI();
}

function updateCountUI() {
	if (countEl) countEl.innerText = getCount().toLocaleString();
	if (personalCountEl) personalCountEl.innerText = getLocalCount().toLocaleString();
}

function increment() {
	setCount(getCount() + 1);
}

// Utility: safe JSON parse
function safeParseJSON(text) {
	try {
		return JSON.parse(text);
	} catch (err) {
		console.error(`Invalid JSON event from server: ${text}`);
		return null;
	}
}

// WEBSOCKET HELPERS

function createWebSocket(petId) {
	if (!petId) return null;
	const token = localStorage.getItem('token');
	const url = `ws://localhost:8082/pet/${petId}/ws?token=${token}`;
	const s = new WebSocket(url);

	s.onopen = () => {
		console.log('Connected to websocket');
	};

	s.onmessage = (e) => handleSocketMessage(e.data);

	s.onerror = (e) => {
		// TODO: handle error (log for now)
		console.error('WebSocket error', e);
	};

	s.onclose = () => {
		// TODO: handle close
		console.log('WebSocket closed');
	};

	return s;
}

function handleSocketMessage(rawData) {
	console.log(rawData);
	const event = safeParseJSON(rawData);
	if (!event) return;

	console.table(event);
	const eventType = event.type;
	switch (eventType) {
		case 'pet':
			if (optimisticIgnoreCount > 0) {
				optimisticIgnoreCount--;
				return;
			}
			increment();
			break;
		case 'petcount':
			setCount(Number(event.data));
			break;
		case 'chat': {
			const msgData = event.data;
			if (msgData) appendChatMessage(msgData.author, msgData.msg);
			break;
		}
		default:
			console.log('Unknown message type', eventType);
	}
}

function sendSocketMessage(obj) {
	if (!socket || socket.readyState !== WebSocket.OPEN) {
		console.warn('WebSocket not ready, cannot send message', obj);
		return;
	}
	socket.send(JSON.stringify(obj));
}

function renderOptimistic() {
	increment();
	optimisticIgnoreCount++;
}

const petCommand = {type: 'pet', data: null};
function pet() {
	if (countEl) {
		countEl.classList.add('count-bump');
		setTimeout(() => countEl.classList.remove('count-bump'), 320);
	}
	setLocalCount(getLocalCount() + 1)
	renderOptimistic();
	sendSocketMessage(petCommand);
}

if (petContainer) {
	petContainer.addEventListener('click', pet);
}

// Chat UI
function appendChatMessage(user, msg) {
	if (!chatMessages) return;
	const msgEl = document.createElement('div');
	msgEl.innerText = `${user}: ${msg}`;
	msgEl.classList.add('message');
	msgEl.classList.add(user === 'You' ? 'me' : 'them');
	chatMessages.appendChild(msgEl);
}

function sendChat() {
	if (!chatInput) return;
	const msg = chatInput.value;
	if (!msg) return;

	// optimistic rendering
	appendChatMessage('You', msg);

	const msgData = {
		type: 'chat',
		data: {
			msg: msg,
			author: 'test'
		}
	};

	sendSocketMessage(msgData);
	chatInput.value = '';
}

if (chatSend) chatSend.addEventListener('click', sendChat);
if (chatInput) {
	chatInput.addEventListener('keydown', (e) => {
		if (e.key === 'Enter') {
			e.preventDefault();
			sendChat();
		}
	});
}

// Visual: gradient based on pet image position
function setGradientPosition(dog = 'daisy') {
	if (!petContainer) return;
	const rect = petContainer.getBoundingClientRect();
	const centerX = rect.left + rect.width / 2;
	const centerY = rect.top + rect.height / 2;
	document.body.style.background = `radial-gradient(circle at ${centerX}px ${centerY}px, var(--${dog}-gradient-start) 1%, var(--${dog}-gradient-end) 100%`;
}

// Initialize
updateCountUI();
setGradientPosition();

// Replace `pet_id` with the actual id variable/name in your scope.
socket = createWebSocket(typeof pet_id !== 'undefined' ? pet_id : null);
