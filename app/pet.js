// Intentionally stripped business logic.
// Keep structure and function names so you can implement later.
// WebSocket placeholder
let ws = new WebSocket(`ws://localhost:8082/pet/${pet_id}/ws`)

if (ws) {
  ws.onopen = () => {
	  console.log("Connected to websocket")
  };

  ws.onmessage = (event) => {
	  console.log(event.data)
  };

  ws.onerror = (e) => {
    // TODO: handle error
  };
  ws.onclose = () => {
    // TODO: handle close
  };
}

// Counter logic (local fallback)
const countEl = document.getElementById('pet-count');
const personalCountEl = document.getElementById('user-pet-count');
const petBtn = document.getElementById('pet-container');

let count = 1000;
function getCount() {
  return count;
}
function setCount(n) {
	count = n;
}

function increment() {
  // TODO: implement local increment logic and optionally send to server
	setCount(getCount() + 1);
	countEl.innerText = getCount().toLocaleString();
	personalCountEl.innerText = getCount().toLocaleString();
}

if (petBtn) {
  petBtn.addEventListener('click', () => {
	  countEl.classList.add('count-bump');
	  setTimeout(() => countEl.classList.remove('count-bump'), 320);
	  increment();
	  ws.send('c')
  });
}

// Chat UI
const chatMessages = document.getElementById('chat-message-container');
const chatInput = document.getElementById('chat-input');
const chatSend = document.getElementById('chat-send');

function appendChatMessage(user, msg) {
	const msgEl = document.createElement('div');
	msgEl.innerText = `${user}: ${msg}`;

	msgEl.classList.add('message');
	msgEl.classList.add((user === 'You' ? 'me' : 'them'));

	chatMessages.appendChild(msgEl);
}

function sendChat() {
	const msg = chatInput.value;
	// some crisp optimistic rendering
	appendChatMessage('You', msg);

	data = {msg: msg}
	ws.send(JSON.stringify(data));
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


function setGradientPosition(dog) {

	if (dog == undefined) {
		dog = "daisy"
	}

	const daisyRect = petBtn.getBoundingClientRect();

	const centerX = daisyRect.left + daisyRect.width / 2;
	const centerY = daisyRect.top + daisyRect.height / 2;

	document.body.style.background = `
        radial-gradient(circle at ${centerX}px ${centerY}px, var(--${dog}-gradient-start) 1%, var(--${dog}-gradient-end) 100%
    `;
}

setGradientPosition();
